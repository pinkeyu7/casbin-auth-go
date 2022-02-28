package seed

import (
	"casbin-auth-go/dto/model"
	"casbin-auth-go/pkg/helper"
	"github.com/brianvoe/gofakeit/v4"
	"time"
	"xorm.io/xorm"
)

func CreateSysAccount(engine *xorm.Engine, systemId int, account, email, password, name string) error {
	sa := model.SysAccount{
		SystemId: systemId,
		Account:  account,
		Phone:    gofakeit.Phone(),
		Email:    email,
		Password: helper.ScryptStr(password),
		Name:     name,
		VerifyAt: time.Now(),
	}

	_, err := engine.Insert(&sa)

	return err
}

func AllSysAccount() []Seed {
	return []Seed{
		{
			Name: "Create System Account - Admin Permission - 1",
			Run: func(engine *xorm.Engine) error {
				err := CreateSysAccount(engine, 1, "pinke", "pinke.yu7@gmail.com", "123456", "Pinke")
				return err
			},
		},
		{
			Name: "Create System Account - Admin Permission - 2",
			Run: func(engine *xorm.Engine) error {
				err := CreateSysAccount(engine, 1, "test_account", "test_account@testmail.com", "123456", gofakeit.Name())
				return err
			},
		},
		{
			Name: "Create System Account - Address Book System - 1",
			Run: func(engine *xorm.Engine) error {
				err := CreateSysAccount(engine, 2, "admin1", "admin1@address-book.com", "123456", gofakeit.Name())
				return err
			},
		},
		{
			Name: "Create System Account - Address Book System - 2",
			Run: func(engine *xorm.Engine) error {
				err := CreateSysAccount(engine, 2, "admin2", "admin2@address-book.com", "123456", gofakeit.Name())
				return err
			},
		},
		{
			Name: "Create System Account - Address Book System - 3",
			Run: func(engine *xorm.Engine) error {
				err := CreateSysAccount(engine, 2, "operator", "operator@address-book.com", "123456", gofakeit.Name())
				return err
			},
		},
	}
}
