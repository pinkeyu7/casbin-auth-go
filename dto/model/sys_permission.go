package model

import "time"

type SysPermission struct {
	Id           int       `xorm:"not null pk autoincr INT(11) id" json:"id"`
	SystemId     int       `xorm:"not null INT(11) system_id" json:"system_id"`
	AllowApiPath string    `xorm:"not null default '' VARCHAR(255) allow_api_path" json:"allow_api_path" `
	Action       string    `xorm:"not null default '' VARCHAR(128) action" json:"action"`
	Slug         string    `xorm:"not null default '' VARCHAR(128) slug" json:"slug"`
	Description  string    `xorm:"not null default '' VARCHAR(255) description" json:"description"`
	IsDisable    bool      `xorm:"not null is_disable" json:"is_disable"`
	CreatedAt    time.Time `xorm:"not null created DATETIME" json:"created_at"`
	UpdatedAt    time.Time `xorm:"not null updated DATETIME" json:"updated_at"`
}
