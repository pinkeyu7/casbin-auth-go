package seed

import (
	"casbin-auth-go/dto/model"
	"xorm.io/xorm"
)

func CreateSysAccountRole(engine *xorm.Engine, accId, roleId int) error {
	accountRole := model.SysAccountRole{
		SysAccountId: accId,
		SysRoleId:    roleId,
	}

	_, err := engine.Insert(&accountRole)

	return err
}

func AllSysAccountRole() []Seed {
	return []Seed{
		{
			Name: "Create System Account Role - Root",
			Run: func(engine *xorm.Engine) error {
				err := CreateSysAccountRole(engine, 1, 1)
				return err
			},
		},
		{
			Name: "Create System Account Role - Admin1",
			Run: func(engine *xorm.Engine) error {
				err := CreateSysAccountRole(engine, 3, 2)
				return err
			},
		},
		{
			Name: "Create System Account Role - Admin2",
			Run: func(engine *xorm.Engine) error {
				err := CreateSysAccountRole(engine, 4, 3)
				return err
			},
		},
		{
			Name: "Create System Account Role - Operator",
			Run: func(engine *xorm.Engine) error {
				err := CreateSysAccountRole(engine, 5, 4)
				return err
			},
		},
	}
}
