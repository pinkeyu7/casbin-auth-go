package model

import "time"

type SysAccountRole struct {
	SysAccountId int       `xorm:"not null INT(11)" json:"sys_account_id" `
	SysRoleId    int       `xorm:"not null INT(11)" json:"sys_role_id" `
	CreatedAt    time.Time `xorm:"not null created DATETIME" json:"created_at"`
}
