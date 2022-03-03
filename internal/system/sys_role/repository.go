package sys_role

import "casbin-auth-go/dto/model"

type Repository interface {
	InsertWithPermission(m *model.SysRole, permIds []int) (*model.SysRole, error)
	UpdateWithPermission(m *model.SysRole, permIds []int) error
	FindOne(m *model.SysRole) (*model.SysRole, error)
	Find(sysId, offset, limit int) ([]*model.SysRole, error)
	Count(sysId int) (int, error)
	FindBySysId(sysId int) ([]*model.SysRole, error)
	FindByPermId(permId int) ([]*model.SysRole, error)
}
