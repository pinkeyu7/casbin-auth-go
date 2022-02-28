package seed

import (
	"casbin-auth-go/pkg/casbin"
	"xorm.io/xorm"
)

func CreateCasbin(accountStr, roleStr string) error {
	e := casbin.Init()

	_, _ = e.AddRoleForUser(accountStr, roleStr)
	_, err := e.AddPolicy(roleStr, ".*", ".*")

	return err
}

func AllCabin() []Seed {
	return []Seed{
		{
			Name: "Create Casbin Root",
			Run: func(engine *xorm.Engine) error {
				err := CreateCasbin("1", "1")
				return err
			},
		},
	}
}
