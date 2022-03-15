package v1

import (
	"casbin-auth-go/api"
	"casbin-auth-go/dto/apireq"
	sysAccRepo "casbin-auth-go/internal/system/sys_account/repository"
	tokenRepo "casbin-auth-go/internal/token/repository"
	tokenSrv "casbin-auth-go/internal/token/service"
	"casbin-auth-go/pkg/er"
	"casbin-auth-go/pkg/valider"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetToken
// @Summary Get Token 取得 token
// @Produce json
// @Accept json
// @Tags Token
// @Param Body body apireq.GetSysAccountToken true "Request Get Sys Account Token"
// @Success 200 {object} apires.SysAccountToken
// @Failure 400 {object} er.AppErrorMsg "{"code":"400400","message":"Wrong parameter format or invalid"}"
// @Failure 401 {object} er.AppErrorMsg "{"code":"400401","message":"Unauthorized"}"
// @Failure 403 {object} er.AppErrorMsg "{"code":"400403","message":"Permission denied"}"
// @Failure 404 {object} er.AppErrorMsg "{"code":"400404","message":"Resource not found"}"
// @Failure 500 {object} er.AppErrorMsg "{"code":"500000","message":"Database unknown error"}"
// @Router /v1/token [post]
func GetToken(c *gin.Context) {
	req := apireq.GetSysAccountToken{}
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
	sar := sysAccRepo.NewRepository(env.Orm)
	tc := tokenRepo.NewCache(env.RedisCluster)
	ts := tokenSrv.NewService(sar, tc)
	res, err := ts.GenToken(&req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, res)
}
