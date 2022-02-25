package service

import (
	"casbin-auth-go/dto/apireq"
	"casbin-auth-go/dto/apires"
	"casbin-auth-go/dto/model"
	"casbin-auth-go/internal/system/sys_account"
	"casbin-auth-go/internal/token"
	tokenLibrary "casbin-auth-go/internal/token/library"
	"casbin-auth-go/pkg/er"
	"casbin-auth-go/pkg/helper"
	"net/http"
	"time"
)

type Service struct {
	sysAccRepo sys_account.Repository
	tokenCache token.Cache
}

func NewService(sar sys_account.Repository, tc token.Cache) token.Service {
	return &Service{
		sysAccRepo: sar,
		tokenCache: tc,
	}
}

func (s *Service) GenToken(req *apireq.GetSysAccountToken) (*apires.SysAccountToken, error) {
	// Check Account Exist
	acc, err := s.sysAccRepo.FindOne(&model.SysAccount{Account: req.Account})
	if err != nil {
		findErr := er.NewAppErr(http.StatusInternalServerError, er.UnknownError, "find account error.", err)
		return nil, findErr
	}
	if acc == nil || acc.IsDisable {
		authErr := er.NewAppErr(http.StatusUnauthorized, er.UnauthorizedError, "", nil)
		return nil, authErr
	}

	// Password not matched
	pw := helper.ScryptStr(req.Password)
	if acc.Password != pw {
		authErr := er.NewAppErr(http.StatusUnauthorized, er.UnauthorizedError, "", nil)
		return nil, authErr
	}

	oToken, expiredAt, err := tokenLibrary.GenToken(acc.Id)
	if err != nil {
		tokenErr := er.NewAppErr(http.StatusUnauthorized, er.UnauthorizedError, "", err)
		return nil, tokenErr
	}

	// Set iat
	iat := time.Now().UTC().Unix()
	_ = s.tokenCache.SetTokenIat(acc.Id, iat)

	mapData := map[string]interface{}{}
	mapData["name"] = acc.Name
	mapData["email"] = acc.Email

	res := apires.SysAccountToken{
		Token:     oToken,
		ExpiredAt: expiredAt,
		Data:      mapData,
	}

	return &res, nil
}
