package model

import "time"

type System struct {
	Id            int       `xorm:"not null pk autoincr INT(11) id" json:"id"`
	Name          string    `xorm:"not null default '' VARCHAR(128) name" json:"name"`
	SystemType    string    `xorm:"not null default '' VARCHAR(20) system_type" json:"system_type"`
	Tag           string    `xorm:"not null default '' VARCHAR(20) tag" json:"tag"`
	Email         string    `xorm:"not null default '' VARCHAR(64) email" json:"email"`
	Address       string    `xorm:"not null default '' VARCHAR(255) address" json:"address"`
	Tel           string    `xorm:"not null default '' VARCHAR(20) tel" json:"tel"`
	Uuid          string    `xorm:"not null default '' VARCHAR(16) uuid" json:"uuid"`
	Quota         int       `xorm:"not null default 0 INT(11) quota" json:"quota"`
	IpAddress     string    `xorm:"not null default '' VARCHAR(191) ip_address" json:"ip_address"`
	MacAddress    string    `xorm:"not null default '' VARCHAR(512) mac_address" json:"mac_address"`
	IsDisable     bool      `xorm:"not null is_disable" json:"is_disable"`
	Principal     string    `xorm:"not null default '' comment('負責人') VARCHAR(64)" json:"principal"`
	Salesman      string    `xorm:"not null default '' comment('業務') VARCHAR(64)" json:"salesman"`
	SalesmanPhone string    `xorm:"not null default '' comment('業務電話') VARCHAR(20)" json:"salesman_phone"`
	CreatedAt     time.Time `xorm:"not null created DATETIME" json:"created_at"`
	UpdatedAt     time.Time `xorm:"not null updated DATETIME" json:"updated_at"`
}
