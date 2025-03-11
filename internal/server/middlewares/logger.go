package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func (bm BaseMiddleware) Logger(ctx *gin.Context) {
	start := time.Now()

	ctx.Next()

	method := ctx.Request.Method
	uri := ctx.Request.URL
	duration := time.Since(start)
	statusCode := ctx.Writer.Status()
	size := ctx.Writer.Size()

	bm.log.Infof("%s Request - URI: \"%s\" - LeadTime: %v - StatusCode: %d (%s) - BodySize: %d", method, uri, duration, statusCode, http.StatusText(statusCode), size)
}
