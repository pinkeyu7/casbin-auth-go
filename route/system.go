package route

import (
	apiV1 "casbin-auth-go/api/v1"
	"casbin-auth-go/middleware"
	"github.com/gin-gonic/gin"
)

func SystemV1(r *gin.Engine) {
	v1Auth := r.Group("/v1/systems")
	v1Auth.Use(middleware.TokenAuth())

	v1Auth.POST("/", func(c *gin.Context) {
		apiV1.AddSystem(c)
	})

	v1Auth.GET("/", func(c *gin.Context) {
		apiV1.ListSystem(c)
	})

	v1Auth.PUT("/:id", func(c *gin.Context) {
		apiV1.EditSystem(c)
	})
}
