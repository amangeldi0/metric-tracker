package handlers

import (
	"github.com/amangeldi0/metric-tracker/internal/server/models"
	"github.com/amangeldi0/metric-tracker/internal/server/storage"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (bh BaseHandler) ValueByURI() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if !bh.validateContentType(ctx, "text/plain", true) {
			bh.handleBadRequest(ctx)
			return
		}

		id := ctx.Param("name")
		if id == "" {
			ctx.Status(http.StatusNotFound)
			ctx.Abort()

			return
		}

		var response interface{}
		storageType := ctx.Param("type")

		if storageType == string(storage.GaugeType) {
			value, err := bh.storage.GetGauge(id)
			if err != nil {
				ctx.Status(http.StatusNotFound)
				ctx.Abort()

				return
			}

			response = strconv.FormatFloat(value, 'f', -1, 64)
		} else if storageType == string(storage.CounterType) {
			value, err := bh.storage.GetCounter(id)
			if err != nil {
				ctx.Status(http.StatusNotFound)
				ctx.Abort()

				return
			}

			response = value
		} else {
			bh.handleBadRequest(ctx)
			return
		}

		ctx.String(http.StatusOK, "%v", response)
		ctx.Abort()
	}
}

func (bh BaseHandler) ValueByBody() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if !bh.validateContentType(ctx, "application/json", false) {
			bh.handleBadRequest(ctx)
			return
		}

		var obj models.MetricsValue
		if response, statusCode, err := bh.validateAndShouldBindJSON(ctx, &obj); err != nil {
			if statusCode == http.StatusInternalServerError {
				bh.log.Errorf("Error decoding object request: %s (%T)", err, err)
			}

			if response == nil {
				ctx.Status(statusCode)
			} else {
				ctx.JSON(statusCode, response)
			}

			ctx.Abort()

			return
		}

		if obj.MType == string(storage.GaugeType) {
			value, err := bh.storage.GetGauge(obj.ID)
			if err != nil {
				ctx.Status(http.StatusNotFound)
				ctx.Abort()

				return
			}

			obj.Value = &value
		} else if obj.MType == string(storage.CounterType) {
			delta, err := bh.storage.GetCounter(obj.ID)
			if err != nil {
				ctx.Status(http.StatusNotFound)
				ctx.Abort()

				return
			}

			obj.Delta = &delta
		}

		ctx.JSON(http.StatusOK, obj)
		ctx.Abort()
	}
}
