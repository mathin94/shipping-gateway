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
	"shipping-gateway/internal/model"
	"shipping-gateway/internal/model/converter"
	"shipping-gateway/internal/repository"
	"time"
)

type TrackingUseCase struct {
	DB              *gorm.DB
	Log             *logrus.Logger
	BiteshipClient  *biteship.Client
	Redis           *redis.Client
	TrackingLogRepo *repository.TrackingLogRepository
}

func NewTrackingUseCase(db *gorm.DB, log *logrus.Logger, biteshipClient *biteship.Client, redis *redis.Client, trackingLogRepo *repository.TrackingLogRepository) *TrackingUseCase {
	return &TrackingUseCase{
		DB:              db,
		Log:             log,
		BiteshipClient:  biteshipClient,
		Redis:           redis,
		TrackingLogRepo: trackingLogRepo,
	}
}

func (uc *TrackingUseCase) GetTrackingByWaybill(ctx context.Context, waybill, courier string) (*model.ServiceResponse, *model.ShipmentTrackingResponse) {
	log := uc.Log.WithField("traceId", ctx.Value("traceId"))
	log.Infof("GetTrackingByWaybill called with waybill: %s, courier: %s", waybill, courier)
	var responseData *model.ShipmentTrackingResponse

	// batas waktu harus refetch data dari biteship
	// misal 1 jam, jika sudah lebih dari 1 jam, maka harus refetch data dari biteship
	refetchDuration := 2 * time.Hour
	if waybill == "" || courier == "" {
		return model.BadRequest("Waybill and Courier are required", nil), nil
	}

	rdsKey := fmt.Sprintf("tracking::%s::%s", waybill, courier)
	cachedData, err := uc.Redis.Get(ctx, rdsKey).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		log.Errorf("Error getting cached data: %v", err)
	}

	if cachedData != "" {
		responseData, err = converter.ShipmentTrackingResponseFromJSONString(cachedData)
		if err != nil {
			log.Errorf("Error unmarshalling cached data: %v", err)
		} else {
			if time.Since(responseData.LastTrackedAt) < refetchDuration {
				log.Debugf("Cache hit for waybill: %s, courier: %s", waybill, courier)
				return model.Success(), responseData
			}
		}
	}

	// If cache miss or data is stale, fetch from database first
	trackingLog, err := uc.TrackingLogRepo.FindByWaybillAndCourierCode(uc.DB, waybill, courier)
	if err == nil && trackingLog != nil {
		log.Debugf("Found tracking log in DB for waybill: %s, courier: %s", waybill, courier)
		// convert tracking log to response data
		responseData, errConvert := converter.ShipmentTrackingLogToResponse(trackingLog)
		if errConvert != nil {
			log.Errorf("Error converting tracking log to response: %v", errConvert)
		} else {
			cachedData, errMarshal := json.Marshal(responseData)
			if errMarshal != nil {
				log.Errorf("Error marshalling response data: %v", errMarshal)
			} else {
				errRedis := uc.Redis.Set(ctx, rdsKey, cachedData, refetchDuration).Err()
				if errRedis != nil {
					log.Errorf("Error setting cache: %v", errRedis)
				} else {
					return model.Success(), responseData
				}
			}
		}

		if trackingLog.TrackingStatus == "delivered" || trackingLog.TrackingStatus == "cancelled" {
			log.Infof("Tracking status is %s, no need to refetch from Biteship", trackingLog.TrackingStatus)
			return model.Success(), responseData
		}
	}

	if err != nil {
		log.Errorf("Error finding tracking log: %v", err)
	}

	biteshipResp, biteshipErr := uc.BiteshipClient.GetTrackingByWaybill(waybill, courier)
	if biteshipErr != nil {
		return biteshipErr.ToServiceResponse(), nil
	}

	responseData, err = converter.BiteshipTrackingToResponse(biteshipResp)
	if err != nil {
		log.Errorf("Error converting Biteship response to ShipmentTrackingResponse: %v", err)
		return model.NotFound("Tracking not found"), nil
	}

	trackingLogData, errEntity := converter.ShipmentTrackingLogResponseToEntity(responseData)
	if errEntity != nil {
		log.Errorf("Error converting response data to entity: %v", errEntity)
	} else {
		errSave := uc.TrackingLogRepo.CreateOrUpdate(uc.DB, trackingLogData)
		if errSave != nil {
			log.Warnf("Error saving tracking log: %v", errSave)
		}
	}

	bCache, err := json.Marshal(responseData)
	if err != nil {
		log.Errorf("Error marshalling response data: %v", err)
	}
	err = uc.Redis.Set(ctx, rdsKey, bCache, refetchDuration).Err()
	if err != nil {
		log.Errorf("Error setting cache: %v", err)
	}

	return model.Success(), responseData
}
