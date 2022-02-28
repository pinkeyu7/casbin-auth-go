package model

import "time"

type SysWhiteList struct {
	Id        int       `xorm:"not null pk autoincr INT(11) id" json:"id"`
	SystemId  int       `xorm:"not null INT(11) system_id" json:"system_id"`
	Phone     string    `xorm:"not null default '' VARCHAR(20) phone" json:"phone"`
	Email     string    `xorm:"not null default '' VARCHAR(64) email" json:"email"`
	Type      string    `xorm:"not null enum('email','phone')" json:"type"`
	IsDisable bool      `xorm:"is_disable" json:"is_disable"`
	CreatedAt time.Time `xorm:"not null created DATETIME" json:"created_at"`
	UpdatedAt time.Time `xorm:"not null updated DATETIME" json:"updated_at"`
}
