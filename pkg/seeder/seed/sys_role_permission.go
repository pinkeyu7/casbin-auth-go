package seed

import (
	"casbin-auth-go/dto/model"
	"xorm.io/xorm"
)

func CreateSysRolePermission(engine *xorm.Engine, roleId, permissionId int) error {
	rolePermission := model.SysRolePermission{
		SysRoleId:       roleId,
		SysPermissionId: permissionId,
	}

	_, err := engine.Insert(&rolePermission)

	return err
}

func AllSysRolePermission() []Seed {
	return []Seed{
		{
			Name: "Create System Role Permission 1 for root",
			Run: func(engine *xorm.Engine) error {
				err := CreateSysRolePermission(engine, 1, 1)
				return err
			},
		},
		{
			Name: "Create System Role Permission 2 for root",
			Run: func(engine *xorm.Engine) error {
				err := CreateSysRolePermission(engine, 1, 2)
				return err
			},
		},
		{
			Name: "Create System Role Permission 3 for admin1",
			Run: func(engine *xorm.Engine) error {
				err := CreateSysRolePermission(engine, 2, 3)
				return err
			},
		},
		{
			Name: "Create System Role Permission 4 for admin1",
			Run: func(engine *xorm.Engine) error {
				err := CreateSysRolePermission(engine, 2, 4)
				return err
			},
		},
		{
			Name: "Create System Role Permission 5 for admin1",
			Run: func(engine *xorm.Engine) error {
				err := CreateSysRolePermission(engine, 2, 5)
				return err
			},
		},
		{
			Name: "Create System Role Permission 6 for admin1",
			Run: func(engine *xorm.Engine) error {
				err := CreateSysRolePermission(engine, 2, 6)
				return err
			},
		},
	}
}
