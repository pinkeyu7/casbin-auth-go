package seed

import (
	"casbin-auth-go/dto/model"
	"xorm.io/xorm"
)

func CreateSysRole(engine *xorm.Engine, name string, systemId int) error {
	role := model.SysRole{
		Sort:        0,
		SystemId:    systemId,
		Name:        name,
		DisplayName: name,
	}

	_, err := engine.Insert(&role)

	return err
}

func AllSysRole() []Seed {
	return []Seed{
		{
			Name: "Create System Role - Admin Permission - Root",
			Run: func(engine *xorm.Engine) error {
				err := CreateSysRole(engine, "Root", 1)
				return err
			},
		},
		{
			Name: "Create System Role - Address Book System - Admin1",
			Run: func(engine *xorm.Engine) error {
				err := CreateSysRole(engine, "Admin1", 2)
				return err
			},
		},
		{
			Name: "Create System Role - Address Book System - Admin2",
			Run: func(engine *xorm.Engine) error {
				err := CreateSysRole(engine, "Admin2", 2)
				return err
			},
		},
		{
			Name: "Create System Role - Address Book System - Operator",
			Run: func(engine *xorm.Engine) error {
				err := CreateSysRole(engine, "Operator", 2)
				return err
			},
		},
	}
}
