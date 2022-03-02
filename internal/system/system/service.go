package system

import (
	"casbin-auth-go/dto/apireq"
	"casbin-auth-go/dto/apires"
	"casbin-auth-go/dto/model"
)

type Service interface {
	ListSystem(req *apireq.ListSystem) (*apires.ListSystem, error)
	GetSystem(sysId int) (*model.System, error)
	AddSystem(req *apireq.AddSystem) error
	EditSystem(sysId int, req *apireq.EditSystem) error
}
