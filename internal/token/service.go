package token

import (
	"casbin-auth-go/dto/apireq"
	"casbin-auth-go/dto/apires"
)

type Service interface {
	GenToken(req *apireq.GetSysAccountToken) (*apires.SysAccountToken, error)
}
