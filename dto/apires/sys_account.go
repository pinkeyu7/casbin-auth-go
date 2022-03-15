package apires

import (
	"time"

	"gopkg.in/guregu/null.v4"
)

type ListSysAccount struct {
	List        []*ListSysAccountItem `json:"list"`
	Total       int                   `json:"total"`
	CurrentPage int                   `json:"current_page"`
	PerPage     int                   `json:"per_page"`
	NextPage    null.Int              `json:"next_page" swaggertype:"string"`
}

type ListSysAccountItem struct {
	Id        int       `json:"id"`
	SystemId  int       `json:"system_id"`
	Account   string    `json:"account"`
	SysRoleId int       `json:"sys_role_id"`
	Phone     string    `json:"phone"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	IsDisable bool      `json:"is_disable"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type SysAccount struct {
	Id        int       `xorm:"not null pk autoincr INT(11) id" json:"id"`
	SystemId  int       `xorm:"not null INT(11) system_id" json:"system_id"`
	Account   string    `xorm:"not null default '' VARCHAR(64)" json:"account"`
	Phone     string    `xorm:"not null default '' VARCHAR(20)" json:"phone"`
	Email     string    `xorm:"not null default '' VARCHAR(64)" json:"email"`
	Name      string    `xorm:"not null default '' VARCHAR(64)" json:"name"`
	IsDisable bool      `xorm:"not null is_disable" json:"is_disable"`
	VerifyAt  time.Time `xorm:"DATETIME" json:"verify_at"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type SysAccountWithRole struct {
	SysAccount `xorm:"extends"`
	SysRoleId  int `xorm:"sys_role_id" json:"sys_role_id"`
}
