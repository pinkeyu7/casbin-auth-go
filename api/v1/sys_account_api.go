package v1

import (
	"casbin-auth-go/api"
	"casbin-auth-go/dto/apireq"
	sysAccRepo "casbin-auth-go/internal/system/sys_account/repository"
	sysAccSrv "casbin-auth-go/internal/system/sys_account/service"
	sysRoleRepo "casbin-auth-go/internal/system/sys_role/repository"
	sysRepo "casbin-auth-go/internal/system/system/repository"
	tokenLibrary "casbin-auth-go/internal/token/library"
	tokenRepo "casbin-auth-go/internal/token/repository"
	"casbin-auth-go/pkg/er"
	"casbin-auth-go/pkg/helper"
	"casbin-auth-go/pkg/valider"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ListSysAccount
// @Summary List System Account 取得帳號列表
// @Produce json
// @Accept json
// @Tags Account
// @Security Bearer
// @Param Bearer header string true "JWT Token"
// @Param account_id query int true "Account ID"
// @Param page query string true "Page"
// @Param per_page query string true "Per Page"
// @Param system_id query int false "System ID"
// @Success 200 {object} apires.ListSysAccount
// @Failure 400 {object} er.AppErrorMsg "{"code":"400400","message":"Wrong parameter format or invalid"}"
// @Failure 401 {object} er.AppErrorMsg "{"code":"400401","message":"Unauthorized"}"
// @Failure 403 {object} er.AppErrorMsg "{"code":"400403","message":"Permission denied"}"
// @Failure 404 {object} er.AppErrorMsg "{"code":"400404","message":"Resource not found"}"
// @Failure 500 {object} er.AppErrorMsg "{"code":"500000","message":"Database unknown error"}"
// @Router /v1/accounts [get]
func ListSysAccount(c *gin.Context) {
	req := apireq.ListSysAccount{}
	err := c.Bind(&req)
	if err != nil {
		err = er.NewAppErr(http.StatusBadRequest, er.ErrorParamInvalid, err.Error(), err)
		_ = c.Error(err)
		return
	}

	// 參數驗證
	err = valider.Validate.Struct(req)
	if err != nil {
		err = er.NewAppErr(http.StatusBadRequest, er.ErrorParamInvalid, err.Error(), err)
		_ = c.Error(err)
		return
	}

	// 驗證 jwt user == user_id
	err = tokenLibrary.CheckJWTAccountId(c, req.AccountId)
	if err != nil {
		_ = c.Error(err)
		return
	}

	env := api.GetEnv()
	tc := tokenRepo.NewCache(env.RedisCluster)
	sc := sysRepo.NewCache(env.RedisCluster)
	sr := sysRepo.NewRepository(env.Orm, sc)
	srr := sysRoleRepo.NewRepository(env.Orm)
	sar := sysAccRepo.NewRepository(env.Orm)
	sas := sysAccSrv.NewService(sar, sr, srr, tc)

	res, err := sas.ListSysAccount(&req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, res)
}

// AddSysAccount
// @Summary Add System Account 新增帳號
// @Produce json
// @Accept json
// @Tags Account
// @Security Bearer
// @Param Bearer header string true "JWT Token"
// @Param Body body apireq.AddSysAccount true "Request Add System Account"
// @Success 200 {string} string "{}"
// @Failure 400 {object} er.AppErrorMsg "{"code":"400400","message":"Wrong parameter format or invalid"}"
// @Failure 401 {object} er.AppErrorMsg "{"code":"400401","message":"Unauthorized"}"
// @Failure 403 {object} er.AppErrorMsg "{"code":"400403","message":"Permission denied"}"
// @Failure 404 {object} er.AppErrorMsg "{"code":"400404","message":"Resource not found"}"
// @Failure 500 {object} er.AppErrorMsg "{"code":"500000","message":"Database unknown error"}"
// @Router /v1/accounts [post]
func AddSysAccount(c *gin.Context) {
	req := apireq.AddSysAccount{}
	err := c.BindJSON(&req)
	if err != nil {
		paramErr := er.NewAppErr(http.StatusBadRequest, er.ErrorParamInvalid, err.Error(), err)
		_ = c.Error(paramErr)
		return
	}

	// 參數驗證
	err = valider.Validate.Struct(req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	// 驗證 jwt user == user_id
	err = tokenLibrary.CheckJWTAccountId(c, req.AccountId)
	if err != nil {
		_ = c.Error(err)
		return
	}

	env := api.GetEnv()
	tc := tokenRepo.NewCache(env.RedisCluster)
	sc := sysRepo.NewCache(env.RedisCluster)
	sr := sysRepo.NewRepository(env.Orm, sc)
	srr := sysRoleRepo.NewRepository(env.Orm)
	sar := sysAccRepo.NewRepository(env.Orm)
	sas := sysAccSrv.NewService(sar, sr, srr, tc)

	err = sas.AddSysAccount(&req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{})
}

// EditSysAccount
// @Summary Edit System Account 編輯帳號
// @Produce json
// @Accept json
// @Tags Account
// @Security Bearer
// @Param Bearer header string true "JWT Token"
// @Param target_id path int true "Target Account ID"
// @Param Body body apireq.EditSysAccount true "Request Edit System Account"
// @Success 200 {string} string "{}"
// @Failure 400 {object} er.AppErrorMsg "{"code":"400400","message":"Wrong parameter format or invalid"}"
// @Failure 401 {object} er.AppErrorMsg "{"code":"400401","message":"Unauthorized"}"
// @Failure 403 {object} er.AppErrorMsg "{"code":"400403","message":"Permission denied"}"
// @Failure 404 {object} er.AppErrorMsg "{"code":"400404","message":"Resource not found"}"
// @Failure 500 {object} er.AppErrorMsg "{"code":"500000","message":"Database unknown error"}"
// @Router /v1/accounts/{target_id} [put]
func EditSysAccount(c *gin.Context) {
	id := c.Param("id")
	sysAccId, err := strconv.Atoi(id)
	if err != nil {
		err = er.NewAppErr(http.StatusBadRequest, er.ErrorParamInvalid, "system account id format error.", err)
		_ = c.Error(err)
		return
	}

	req := apireq.EditSysAccount{}
	err = c.BindJSON(&req)
	if err != nil {
		paramErr := er.NewAppErr(http.StatusBadRequest, er.ErrorParamInvalid, err.Error(), err)
		_ = c.Error(paramErr)
		return
	}

	// 參數驗證
	err = valider.Validate.Struct(req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	// 驗證 jwt user == user_id
	err = tokenLibrary.CheckJWTAccountId(c, req.AccountId)
	if err != nil {
		_ = c.Error(err)
		return
	}

	env := api.GetEnv()
	tc := tokenRepo.NewCache(env.RedisCluster)
	sc := sysRepo.NewCache(env.RedisCluster)
	sr := sysRepo.NewRepository(env.Orm, sc)
	srr := sysRoleRepo.NewRepository(env.Orm)
	sar := sysAccRepo.NewRepository(env.Orm)
	sas := sysAccSrv.NewService(sar, sr, srr, tc)

	err = sas.EditSysAccount(sysAccId, &req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{})
}

// GetSysAccount
// @Summary Get System Account 取得帳號資料
// @Produce json
// @Accept json
// @Tags Account
// @Security Bearer
// @Param Bearer header string true "JWT Token"
// @Param target_id path int true "Target Account ID"
// @Param account_id query int true "Account ID"
// @Success 200 {object} apires.SysAccount
// @Failure 400 {object} er.AppErrorMsg "{"code":"400400","message":"Wrong parameter format or invalid"}"
// @Failure 401 {object} er.AppErrorMsg "{"code":"400401","message":"Unauthorized"}"
// @Failure 403 {object} er.AppErrorMsg "{"code":"400403","message":"Permission denied"}"
// @Failure 404 {object} er.AppErrorMsg "{"code":"400404","message":"Resource not found"}"
// @Failure 500 {object} er.AppErrorMsg "{"code":"500000","message":"Database unknown error"}"
// @Router /v1/accounts/{target_id} [get]
func GetSysAccount(c *gin.Context) {
	id := c.Param("id")
	sysAccId, err := strconv.Atoi(id)
	if err != nil {
		err = er.NewAppErr(http.StatusBadRequest, er.ErrorParamInvalid, "system account id format error.", err)
		_ = c.Error(err)
		return
	}

	req := apireq.GetSysRole{}
	err = c.Bind(&req)
	if err != nil {
		paramErr := er.NewAppErr(http.StatusBadRequest, er.ErrorParamInvalid, err.Error(), err)
		_ = c.Error(paramErr)
		return
	}

	// 參數驗證
	err = valider.Validate.Struct(req)
	if err != nil {
		paramErr := er.NewAppErr(http.StatusBadRequest, er.ErrorParamInvalid, err.Error(), err)
		_ = c.Error(paramErr)
		return
	}

	// 驗證 jwt user == user_id
	err = tokenLibrary.CheckJWTAccountId(c, req.AccountId)
	if err != nil {
		_ = c.Error(err)
		return
	}

	env := api.GetEnv()
	tc := tokenRepo.NewCache(env.RedisCluster)
	sc := sysRepo.NewCache(env.RedisCluster)
	sr := sysRepo.NewRepository(env.Orm, sc)
	srr := sysRoleRepo.NewRepository(env.Orm)
	sar := sysAccRepo.NewRepository(env.Orm)
	sas := sysAccSrv.NewService(sar, sr, srr, tc)

	res, err := sas.GetSysAccount(sysAccId)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, res)
}

// ForgotPasswordByAdmin
// @Summary Forgot System Account Password 忘記系統帳號密碼
// @Produce json
// @Accept json
// @Tags Account
// @Security Bearer
// @Param Bearer header string true "JWT Token"
// @Param account_id path int true "Account ID"
// @Param Body body apireq.ForgotSysAccountPassword true "Request Forgot System Account Password"
// @Success 200 {string} string "{}"
// @Failure 400 {object} er.AppErrorMsg "{"code":"400400","message":"Wrong parameter format or invalid"}"
// @Failure 401 {object} er.AppErrorMsg "{"code":"400401","message":"Unauthorized"}"
// @Failure 403 {object} er.AppErrorMsg "{"code":"400403","message":"Permission denied"}"
// @Failure 404 {object} er.AppErrorMsg "{"code":"400404","message":"Resource not found"}"
// @Failure 500 {object} er.AppErrorMsg "{"code":"500000","message":"Database unknown error"}"
// @Router /v1/accounts/{account_id}/forgot-password [put]
func ForgotPasswordByAdmin(c *gin.Context) {
	id := c.Param("id")
	sysAccId, err := strconv.Atoi(id)
	if err != nil {
		err = er.NewAppErr(http.StatusBadRequest, er.ErrorParamInvalid, "system account id format error.", err)
		_ = c.Error(err)
		return
	}

	req := apireq.ForgotSysAccountPassword{}
	err = c.BindJSON(&req)
	if err != nil {
		paramErr := er.NewAppErr(http.StatusBadRequest, er.ErrorParamInvalid, err.Error(), err)
		_ = c.Error(paramErr)
		return
	}

	// 參數驗證
	err = valider.Validate.Struct(req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	// 驗證 jwt user == user_id
	err = tokenLibrary.CheckJWTAccountId(c, req.AccountId)
	if err != nil {
		_ = c.Error(err)
		return
	}

	env := api.GetEnv()
	tc := tokenRepo.NewCache(env.RedisCluster)
	sc := sysRepo.NewCache(env.RedisCluster)
	sr := sysRepo.NewRepository(env.Orm, sc)
	srr := sysRoleRepo.NewRepository(env.Orm)
	sar := sysAccRepo.NewRepository(env.Orm)
	sas := sysAccSrv.NewService(sar, sr, srr, tc)

	err = sas.ForgotPasswordByAdmin(sysAccId, &req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{})
}

// ChangePassword
// @Summary Change Password 更改系統帳號密碼
// @Produce json
// @Accept json
// @Tags Account
// @Security Bearer
// @Param Bearer header string true "JWT Token"
// @Param account_id path int true "Account ID"
// @Param Body body apireq.ChangePassword true "Request Change Password"
// @Success 200 {string} string "{}"
// @Failure 400 {object} er.AppErrorMsg "{"code":"400400","message":"Wrong parameter format or invalid"}"
// @Failure 401 {object} er.AppErrorMsg "{"code":"400401","message":"Unauthorized"}"
// @Failure 403 {object} er.AppErrorMsg "{"code":"400403","message":"Permission denied"}"
// @Failure 404 {object} er.AppErrorMsg "{"code":"400404","message":"Resource not found"}"
// @Failure 500 {object} er.AppErrorMsg "{"code":"500000","message":"Database unknown error"}"
// @Router /v1/accounts/{account_id}/change-password [put]
func ChangePassword(c *gin.Context) {
	id := c.Param("id")
	sysAccId, err := strconv.Atoi(id)
	if err != nil {
		err = er.NewAppErr(http.StatusBadRequest, er.ErrorParamInvalid, "system account id format error.", err)
		_ = c.Error(err)
		return
	}

	req := apireq.ChangePassword{}
	err = c.BindJSON(&req)
	if err != nil {
		paramErr := er.NewAppErr(http.StatusBadRequest, er.ErrorParamInvalid, err.Error(), err)
		_ = c.Error(paramErr)
		return
	}

	// 參數驗證
	err = valider.Validate.Struct(req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	err = helper.PasswordValidation(req.NewPassword)
	if err != nil {
		_ = c.Error(err)
		return
	}

	// 驗證 jwt user == user_id
	err = tokenLibrary.CheckJWTAccountId(c, sysAccId)
	if err != nil {
		_ = c.Error(err)
		return
	}

	env := api.GetEnv()
	tc := tokenRepo.NewCache(env.RedisCluster)
	sc := sysRepo.NewCache(env.RedisCluster)
	sr := sysRepo.NewRepository(env.Orm, sc)
	srr := sysRoleRepo.NewRepository(env.Orm)
	sar := sysAccRepo.NewRepository(env.Orm)
	sas := sysAccSrv.NewService(sar, sr, srr, tc)

	err = sas.ChangePassword(sysAccId, req.OldPassword, req.NewPassword)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{})
}
