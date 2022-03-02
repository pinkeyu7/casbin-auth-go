package system

import (
	"casbin-auth-go/dto/apireq"
	"casbin-auth-go/dto/apires"
	"casbin-auth-go/dto/model"
	"fmt"
	"math/rand"
)

type Repository interface {
	Insert(req *apireq.AddSystem) error
	Find(listType string, offset, limit int) ([]*apires.System, error)
	FindOne(m *model.System) (*model.System, error)
	Count(listType string) (int, error)
	Update(m *model.System) error
}

type Cache interface {
}

func GenIncNo(tag string) string {
	return fmt.Sprintf("%s%03d", tag, rand.Intn(100000))
}
