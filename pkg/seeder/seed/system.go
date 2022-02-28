package seed

import (
	"casbin-auth-go/dto/model"
	"github.com/brianvoe/gofakeit/v4"
	"xorm.io/xorm"
)

// --------------------------------------------- System Seeds ---------------------------------------------
func CreateSystem(engine *xorm.Engine, name, systemType, tag, email, address, tel, uuid string, quota int, ipAddress, macAddress string) error {
	s := model.System{
		Name:          name,
		Tag:           tag,
		Email:         email,
		Address:       address,
		Tel:           tel,
		Uuid:          uuid,
		Quota:         quota,
		IpAddress:     ipAddress,
		MacAddress:    macAddress,
		Principal:     gofakeit.Name(),
		Salesman:      gofakeit.Name(),
		SalesmanPhone: gofakeit.Phone(),
		SystemType:    systemType,
	}
	_, err := engine.Insert(&s)
	return err
}

func AllSystem() []Seed {
	return []Seed{
		{
			Name: "Create System - Admin Permission",
			Run: func(engine *xorm.Engine) error {
				err := CreateSystem(engine, "Admin Permission", "AP", "AP", gofakeit.Name(),
					gofakeit.City(), gofakeit.Phone(), gofakeit.Animal(), 0, "*", "*")
				return err
			},
		},
		{
			Name: "Create System - Address Book System",
			Run: func(engine *xorm.Engine) error {
				err := CreateSystem(engine, "Address Book System", "ABS", "ABS", gofakeit.Name(),
					gofakeit.City(), gofakeit.Phone(), gofakeit.Animal(), 10000, "*", "*")
				return err
			},
		},
	}
}
