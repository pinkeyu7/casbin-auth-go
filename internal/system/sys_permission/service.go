package sys_permission

import (
	"casbin-auth-go/dto/apireq"
	"casbin-auth-go/internal/system/sys_role"
	"casbin-auth-go/internal/system/system"
)

type Service interface {
	AddPermission(req *apireq.AddSysPermission) error
	EditPermission(sysPermId int, req *apireq.EditSysPermission) error
	DeletePermission(sysPermId int, sysRepo system.Repository, sysRoleRepo sys_role.Repository) error
}
