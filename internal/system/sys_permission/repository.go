package sys_permission

import "casbin-auth-go/dto/model"

type Repository interface {
	Insert(m *model.SysPermission) error
	Update(m *model.SysPermission) error
	Delete(sysPermId int) error
	FindOne(m *model.SysPermission) (*model.SysPermission, error)
	FindByIds(permIds []int) ([]*model.SysPermission, error)
	Find(sysId int, offset, limit int) ([]*model.SysPermission, error)
	Count(sysId int) (int, error)
	Exist(m *model.SysPermission) (bool, error)
	FindIdsByRole(sysRoleId int) ([]int, error)
}
