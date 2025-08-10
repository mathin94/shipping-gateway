package route

import (
	"github.com/gin-gonic/gin"
	"shipping-gateway/internal/delivery/http"
)

type RouteConfig struct {
	App *gin.Engine

	// Add controller below
	HealthCheckController *http.HealthCheckController
	CourierRateController *http.CourierRateController

	// Add middleware below
	TraceIDMiddleware gin.HandlerFunc
}

func (c *RouteConfig) Setup() {
	c.SetupGuestRoute()
	c.SetupInternalRoute()
}

// SetupGuestRoute is used to setup routes that can be accessed by guest users without authentication.
// This function should be called to define routes that do not require any authentication or special permissions.
// It is typically used for public endpoints that are accessible to all users, such as health checks
func (c *RouteConfig) SetupGuestRoute() {
	v1 := c.App.Group("/api/v1")

	// Example of setting up a guest route
	v1.GET("/health-check", c.HealthCheckController.HealthCheck)
}

// SetupInternalRoute is used to setup routes that can be accessed by internal users with authentication.
// This function should be called to define routes that require authentication and are intended for internal use only
func (c *RouteConfig) SetupInternalRoute() {
	v1 := c.App.Group("/api/v1")
	v1.Use(c.TraceIDMiddleware)
	// Shipping routes
	shippingV1 := v1.Group("/shipping")
	shippingV1.POST("/rates", c.CourierRateController.GetCourierRates)
}
