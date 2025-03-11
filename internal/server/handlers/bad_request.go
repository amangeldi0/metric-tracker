package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (bh BaseHandler) BadRequest(ctx *gin.Context) {
	bh.handleBadRequest(ctx)
}

func (bh BaseHandler) handleBadRequest(ctx *gin.Context) {
	ctx.Status(http.StatusBadRequest)
	ctx.Abort()
}
