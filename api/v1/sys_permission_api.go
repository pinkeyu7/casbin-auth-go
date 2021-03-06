package v1

import (
	"casbin-auth-go/api"
	"casbin-auth-go/dto/apireq"
	sysPermRepo "casbin-auth-go/internal/system/sys_permission/repository"
	sysPermSrv "casbin-auth-go/internal/system/sys_permission/service"
	sysRoleRepo "casbin-auth-go/internal/system/sys_role/repository"
	sysRepo "casbin-auth-go/internal/system/system/repository"
	tokenLibrary "casbin-auth-go/internal/token/library"
	"casbin-auth-go/pkg/er"
	"casbin-auth-go/pkg/valider"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// ListSysPermission
// @Summary List System Permission 取得系統權限列表
// @Produce json
// @Accept json
// @Tags Permission
// @Security Bearer
// @Param Bearer header string true "JWT Token"
// @Param account_id query int true "Account ID"
// @Param page query string true "Page"
// @Param per_page query string true "Per Page"
// @Param system_id query int false "System ID"
// @Success 200 {object} apires.ListSysPermission
// @Failure 400 {object} er.AppErrorMsg "{"code":"400400","message":"Wrong parameter format or invalid"}"
// @Failure 401 {object} er.AppErrorMsg "{"code":"400401","message":"Unauthorized"}"
// @Failure 403 {object} er.AppErrorMsg "{"code":"400403","message":"Permission denied"}"
// @Failure 404 {object} er.AppErrorMsg "{"code":"400404","message":"Resource not found"}"
// @Failure 500 {object} er.AppErrorMsg "{"code":"500000","message":"Database unknown error"}"
// @Router /v1/permissions [get]
func ListSysPermission(c *gin.Context) {
	req := apireq.ListSysPermission{}
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
	sc := sysRepo.NewCache(env.RedisCluster)
	sr := sysRepo.NewRepository(env.Orm, sc)
	sperr := sysPermRepo.NewRepository(env.Orm)
	spers := sysPermSrv.NewService(sperr, sr)

	res, err := spers.ListSysPermission(&req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, res)
}

// AddSysPermission
// @Summary Add System Permission 新增系統權限
// @Produce json
// @Accept json
// @Tags Permission
// @Security Bearer
// @Param Bearer header string true "JWT Token"
// @Param Body body apireq.AddSysPermission true "Request Add System Permission"
// @Success 200 {string} string "{}"
// @Failure 400 {object} er.AppErrorMsg "{"code":"400400","message":"Wrong parameter format or invalid"}"
// @Failure 401 {object} er.AppErrorMsg "{"code":"400401","message":"Unauthorized"}"
// @Failure 403 {object} er.AppErrorMsg "{"code":"400403","message":"Permission denied"}"
// @Failure 404 {object} er.AppErrorMsg "{"code":"400404","message":"Resource not found"}"
// @Failure 500 {object} er.AppErrorMsg "{"code":"500000","message":"Database unknown error"}"
// @Router /v1/permissions [post]
func AddSysPermission(c *gin.Context) {
	req := apireq.AddSysPermission{}
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
	sc := sysRepo.NewCache(env.RedisCluster)
	sr := sysRepo.NewRepository(env.Orm, sc)
	sperr := sysPermRepo.NewRepository(env.Orm)
	spers := sysPermSrv.NewService(sperr, sr)

	err = spers.AddSysPermission(&req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{})
}

// EditSysPermission
// @Summary Edit System Permission 編輯系統權限
// @Produce json
// @Accept json
// @Tags Permission
// @Security Bearer
// @Param Bearer header string true "JWT Token"
// @Param permission_id path int true "Permission ID"
// @Param Body body apireq.EditSysPermission true "Request Edit System Permission"
// @Success 200 {string} string "{}"
// @Failure 400 {object} er.AppErrorMsg "{"code":"400400","message":"Wrong parameter format or invalid"}"
// @Failure 401 {object} er.AppErrorMsg "{"code":"400401","message":"Unauthorized"}"
// @Failure 403 {object} er.AppErrorMsg "{"code":"400403","message":"Permission denied"}"
// @Failure 404 {object} er.AppErrorMsg "{"code":"400404","message":"Resource not found"}"
// @Failure 500 {object} er.AppErrorMsg "{"code":"500000","message":"Database unknown error"}"
// @Router /v1/permissions/{permission_id} [put]
func EditSysPermission(c *gin.Context) {
	id := c.Param("id")
	permId, err := strconv.Atoi(id)
	if err != nil {
		err = er.NewAppErr(http.StatusBadRequest, er.ErrorParamInvalid, "system permission id format error.", err)
		_ = c.Error(err)
		return
	}

	req := apireq.EditSysPermission{}
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
	sc := sysRepo.NewCache(env.RedisCluster)
	sr := sysRepo.NewRepository(env.Orm, sc)
	sperr := sysPermRepo.NewRepository(env.Orm)
	spers := sysPermSrv.NewService(sperr, sr)

	err = spers.EditSysPermission(permId, &req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{})
}

// DeleteSysPermission
// @Summary Delete System Permission 刪除系統權限
// @Produce json
// @Accept json
// @Tags Permission
// @Security Bearer
// @Param Bearer header string true "JWT Token"
// @Param permission_id path int true "Permission ID"
// @Param Body body apireq.DeleteSysPermission true "Request Delete System Permission"
// @Success 200 {string} string "{}"
// @Failure 400 {object} er.AppErrorMsg "{"code":"400400","message":"Wrong parameter format or invalid"}"
// @Failure 401 {object} er.AppErrorMsg "{"code":"400401","message":"Unauthorized"}"
// @Failure 403 {object} er.AppErrorMsg "{"code":"400403","message":"Permission denied"}"
// @Failure 404 {object} er.AppErrorMsg "{"code":"400404","message":"Resource not found"}"
// @Failure 500 {object} er.AppErrorMsg "{"code":"500000","message":"Database unknown error"}"
// @Router /v1/permissions/{permission_id} [delete]
func DeleteSysPermission(c *gin.Context) {
	id := c.Param("id")
	permId, err := strconv.Atoi(id)
	if err != nil {
		err = er.NewAppErr(http.StatusBadRequest, er.ErrorParamInvalid, "system permission id format error.", err)
		_ = c.Error(err)
		return
	}

	req := apireq.DeleteSysPermission{}
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
	sc := sysRepo.NewCache(env.RedisCluster)
	sr := sysRepo.NewRepository(env.Orm, sc)
	srr := sysRoleRepo.NewRepository(env.Orm)
	sperr := sysPermRepo.NewRepository(env.Orm)
	spers := sysPermSrv.NewService(sperr, sr)

	err = spers.DeleteSysPermission(permId, srr)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{})
}
