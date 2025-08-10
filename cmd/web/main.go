package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"shipping-gateway/internal/config"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	viperConfig := config.NewViper()
	log := config.NewLogger(viperConfig)
	app := config.NewGinEngine(viperConfig)
	db := config.NewDatabase(viperConfig, log)
	rds := config.InitRedis(viperConfig, log)
	validate := config.NewValidator(viperConfig)

	config.Bootstrap(&config.BootstrapConfig{
		DB:       db,
		Rds:      rds,
		App:      app,
		Log:      log,
		Validate: validate,
		Config:   viperConfig,
	})

	webPort := viperConfig.GetString("web.port")
	err := app.Run(fmt.Sprintf(":%s", webPort))
	if err != nil {
		log.Fatalf("Failed to start web server: %v", err)
	}
}
