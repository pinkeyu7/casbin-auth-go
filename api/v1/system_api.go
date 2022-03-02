package v1

import (
	"casbin-auth-go/api"
	"casbin-auth-go/dto/apireq"
	sysRepo "casbin-auth-go/internal/system/system/repository"
	sysSrv "casbin-auth-go/internal/system/system/service"
	"casbin-auth-go/pkg/er"
	"casbin-auth-go/pkg/valider"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// AddSystem
// @Summary Add System 新增系統
// @Produce json
// @Accept json
// @Tags System
// @Security Bearer
// @Param Bearer header string true "JWT Token"
// @Param Body body apireq.AddSystem true "Request Add System"
// @Success 200 {string} string "{}"
// @Failure 400 {object} er.AppErrorMsg "{"code":"400400","message":"Wrong parameter format or invalid"}"
// @Failure 401 {object} er.AppErrorMsg "{"code":"400401","message":"Unauthorized"}"
// @Failure 403 {object} er.AppErrorMsg "{"code":"400403","message":"Permission denied"}"
// @Failure 404 {object} er.AppErrorMsg "{"code":"400404","message":"Resource not found"}"
// @Failure 500 {object} er.AppErrorMsg "{"code":"500000","message":"Database unknown error"}"
// @Router /v1/systems [post]
func AddSystem(c *gin.Context) {
	req := apireq.AddSystem{}
	err := c.BindJSON(&req)
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

	env := api.GetEnv()
	sc := sysRepo.NewCache(env.RedisCluster)
	sr := sysRepo.NewRepository(env.Orm, sc)
	ss := sysSrv.NewService(sr)
	err = ss.AddSystem(&req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{})
}

// ListSystem
// @Summary List System 系統列表
// @Produce json
// @Accept json
// @Tags System
// @Security Bearer
// @Param Bearer header string true "JWT Token"
// @Param page query string true "Page"
// @Param per_page query string true "Per Page"
// @Success 200 {object} apires.ListSystem
// @Failure 400 {object} er.AppErrorMsg "{"code":"400400","message":"Wrong parameter format or invalid"}"
// @Failure 401 {object} er.AppErrorMsg "{"code":"400401","message":"Unauthorized"}"
// @Failure 403 {object} er.AppErrorMsg "{"code":"400403","message":"Permission denied"}"
// @Failure 404 {object} er.AppErrorMsg "{"code":"400404","message":"Resource not found"}"
// @Failure 500 {object} er.AppErrorMsg "{"code":"500000","message":"Database unknown error"}"
// @Router /v1/systems [get]
func ListSystem(c *gin.Context) {
	req := apireq.ListSystem{}
	err := c.Bind(&req)
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

	env := api.GetEnv()
	sc := sysRepo.NewCache(env.RedisCluster)
	sr := sysRepo.NewRepository(env.Orm, sc)
	ss := sysSrv.NewService(sr)
	res, err := ss.ListSystem(&req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, res)
}

// GetSystem
// @Summary Get System 取得系統
// @Produce json
// @Accept json
// @Tags System
// @Security Bearer
// @Param Bearer header string true "JWT Token"
// @Param system_id path int true "系統id e.g. 11"
// @Success 200 {object} model.System
// @Failure 400 {object} er.AppErrorMsg "{"code":"400400","message":"Wrong parameter format or invalid"}"
// @Failure 401 {object} er.AppErrorMsg "{"code":"400401","message":"Unauthorized"}"
// @Failure 403 {object} er.AppErrorMsg "{"code":"400403","message":"Permission denied"}"
// @Failure 404 {object} er.AppErrorMsg "{"code":"400404","message":"Resource not found"}"
// @Failure 500 {object} er.AppErrorMsg "{"code":"500000","message":"Database unknown error"}"
// @Router /v1/systems/{system_id} [get]
func GetSystem(c *gin.Context) {
	sysIdStr := c.Param("id")
	sysId, err := strconv.Atoi(sysIdStr)
	if err != nil {
		paramErr := er.NewAppErr(http.StatusBadRequest, er.ErrorParamInvalid, "system id format error.", err)
		_ = c.Error(paramErr)
		return
	}

	env := api.GetEnv()
	sc := sysRepo.NewCache(env.RedisCluster)
	sr := sysRepo.NewRepository(env.Orm, sc)
	ss := sysSrv.NewService(sr)
	res, err := ss.GetSystem(sysId)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, res)
}

// EditSystem
// @Summary Edit System 編輯系統
// @Produce json
// @Accept json
// @Tags System
// @Security Bearer
// @Param Bearer header string true "JWT Token"
// @Param system_id path int true "系統id e.g. 11"
// @Param Body body apireq.EditSystem true "Request Edit System"
// @Success 200 {string} string "{}"
// @Failure 400 {object} er.AppErrorMsg "{"code":"400400","message":"Wrong parameter format or invalid"}"
// @Failure 401 {object} er.AppErrorMsg "{"code":"400401","message":"Unauthorized"}"
// @Failure 403 {object} er.AppErrorMsg "{"code":"400403","message":"Permission denied"}"
// @Failure 404 {object} er.AppErrorMsg "{"code":"400404","message":"Resource not found"}"
// @Failure 500 {object} er.AppErrorMsg "{"code":"500000","message":"Database unknown error"}"
// @Router /v1/systems/{system_id} [put]
func EditSystem(c *gin.Context) {
	sysIdStr := c.Param("id")
	sysId, err := strconv.Atoi(sysIdStr)
	if err != nil {
		paramErr := er.NewAppErr(http.StatusBadRequest, er.ErrorParamInvalid, "system id format error.", err)
		_ = c.Error(paramErr)
		return
	}

	req := apireq.EditSystem{}
	err = c.BindJSON(&req)
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

	env := api.GetEnv()
	sc := sysRepo.NewCache(env.RedisCluster)
	sr := sysRepo.NewRepository(env.Orm, sc)
	ss := sysSrv.NewService(sr)
	err = ss.EditSystem(sysId, &req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{})
}
