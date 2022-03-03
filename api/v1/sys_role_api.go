package v1

import (
	"casbin-auth-go/api"
	"casbin-auth-go/dto/apireq"
	sysPermRepo "casbin-auth-go/internal/system/sys_permission/repository"
	sysRoleRepo "casbin-auth-go/internal/system/sys_role/repository"
	sysRoleSrv "casbin-auth-go/internal/system/sys_role/service"
	sysRepo "casbin-auth-go/internal/system/system/repository"
	tokenLibrary "casbin-auth-go/internal/token/library"
	"casbin-auth-go/pkg/er"
	"casbin-auth-go/pkg/valider"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// ListSysRole
// @Summary List System Role 取得角色列表
// @Produce json
// @Accept json
// @Tags Role
// @Security Bearer
// @Param Bearer header string true "JWT Token"
// @Param account_id query int true "Account ID"
// @Param page query string true "Page"
// @Param per_page query string true "Per Page"
// @Param system_id query int false "System ID"
// @Success 200 {object} apires.ListSysRole
// @Failure 400 {object} er.AppErrorMsg "{"code":"400400","message":"Wrong parameter format or invalid"}"
// @Failure 401 {object} er.AppErrorMsg "{"code":"400401","message":"Unauthorized"}"
// @Failure 403 {object} er.AppErrorMsg "{"code":"400403","message":"Permission denied"}"
// @Failure 404 {object} er.AppErrorMsg "{"code":"400404","message":"Resource not found"}"
// @Failure 500 {object} er.AppErrorMsg "{"code":"500000","message":"Database unknown error"}"
// @Router /v1/roles [get]
func ListSysRole(c *gin.Context) {
	req := apireq.ListSysRole{}
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
	spr := sysPermRepo.NewRepository(env.Orm)
	srr := sysRoleRepo.NewRepository(env.Orm)
	srs := sysRoleSrv.NewService(srr, sr, spr)

	res, err := srs.ListSysRole(&req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, res)
}

// AddSysRole
// @Summary Add System Role 新增角色
// @Produce json
// @Accept json
// @Tags Role
// @Security Bearer
// @Param Bearer header string true "JWT Token"
// @Param Body body apireq.AddSysRole true "Request Add System Role"
// @Success 200 {string} string "{}"
// @Failure 400 {object} er.AppErrorMsg "{"code":"400400","message":"Wrong parameter format or invalid"}"
// @Failure 401 {object} er.AppErrorMsg "{"code":"400401","message":"Unauthorized"}"
// @Failure 403 {object} er.AppErrorMsg "{"code":"400403","message":"Permission denied"}"
// @Failure 404 {object} er.AppErrorMsg "{"code":"400404","message":"Resource not found"}"
// @Failure 500 {object} er.AppErrorMsg "{"code":"500000","message":"Database unknown error"}"
// @Router /v1/roles [post]
func AddSysRole(c *gin.Context) {
	req := apireq.AddSysRole{}
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
	spr := sysPermRepo.NewRepository(env.Orm)
	srr := sysRoleRepo.NewRepository(env.Orm)
	srs := sysRoleSrv.NewService(srr, sr, spr)

	err = srs.AddSysRoleWithPermission(&req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{})
}

// EditSysRole
// @Summary Edit System Role 編輯角色
// @Produce json
// @Accept json
// @Tags Role
// @Security Bearer
// @Param Bearer header string true "JWT Token"
// @Param role_id path int true "Role ID"
// @Param Body body apireq.EditSysRole true "Request Edit System Role"
// @Success 200 {string} string "{}"
// @Failure 400 {object} er.AppErrorMsg "{"code":"400400","message":"Wrong parameter format or invalid"}"
// @Failure 401 {object} er.AppErrorMsg "{"code":"400401","message":"Unauthorized"}"
// @Failure 403 {object} er.AppErrorMsg "{"code":"400403","message":"Permission denied"}"
// @Failure 404 {object} er.AppErrorMsg "{"code":"400404","message":"Resource not found"}"
// @Failure 500 {object} er.AppErrorMsg "{"code":"500000","message":"Database unknown error"}"
// @Router /v1/roles/{role_id} [put]
func EditSysRole(c *gin.Context) {
	id := c.Param("id")
	sysRoleId, err := strconv.Atoi(id)
	if err != nil {
		err = er.NewAppErr(http.StatusBadRequest, er.ErrorParamInvalid, "system role id format error.", err)
		_ = c.Error(err)
		return
	}

	req := apireq.EditSysRole{}
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
	spr := sysPermRepo.NewRepository(env.Orm)
	srr := sysRoleRepo.NewRepository(env.Orm)
	srs := sysRoleSrv.NewService(srr, sr, spr)

	err = srs.EditSysRoleWithPermission(sysRoleId, &req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{})
}

// GetSysRoleWithPermission
// @Summary System Role Permission 角色權限詳情
// @Produce json
// @Accept json
// @Tags Role
// @Security Bearer
// @Param Bearer header string true "JWT Token"
// @Param role_id path int true "Role ID"
// @Param account_id query int true "Account ID"
// @Success 200 {object} apires.SysRoleWithPermissionIds
// @Failure 400 {object} er.AppErrorMsg "{"code":"400400","message":"Wrong parameter format or invalid"}"
// @Failure 401 {object} er.AppErrorMsg "{"code":"400401","message":"Unauthorized"}"
// @Failure 403 {object} er.AppErrorMsg "{"code":"400403","message":"Permission denied"}"
// @Failure 404 {object} er.AppErrorMsg "{"code":"400404","message":"Resource not found"}"
// @Failure 500 {object} er.AppErrorMsg "{"code":"500000","message":"Database unknown error"}"
// @Router /v1/roles/{role_id} [get]
func GetSysRoleWithPermission(c *gin.Context) {
	id := c.Param("id")
	sysRoleId, err := strconv.Atoi(id)
	if err != nil {
		err = er.NewAppErr(http.StatusBadRequest, er.ErrorParamInvalid, "system role id format error.", err)
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
	sc := sysRepo.NewCache(env.RedisCluster)
	sr := sysRepo.NewRepository(env.Orm, sc)
	spr := sysPermRepo.NewRepository(env.Orm)
	srr := sysRoleRepo.NewRepository(env.Orm)
	srs := sysRoleSrv.NewService(srr, sr, spr)

	res, err := srs.GetSysRoleWithPermission(sysRoleId)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, res)
}
