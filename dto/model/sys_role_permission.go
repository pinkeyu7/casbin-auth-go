package model

import "time"

type SysRolePermission struct {
	SysRoleId       int       `xorm:"not null INT(11) sys_role_id" json:"sys_role_id"`
	SysPermissionId int       `xorm:"not null INT(11) sys_permission_id" json:"sys_permission_id"`
	CreatedAt       time.Time `xorm:"not null created DATETIME" json:"created_at"`
}
