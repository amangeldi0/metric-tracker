package handlers

import (
	"github.com/amangeldi0/metric-tracker/internal/server/models"
	"github.com/amangeldi0/metric-tracker/internal/server/storage"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (bh BaseHandler) UpdateByURI() gin.HandlerFunc {
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

		storageType := ctx.Param("type")
		if storageType == string(storage.GaugeType) {
			value, err := strconv.ParseFloat(ctx.Param("value"), 64)
			if err != nil {
				bh.handleBadRequest(ctx)
				return
			}

			bh.storage.SetGauge(id, value)
		} else if storageType == string(storage.CounterType) {
			value, err := strconv.ParseInt(ctx.Param("value"), 0, 64)
			if err != nil {
				bh.handleBadRequest(ctx)
				return
			}

			bh.storage.AddCounter(id, value)
		} else {
			bh.handleBadRequest(ctx)
			return
		}

		ctx.Status(http.StatusOK)
		ctx.Abort()
	}
}

func (bh BaseHandler) UpdateByBody() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if !bh.validateContentType(ctx, "application/json", false) {
			bh.handleBadRequest(ctx)
			return
		}

		var obj models.MetricsUpdate
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
			bh.storage.SetGauge(obj.ID, *obj.Value)
		} else if obj.MType == string(storage.CounterType) {
			bh.storage.AddCounter(obj.ID, *obj.Delta)

			counter, err := bh.storage.GetCounter(obj.ID)
			if err != nil {
				bh.log.Errorf("Failed to get updated counter value: %s", err)
			} else {
				obj.Delta = &counter
			}
		}

		ctx.JSON(http.StatusOK, obj)
		ctx.Abort()
	}
}
