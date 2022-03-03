package route

import (
	apiV1 "casbin-auth-go/api/v1"
	"casbin-auth-go/middleware"
	"github.com/gin-gonic/gin"
)

func SysRoleV1(r *gin.Engine) {
	v1Auth := r.Group("/v1/roles")
	v1Auth.Use(middleware.TokenAuth())

	// 列表 系統角色
	v1Auth.GET("/", func(c *gin.Context) {
		apiV1.ListSysRole(c)
	})

	// 取得 系統角色 及 權限列表
	v1Auth.GET("/:id", func(c *gin.Context) {
		apiV1.GetSysRoleWithPermission(c)
	})

	// 新增 系統角色
	v1Auth.POST("/", func(c *gin.Context) {
		apiV1.AddSysRole(c)
	})

	// 編輯 系統角色
	v1Auth.PUT("/:id", func(c *gin.Context) {
		apiV1.EditSysRole(c)
	})
}
