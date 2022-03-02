package system

import (
	"casbin-auth-go/dto/apireq"
	"casbin-auth-go/dto/apires"
)

type Service interface {
	ListSystem(req *apireq.ListSystem) (*apires.ListSystem, error)
	AddSystem(req *apireq.AddSystem) error
	EditSystem(sysId int, req *apireq.EditSystem) error
}
