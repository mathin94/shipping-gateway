package http

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"shipping-gateway/internal/delivery/http/validator"
	"shipping-gateway/internal/model"
	"shipping-gateway/internal/usecase"
)

type CourierRateController struct {
	Log             *logrus.Logger
	ShippingUseCase *usecase.ShippingUseCase
}

func NewCourierRateController(log *logrus.Logger, shippingUseCase *usecase.ShippingUseCase) *CourierRateController {
	return &CourierRateController{
		Log:             log,
		ShippingUseCase: shippingUseCase,
	}
}

func (c *CourierRateController) GetCourierRates(ctx *gin.Context) {
	log := c.Log.WithField("traceId", ctx.Value("traceId"))

	var req model.CourierRateRequest
	err := validator.ValidateCourierRateRequest(ctx, &req)
	if err != nil {
		log.Errorf("Invalid request parameters, error: %v", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, model.Response{Message: err.Error(), Status: "failed"})
		return
	}

	ucResp, resp := c.ShippingUseCase.GetCourierRates(ctx, &req)
	if ucResp.StatusCode != http.StatusOK {
		ctx.AbortWithStatusJSON(ucResp.StatusCode, model.Response{
			Status:  "failed",
			Code:    ucResp.StatusCode,
			Message: ucResp.Message,
		})
		return
	}
	ctx.JSON(ucResp.StatusCode, resp)
}
