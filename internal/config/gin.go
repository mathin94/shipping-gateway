package config

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net/http"
)

func NewGinEngine(config *viper.Viper) *gin.Engine {
	app := gin.New()
	app.Use(gin.Recovery())
	app.Use(NewErrorHandler())
	// You can set app name in context or use it elsewhere as needed
	return app
}

func NewErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if len(c.Errors) > 0 {
			c.JSON(-1, gin.H{
				"errors": c.Errors[0].Error(),
			})
			c.AbortWithStatus(http.StatusInternalServerError)
		}
	}
}
