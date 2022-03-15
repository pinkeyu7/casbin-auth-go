package middleware

import (
	"casbin-auth-go/pkg/er"
	"casbin-auth-go/pkg/logr"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

func ErrorResponse() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		err := c.Errors.Last()

		if err == nil {
			return
		}

		logger := logr.NewCtxLogger(c)

		switch v := err.Err; v.(type) {
		case validator.ValidationErrors:
			var msgs []string
			vErr := err.Err.(validator.ValidationErrors)
			for _, e := range vErr {
				msgs = append(msgs, e.Field())
			}

			msg := strings.Join(msgs, " , ") + " is wrong format or invalid"
			appErr := er.NewAppErr(http.StatusBadRequest, er.ErrorParamInvalid, msg, vErr)
			errMsgBytes, _ := json.Marshal(appErr.GetMsg())

			logger.Error(
				string(errMsgBytes),
				zap.String("type", "API"),
				zap.Int("e_status", appErr.GetStatus()),
				zap.String("e_code", appErr.Code),
				zap.String("cause", vErr.Error()),
			)

			c.AbortWithStatusJSON(appErr.GetStatus(), appErr.GetMsg())
			return
		case *er.AppError:
			appErr := err.Err.(*er.AppError)

			// auto log error, only 500 error will log now
			if appErr.StatusCode == http.StatusInternalServerError || appErr.StatusCode == http.StatusBadRequest {
				errMsgBytes, _ := json.Marshal(appErr.GetMsg())
				causeBytes, _ := json.Marshal(appErr.CauseErr)

				logger.Error(
					string(errMsgBytes),
					zap.String("type", "API"),
					zap.Int("e_status", appErr.GetStatus()),
					zap.String("e_code", appErr.Code),
					zap.String("cause", string(causeBytes)),
				)
			}

			c.AbortWithStatusJSON(appErr.GetStatus(), appErr.GetMsg())
			return
		default:
		}
	}
}
