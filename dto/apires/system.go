package apires

import (
	"gopkg.in/guregu/null.v4"
)

type ListSystem struct {
	List        []*System `json:"list"`
	Total       int       `json:"total"`
	CurrentPage int       `json:"current_page"`
	PerPage     int       `json:"per_page"`
	NextPage    null.Int  `json:"next_page" swaggertype:"string"`
}

type System struct {
	Id         int    `xorm:"not null pk autoincr INT(11) id" json:"id"`
	Name       string `xorm:"not null default '' VARCHAR(128) name" json:"name"`
	SystemType string `xorm:"not null default '' VARCHAR(20) system_type" json:"system_type"`
	Tag        string `xorm:"not null default '' VARCHAR(20) tag" json:"tag"`
	Email      string `xorm:"not null default '' VARCHAR(64) email" json:"email"`
	Tel        string `xorm:"not null default '' VARCHAR(20) tel" json:"tel"`
	Uuid       string `xorm:"not null default '' VARCHAR(16) uuid" json:"uuid"`
	Quota      int    `xorm:"not null default 0 INT(11) quota" json:"quota"`
	IsDisable  bool   `xorm:"not null is_disable" json:"is_disable"`
}
