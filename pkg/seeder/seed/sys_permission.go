package seed

import (
	"casbin-auth-go/dto/model"
	"xorm.io/xorm"
)

func CreateSysPermission(engine *xorm.Engine, sysId int, apiPath, action, slug, desc string, isDisable bool) error {
	s := model.SysPermission{
		SystemId:     sysId,
		AllowApiPath: apiPath,
		Action:       action,
		Slug:         slug,
		Description:  desc,
		IsDisable:    isDisable,
	}

	_, err := engine.Insert(&s)

	return err
}

func AllSysPermission() []Seed {
	return []Seed{
		{
			Name: "Create System Permission - Admin Permission - 1",
			Run: func(engine *xorm.Engine) error {
				err := CreateSysPermission(engine, 1, "/v1/admin/system-permissions", "POST",
					"v1-admin-system-permission-add", "新增系統權限", false)
				return err
			},
		},
		{
			Name: "Create System Permission - Admin Permission - 2",
			Run: func(engine *xorm.Engine) error {
				err := CreateSysPermission(engine, 1, "/v1/admin/system-permissions/:id", "PUT",
					"v1-admin-system-permission-edit", "編輯系統權限", false)
				return err
			},
		},
		{
			Name: "Create System Permission - Address Book System - 1",
			Run: func(engine *xorm.Engine) error {
				err := CreateSysPermission(engine, 2, "/v1/contacts", "GET",
					"v1-contact-list", "取得聯絡人列表", false)
				return err
			},
		},
		{
			Name: "Create System Permission - Address Book System - 2",
			Run: func(engine *xorm.Engine) error {
				err := CreateSysPermission(engine, 2, "/v1/contacts/:id", "GET",
					"v1-contact-id", "取得聯絡人", false)
				return err
			},
		},
		{
			Name: "Create System Permission - Address Book System - 3",
			Run: func(engine *xorm.Engine) error {
				err := CreateSysPermission(engine, 2, "/v1/contacts", "POST",
					"v1-contact-add", "新增聯絡人", false)
				return err
			},
		},
		{
			Name: "Create System Permission - Address Book System - 4",
			Run: func(engine *xorm.Engine) error {
				err := CreateSysPermission(engine, 2, "/v1/contacts/:id", "PUT",
					"v1-contact-edit", "編輯聯絡人", false)
				return err
			},
		},
	}
}
