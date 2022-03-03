package sys_permission

import (
	"casbin-auth-go/dto/apireq"
	"casbin-auth-go/dto/apires"
	"casbin-auth-go/internal/system/sys_role"
)

type Service interface {
	ListPermission(req *apireq.ListSysPermission) (*apires.ListSysPermission, error)
	AddPermission(req *apireq.AddSysPermission) error
	EditPermission(sysPermId int, req *apireq.EditSysPermission) error
	DeletePermission(sysPermId int, sysRoleRepo sys_role.Repository) error
}
