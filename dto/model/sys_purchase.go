package model

import "time"

type SysPurchase struct {
	Id         int       `xorm:"not null pk autoincr INT(11) id" json:"id"`
	SystemId   int       `xorm:"not null INT(11) system_id" json:"system_id"`
	TrackingNo string    `xorm:"not null default '' VARCHAR(32)" json:"tracking_no"`
	Quota      int       `xorm:"not null default 0 INT(11)" json:"quota"`
	Salesman   string    `xorm:"not null default '' comment('業務') VARCHAR(64)" json:"salesman"`
	AppliedAt  time.Time `xorm:"not null comment('申請日期') DATETIME" json:"applied_at"`
	CreatedAt  time.Time `xorm:"not null created DATETIME" json:"created_at"`
}
