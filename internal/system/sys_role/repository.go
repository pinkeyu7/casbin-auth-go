package sys_role

import "casbin-auth-go/dto/model"

type Repository interface {
	InsertWithPermission(sysRole *model.SysRole, permIds []int) error
	UpdateWithPermission(sysRole *model.SysRole, permIds []int) error
	FindOne(m *model.SysRole) (*model.SysRole, error)
	Find(sysId int, offset, limit int) ([]*model.SysRole, error)
	Count(sysId int) (int, error)
	FindBySysId(sysId int) ([]*model.SysRole, error)
	FindByPermId(permId int) ([]*model.SysRole, error)
}
