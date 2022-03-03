package route

import (
	apiV1 "casbin-auth-go/api/v1"
	"casbin-auth-go/middleware"
	"github.com/gin-gonic/gin"
)

func SysPermissionV1(r *gin.Engine) {
	v1Auth := r.Group("/v1/permissions")
	v1Auth.Use(middleware.TokenAuth())

	// 列表 系統權限
	v1Auth.GET("/", func(c *gin.Context) {
		apiV1.ListPermission(c)
	})

	// 新增 系統權限
	v1Auth.POST("/", func(c *gin.Context) {
		apiV1.AddPermission(c)
	})

	// 編輯 系統權限
	v1Auth.PUT("/:id", func(c *gin.Context) {
		apiV1.EditPermission(c)
	})

	// 刪除 系統權限
	v1Auth.DELETE("/:id", func(c *gin.Context) {
		apiV1.DeletePermission(c)
	})
}
