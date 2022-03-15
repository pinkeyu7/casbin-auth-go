package sys_account

import (
	"casbin-auth-go/dto/apireq"
	"casbin-auth-go/dto/apires"
)

type Service interface {
	ListSysAccount(req *apireq.ListSysAccount) (*apires.ListSysAccount, error)
	GetSysAccount(sysAccId int) (*apires.SysAccount, error)
	AddSysAccount(req *apireq.AddSysAccount) error
	EditSysAccount(sysAccId int, req *apireq.EditSysAccount) error
	ChangePassword(sysAccId int, oldPw, newPw string) error
	ForgotPasswordByAdmin(sysAccId int, req *apireq.ForgotSysAccountPassword) error
}
