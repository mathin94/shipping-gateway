package config

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"shipping-gateway/external/biteship"
	"shipping-gateway/internal/delivery/http"
	"shipping-gateway/internal/delivery/http/middleware"
	"shipping-gateway/internal/delivery/http/route"
	"shipping-gateway/internal/repository"
	"shipping-gateway/internal/usecase"
)

type BootstrapConfig struct {
	DB       *gorm.DB
	Rds      *redis.Client
	App      *gin.Engine
	Log      *logrus.Logger
	Validate *validator.Validate
	Config   *viper.Viper
}

func Bootstrap(config *BootstrapConfig) {
	// External dependencies
	biteshipClient := biteship.NewClient(config.Config, config.Log)

	// setup repositories
	areaRepository := repository.NewAreaRepository()
	trackingLogRepository := repository.NewTrackingLogRepository()

	//trackingLogRepository := repository.NewTrackingLogRepository()

	// setup use cases
	areaUseCase := usecase.NewAreaUseCase(biteshipClient, config.DB, config.Rds, areaRepository, config.Log)
	shippingUseCase := usecase.NewShippingUseCase(config.DB, config.Log, config.Validate, areaUseCase, biteshipClient, config.Rds)
	trackingUseCase := usecase.NewTrackingUseCase(config.DB, config.Log, biteshipClient, config.Rds, trackingLogRepository)

	// setup controller
	healthCheckController := http.NewHealthCheckController(config.Log)
	courierRateController := http.NewCourierRateController(config.Log, shippingUseCase)
	trackingController := http.NewTrackingController(config.Log, trackingUseCase)

	// setup middleware
	traceIDMiddleware := middleware.TraceIDMiddleware()

	routeConfig := route.RouteConfig{
		App:                   config.App,
		HealthCheckController: healthCheckController,
		CourierRateController: courierRateController,
		TrackingController:    trackingController,
		TraceIDMiddleware:     traceIDMiddleware,
	}

	routeConfig.Setup()
}
