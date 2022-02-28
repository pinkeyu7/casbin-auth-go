package model

import "time"

type SysAccount struct {
	Id        int       `xorm:"not null pk autoincr INT(11) id" json:"id"`
	SystemId  int       `xorm:"not null INT(11) system_id" json:"system_id"`
	Account   string    `xorm:"not null default '' VARCHAR(64)" json:"account"`
	Phone     string    `xorm:"not null default '' VARCHAR(20)" json:"phone"`
	Email     string    `xorm:"not null default '' VARCHAR(64)" json:"email"`
	Password  string    `xorm:"not null default '' VARCHAR(255)" json:"password"`
	Name      string    `xorm:"not null default '' VARCHAR(64)" json:"name"`
	IsDisable bool      `xorm:"not null is_disable" json:"is_disable"`
	VerifyAt  time.Time `xorm:"DATETIME" json:"verify_at"`
	CreatedAt time.Time `xorm:"not null created DATETIME" json:"created_at"`
	UpdatedAt time.Time `xorm:"not null updated DATETIME" json:"updated_at"`
}
