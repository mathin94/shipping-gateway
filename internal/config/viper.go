package config

import (
	"fmt"
	"github.com/spf13/viper"
)

// NewViper is a function to load config from config.json
// You can change the implementation, for example load from env file, consul, etcd, etc
func NewViper() *viper.Viper {
	config := viper.New()

	config.SetConfigName("config")
	config.SetConfigType("yaml")
	config.AddConfigPath("./")
	err := config.ReadInConfig()

	SetDefaultValues(config)

	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %w \n", err))
	}

	return config
}

func SetDefaultValues(config *viper.Viper) {
	// Web Server Configuration
	config.SetDefault("web.port", "8080")

	// Database Configuration
	config.SetDefault("db.host", "localhost")
	config.SetDefault("db.port", "3306")
	config.SetDefault("db.user", "root")
	config.SetDefault("db.password", "")
	config.SetDefault("db.name", "mydb")
	config.SetDefault("db.max_open_conns", 100)
	config.SetDefault("db.max_idle_conns", 10)
	config.SetDefault("db.conn_max_lifetime", "30s")

	// Redis Configuration
	config.SetDefault("redis.host", "localhost")
	config.SetDefault("redis.port", "6379")
	config.SetDefault("redis.password", "")
	config.SetDefault("redis.db", 0)

	// Logger Configuration
	config.SetDefault("log.level", "info")
	config.SetDefault("log.console_enabled", true)
	config.SetDefault("log.file_path", "logs/app.log")
	config.SetDefault("log.max_size", 100) // in MB
	config.SetDefault("log.max_backups", 3)
	config.SetDefault("log.max_age", 30) // in days
	config.SetDefault("log.compression", true)

	// Add more default values as needed
}
