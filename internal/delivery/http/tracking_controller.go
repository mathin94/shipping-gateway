package http

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"shipping-gateway/internal/model"
	"shipping-gateway/internal/usecase"
)

type TrackingController struct {
	Log             *logrus.Logger
	TrackingUseCase *usecase.TrackingUseCase
}

func NewTrackingController(log *logrus.Logger, trackingUseCase *usecase.TrackingUseCase) *TrackingController {
	return &TrackingController{
		Log:             log,
		TrackingUseCase: trackingUseCase,
	}
}

func (tc *TrackingController) GetTrackingByWaybill(c *gin.Context) {
	log := tc.Log.WithField("traceId", c.Value("traceId"))

	waybill := c.Param("waybill")
	courier := c.Param("courier")

	ucResp, resp := tc.TrackingUseCase.GetTrackingByWaybill(c, waybill, courier)
	if ucResp.StatusCode != 200 {
		log.Errorf("Error getting tracking by waybill: %s, courier: %s, error: %s", waybill, courier, ucResp.Message)
		c.AbortWithStatusJSON(ucResp.StatusCode, model.Response{
			Status:  "failed",
			Code:    ucResp.StatusCode,
			Message: ucResp.Message,
		})
		return
	}

	response := model.ShipmentTrackingResp{
		Response: model.Response{
			Status:  "success",
			Code:    ucResp.StatusCode,
			Message: "success",
		},
		Data: *resp,
	}

	c.JSON(ucResp.StatusCode, response)
}
