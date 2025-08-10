package usecase

import (
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"shipping-gateway/external/biteship"
	"shipping-gateway/internal/model"
	"strconv"
)

type ShippingUseCase struct {
	DB             *gorm.DB
	Log            *logrus.Logger
	Validate       *validator.Validate
	AreaUseCase    *AreaUseCase
	BiteshipClient *biteship.Client
	Redis          *redis.Client
}

func NewShippingUseCase(db *gorm.DB, log *logrus.Logger, validate *validator.Validate,
	areaUseCase *AreaUseCase, bc *biteship.Client, redis *redis.Client) *ShippingUseCase {
	return &ShippingUseCase{
		AreaUseCase:    areaUseCase,
		BiteshipClient: bc,
		DB:             db,
		Log:            log,
		Validate:       validate,
		Redis:          redis,
	}
}

func (uc *ShippingUseCase) GetCourierRates(ctx context.Context, req *model.CourierRateRequest) (*model.ServiceResponse, *model.CourierRateResponse) {
	log := uc.Log.WithField("traceId", ctx.Value("traceId"))
	log.Infof("GetCourierRates request: %+v", req)

	var bsReq biteship.RateRequest

	originArea, err := uc.AreaUseCase.FindArea(ctx, req.OriginSubdistrictID, req.OriginPostalCode, req.OriginQuery)
	if err != nil || originArea == nil || originArea.ExternalID == "" {
		if req.OriginPostalCode == "" {
			log.Errorf("Origin postal code is required when origin area is not found: %+v", req)
			return model.BadRequest("Origin postal code is required", nil), nil
		}
		postalCode, _ := strconv.Atoi(req.OriginPostalCode)
		bsReq.OriginPostalCode = postalCode
	} else {
		bsReq.OriginAreaID = originArea.ExternalID
	}

	destinationArea, err := uc.AreaUseCase.FindArea(ctx, req.DestinationSubdistrictID, req.DestinationPostalCode, req.DestinationQuery)
	if err != nil || destinationArea == nil || destinationArea.ExternalID == "" {
		if req.DestinationPostalCode == "" {
			log.Errorf("Destination postal code is required when destination area is not found: %+v", req)
			return model.BadRequest("Destination postal code is required", nil), nil
		}
		postalCode, _ := strconv.Atoi(req.DestinationPostalCode)
		bsReq.DestinationPostalCode = postalCode
	} else {
		bsReq.DestinationAreaID = destinationArea.ExternalID
	}

	bsReq.Couriers = req.CourierCode
	bsReq.Items = make([]biteship.Item, 0, len(req.Items))
	for _, item := range req.Items {
		bsReq.Items = append(bsReq.Items, biteship.Item{
			Name:        item.Name,
			Description: item.Description,
			Value:       item.Price,
			Weight:      item.Weight,
			Length:      item.Length,
			Width:       item.Width,
			Height:      item.Height,
			Quantity:    item.Quantity,
		})
	}

	rateResponse, errResp := uc.BiteshipClient.GetCourierRates(bsReq)
	if errResp != nil {
		if errResp.IsEmptyData() {
			log.Infof("No courier rates found for request: %+v", req)
			return model.NotFound("No courier rates found"), nil
		}
		log.Errorf("Error getting courier rates from Biteship: %v", errResp)
		return model.DefaultError("Failed to get courier rates", nil), nil
	}

	resp := rateResponse.ToCourierRateResponse()
	return model.Success(), &resp
}
