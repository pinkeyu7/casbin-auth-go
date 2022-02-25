package middleware

import (
	"casbin-auth-go/pkg/logr"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func LogRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		meta := logr.ExtractReqMeta(c)

		logr.L.Info("[Request]:",
			zap.String("header-token", meta.Token),
			zap.String("header-appVersion", meta.AppVersion),
			zap.String("header-acceptLanguage", meta.AcceptLanguage),
			zap.String("type", "REQ"),
			zap.String("method", meta.Method),
			zap.String("URI", meta.URI),
			zap.Any("body", meta.Body),
			zap.String("qs", meta.QueryString),
		)

		c.Next()
	}
}
