package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (bh BaseHandler) Ping() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		err := bh.dbpool.Ping(ctx)

		if err != nil {
			ctx.Status(http.StatusInternalServerError)
			ctx.Abort()
		}

		ctx.Status(http.StatusOK)
		ctx.Abort()
	}
}
