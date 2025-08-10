package main

import (
	"fmt"
	"shipping-gateway/internal/config"
)

func main() {
	viperConfig := config.NewViper()
	log := config.NewLogger(viperConfig)
	db := config.NewDatabase(viperConfig, log)
	rds := config.InitRedis(viperConfig, log)
	log.Infof("Database connection established successfully")
	validate := config.NewValidator(viperConfig)
	app := config.NewGinEngine(viperConfig)

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
