package sys_permission

import (
	"casbin-auth-go/dto/apireq"
	"casbin-auth-go/dto/apires"
	"casbin-auth-go/internal/system/sys_role"
)

type Service interface {
	ListSysPermission(req *apireq.ListSysPermission) (*apires.ListSysPermission, error)
	AddSysPermission(req *apireq.AddSysPermission) error
	EditSysPermission(sysPermId int, req *apireq.EditSysPermission) error
	DeleteSysPermission(sysPermId int, sysRoleRepo sys_role.Repository) error
}
