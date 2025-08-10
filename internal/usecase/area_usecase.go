package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"shipping-gateway/external/biteship"
	"shipping-gateway/internal/entity"
	"shipping-gateway/internal/model"
	"shipping-gateway/internal/model/converter"
	"shipping-gateway/internal/repository"
	"strconv"
)

type AreaUseCase struct {
	BiteshipClient *biteship.Client
	DB             *gorm.DB
	Redis          *redis.Client
	AreaRepo       *repository.AreaRepository
	Logger         *logrus.Logger
}

func NewAreaUseCase(bs *biteship.Client, db *gorm.DB, redis *redis.Client, areaRepo *repository.AreaRepository, logger *logrus.Logger) *AreaUseCase {
	return &AreaUseCase{
		DB:             db,
		Redis:          redis,
		AreaRepo:       areaRepo,
		Logger:         logger,
		BiteshipClient: bs,
	}
}

func (a *AreaUseCase) FindArea(ctx context.Context, subdistrictID, postalCode, query string) (*entity.Area, error) {
	log := a.Logger.WithField("traceId", ctx.Value("traceId"))
	log.Infof("Finding area with subdistrictID: %s, postalCode: %s, query: %s", subdistrictID, postalCode, query)

	var area *entity.Area

	rdsKey := fmt.Sprintf("area::%s::%s", subdistrictID, postalCode)
	strArea, err := a.Redis.Get(ctx, rdsKey).Result()
	if err == nil {
		log.Infof("Cache hit for area with key: %s", rdsKey)
		if area, err = converter.AreaFromJSONString(strArea); err != nil {
			log.Errorf("Error converting area from JSON string: %v", err)
		}

		if area != nil && area.ExternalID != "" {
			return area, nil
		}
	}

	if err != nil && !errors.Is(err, redis.Nil) {
		log.Errorf("Error getting area from Redis: %v", err)
	}

	// If not found in Redis, query the database
	log.Infof("Cache miss for area with subdistrictID: %s, postalCode: %s", subdistrictID, postalCode)
	area, err = a.AreaRepo.FindByParameter(a.DB, model.FindAreaParams{
		OriginalSubdistrictID: subdistrictID,
		OriginalPostalCode:    postalCode,
	})

	if err != nil && area != nil && area.ExternalID != "" {
		// Area found in database, save it to Redis
		strAreaNew, err := json.Marshal(area)
		if err != nil {
			log.Errorf("Error converting area to JSON string: %v", err)
			goto findFromBiteship
		}
		rdsKey = fmt.Sprintf("area::%d::%s", area.OriginalSubdistrictID, area.OriginalPostalCode)
		if err = a.Redis.Set(ctx, rdsKey, strAreaNew, 0).Err(); err != nil {
			log.Errorf("Error setting area in Redis: %v", err)
		}

		return area, nil
	}

findFromBiteship:
	if query == "" {
		log.Infof("Skipping area search from Biteship due to empty query")
		return nil, fmt.Errorf("area query is empty, cannot search area")
	}

	biteshipArea, errResp := a.BiteshipClient.SearchAreas(query)
	if errResp != nil {
		log.Errorf("Error finding area from Biteship: %s", errResp.Error)
		return nil, fmt.Errorf("failed to find area from Biteship: %s", errResp.Error)
	}

	if len(biteshipArea.Areas) == 0 {
		log.Infof("No area found from Biteship for query: %s", query)
		return nil, fmt.Errorf("no area found for query: %s", query)
	}

	// Get the first area from Biteship response
	biteshipFirstArea := biteshipArea.Areas[0]
	var iSubdistrictID int
	// Create a new entity.Area from Biteship response
	if subdistrictID != "" {
		iSubdistrictID, err = strconv.Atoi(subdistrictID)
		if err != nil {
			log.Errorf("Invalid subdistrict ID: %s, error: %v", subdistrictID, err)
			return nil, fmt.Errorf("invalid subdistrict ID: %s", subdistrictID)
		}
	}

	// convert biteshipFirstArea to json string
	strBiteshipArea, err := json.Marshal(biteshipFirstArea)
	if err != nil {
		log.Errorf("Failed to marshal Biteship area: %v", err)
		return nil, fmt.Errorf("failed to marshal Biteship area: %w", err)
	}

	area = &entity.Area{
		OriginalSubdistrictID: uint(iSubdistrictID),
		OriginalPostalCode:    strconv.Itoa(biteshipFirstArea.PostalCode),
		Description:           biteshipFirstArea.Name,
		ExternalSource:        "biteship",
		ExternalID:            biteshipFirstArea.ID,
		ExternalInfo:          string(strBiteshipArea),
	}

	bArea, _ := json.Marshal(area)
	if iSubdistrictID > 0 {
		// Save the area to the database
		err = a.AreaRepo.SaveOrCreate(a.DB, area)
		if err != nil {
			log.Errorf("Failed to save area to database: %v", err)
			return nil, fmt.Errorf("failed to save area to database: %w", err)
		}
	}

	// Save the area to Redis
	if err = a.Redis.Set(ctx, rdsKey, bArea, 0).Err(); err != nil {
		log.Errorf("Error setting area in Redis: %v", err)
	}

	return area, nil
}
