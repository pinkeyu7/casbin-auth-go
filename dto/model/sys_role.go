package model

import "time"

type SysRole struct {
	Id          int       `xorm:"not null pk autoincr INT(11) id" json:"id"`
	Sort        int       `xorm:"not null INT(11) sort" json:"sort"`
	SystemId    int       `xorm:"not null INT(11) system_id" json:"system_id"`
	Name        string    `xorm:"not null default '' VARCHAR(64) name" json:"name"`
	DisplayName string    `xorm:"not null default '' VARCHAR(64) display_name" json:"display_name"`
	IsDisable   bool      `xorm:"not null is_disable" json:"is_disable"`
	CreatedAt   time.Time `xorm:"not null created DATETIME" json:"created_at"`
	UpdatedAt   time.Time `xorm:"not null updated DATETIME" json:"updated_at"`
}
