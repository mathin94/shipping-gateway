package validator

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"shipping-gateway/internal/model"
)

func ValidateCourierRateRequest(c *gin.Context, req *model.CourierRateRequest) error {
	if err := c.ShouldBindJSON(&req); err != nil {
		return fmt.Errorf("invalid request format: %w", err)
	}

	if req.OriginSubdistrictID == "" && req.OriginPostalCode == "" && req.OriginQuery == "" {
		return fmt.Errorf("invalid request : field 'origin_subdistrict_id', 'origin_postal_code', or 'origin_query' must be provided")
	}

	if req.DestinationSubdistrictID == "" && req.DestinationPostalCode == "" && req.DestinationQuery == "" {
		return fmt.Errorf("invalid request : field 'destination_subdistrict_id', 'destination_postal_code', or 'destination_query' must be provided")
	}

	if req.CourierCode == "" {
		return fmt.Errorf("invalid request : field 'courier_code' must be provided")
	}

	return nil
}
