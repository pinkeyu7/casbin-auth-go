package sys_role

import (
	"casbin-auth-go/dto/apireq"
	"casbin-auth-go/dto/apires"
)

type Service interface {
	AddSysRoleWithPermission(req *apireq.AddSysRole) error
	EditSysRoleWithPermission(sysRoleId int, req *apireq.EditSysRole) error
	GetSysRoleWithPermission(sysRoleId int) (*apires.SysRoleWithPermissionIds, error)
}
