package repository

import (
	"casbin-auth-go/dto/model"
	"casbin-auth-go/internal/system/sys_permission"
	"xorm.io/xorm"
)

type Repository struct {
	orm *xorm.EngineGroup
}

func NewRepository(orm *xorm.EngineGroup) sys_permission.Repository {
	return &Repository{
		orm: orm,
	}
}

func (r *Repository) Insert(m *model.SysPermission) error {
	_, err := r.orm.Insert(m)
	return err
}

func (r *Repository) Update(m *model.SysPermission) error {
	_, err := r.orm.ID(m.Id).Omit("system_id", "allow_api_path", "action", "description").Update(m)
	return err
}

func (r *Repository) Delete(sysPermId int) error {
	session := r.orm.NewSession()
	defer session.Close()

	err := session.Begin()
	if err != nil {
		_ = session.Rollback()
		return err
	}

	// ----------------------- start session -----------------------

	_, err = session.ID(sysPermId).Delete(&model.SysPermission{})
	if err != nil {
		_ = session.Rollback()
		return err
	}

	_, err = session.Where(" sys_permission_id = ? ", sysPermId).Delete(&model.SysRolePermission{})
	if err != nil {
		_ = session.Rollback()
		return err
	}

	// ----------------------- end session -----------------------

	_ = session.Commit()

	return nil
}

func (r *Repository) FindOne(m *model.SysPermission) (*model.SysPermission, error) {
	has, err := r.orm.Get(m)
	if err != nil {
		return nil, err
	}

	if !has {
		return nil, nil
	}

	return m, nil
}

func (r *Repository) FindByIds(permIds []int) ([]*model.SysPermission, error) {
	list := make([]*model.SysPermission, 0)

	err := r.orm.Where("is_disable = ?", 0).In("id", permIds).Find(&list)
	return list, err
}

func (r *Repository) Find(sysId int, offset, limit int) ([]*model.SysPermission, error) {
	var err error
	list := make([]*model.SysPermission, 0)

	if sysId == 0 {
		err = r.orm.Table("sys_permission").Limit(limit, offset).Find(&list)
	} else {
		err = r.orm.Table("sys_permission").Where("system_id = ? ", sysId).Limit(limit, offset).Find(&list)
	}

	return list, err
}

func (r *Repository) Count(sysId int) (int, error) {
	var err error
	var total int64

	if sysId == 0 {
		total, err = r.orm.Table("sys_permission").Count(&model.SysPermission{})
	} else {
		total, err = r.orm.Table("sys_permission").Where("system_id = ? ", sysId).Count(&model.SysPermission{})
	}

	if err != nil {
		return 0, err
	}

	return int(total), nil
}

func (r *Repository) Exist(m *model.SysPermission) (bool, error) {
	return r.orm.Exist(m)
}

func (r *Repository) FindIdsByRole(sysRoleId int) ([]int, error) {
	permIds := make([]int, 0)
	err := r.orm.Table("sys_role_permission").Cols("sys_permission_id").Where("sys_role_id = ? ", sysRoleId).Find(&permIds)
	if err != nil {
		return nil, err
	}
	return permIds, nil
}
