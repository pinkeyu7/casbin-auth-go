package route

import (
	apiV1 "casbin-auth-go/api/v1"
	"casbin-auth-go/middleware"

	"github.com/gin-gonic/gin"
)

func SysAccountV1(r *gin.Engine) {
	v1Auth := r.Group("/v1/accounts")
	v1Auth.Use(middleware.TokenAuth())

	// 列表 系統帳號
	v1Auth.GET("/", func(c *gin.Context) {
		apiV1.ListSysAccount(c)
	})

	// 取得 系統帳號
	v1Auth.GET("/:id", func(c *gin.Context) {
		apiV1.GetSysAccount(c)
	})

	// 新增 系統帳號
	v1Auth.POST("/", func(c *gin.Context) {
		apiV1.AddSysAccount(c)
	})

	// 編輯 系統帳號
	v1Auth.PUT("/:id", func(c *gin.Context) {
		apiV1.EditSysAccount(c)
	})

	// 系統管理者幫系統帳號重設密碼
	v1Auth.PUT("/:id/forgot-password", func(c *gin.Context) {
		apiV1.ForgotPasswordByAdmin(c)
	})

	// 系統帳號更改密碼Ï
	v1Auth.PUT("/:id/change-password", func(c *gin.Context) {
		apiV1.ChangePassword(c)
	})
}
