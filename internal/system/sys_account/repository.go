package sys_account

import (
	"casbin-auth-go/dto/apires"
	"casbin-auth-go/dto/model"
)

type Repository interface {
	InsertWithRole(m *model.SysAccount, sysRoleId int) (*model.SysAccount, error)
	UpdateWithRole(m *model.SysAccount, sysRoleId int) error
	UpdatePassword(sysAccId int, newPassword string) error
	FindOne(m *model.SysAccount) (*model.SysAccount, error)
	FindOneWithRole(sysAccId int) (*apires.SysAccountWithRole, error)
	Find(sysId, offset, limit int) ([]*apires.ListSysAccountItem, error)
	Count(sysId int) (int, error)
}
