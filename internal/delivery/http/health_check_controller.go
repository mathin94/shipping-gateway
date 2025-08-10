package http

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

type HealthCheckController struct {
	Log *logrus.Logger
}

func NewHealthCheckController(log *logrus.Logger) *HealthCheckController {
	return &HealthCheckController{
		Log: log,
	}
}

func (h *HealthCheckController) HealthCheck(c *gin.Context) {
	h.Log.Info("Health check endpoint hit")

	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "Service is running",
	})
}
