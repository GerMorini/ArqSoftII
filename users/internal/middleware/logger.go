package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func LoggerMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		ctx.Next()
		elapsed := time.Since(start)

		if ctx.Request.Method == "OPTIONS" {
			return
		}

		status := ctx.Writer.Status()
		method := ctx.Request.Method
		uri := ctx.Request.RequestURI

		logMsg := fmt.Sprintf("%s | %d | %s | %s", method, status, elapsed.String(), uri)

		if len(ctx.Errors) > 0 {
			logMsg += " errors=" + ctx.Errors.String()
		}

		if status >= 500 {
			log.Error(logMsg)
		} else if status >= 400 {
			log.Warn(logMsg)
		} else {
			log.Info(logMsg)
		}
	}
}
