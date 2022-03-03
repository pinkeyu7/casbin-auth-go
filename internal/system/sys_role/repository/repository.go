package repository

import (
	"casbin-auth-go/dto/model"
	"casbin-auth-go/internal/system/sys_role"
	"xorm.io/xorm"
)

type Repository struct {
	orm *xorm.EngineGroup
}

func NewRepository(orm *xorm.EngineGroup) sys_role.Repository {
	return &Repository{
		orm: orm,
	}
}

func (r *Repository) InsertWithPermission(sysRole *model.SysRole, permIds []int) error {
	session := r.orm.NewSession()
	defer session.Close()

	err := session.Begin()
	if err != nil {
		_ = session.Rollback()
		return err
	}

	// ----------------------- start session -----------------------

	_, err = session.Insert(sysRole)
	if err != nil {
		_ = session.Rollback()
		return err
	}

	permLen := len(permIds)
	if permLen > 0 {
		var rolePerm = make([]*model.SysRolePermission, len(permIds))
		for i := range permIds {
			rolePerm[i] = &model.SysRolePermission{
				SysRoleId:       sysRole.Id,
				SysPermissionId: permIds[i],
			}
		}

		_, err = session.Insert(rolePerm)
		if err != nil {
			_ = session.Rollback()
			return err
		}
	}

	// ----------------------- end session -----------------------

	_ = session.Commit()

	return err
}

func (r *Repository) UpdateWithPermission(sysRole *model.SysRole, permIds []int) error {
	session := r.orm.NewSession()
	defer session.Close()

	err := session.Begin()
	if err != nil {
		_ = session.Rollback()
		return err
	}

	// ----------------------- start session -----------------------

	_, err = session.ID(sysRole.Id).Cols("is_disable", "name", "display_name").Update(sysRole)
	if err != nil {
		_ = session.Rollback()
		return err
	}

	permLen := len(permIds)
	if permLen > 0 {
		_, err = session.Where("sys_role_id = ?", sysRole.Id).Delete(&model.SysRolePermission{})
		if err != nil {
			_ = session.Rollback()
			return err
		}

		var rolePerm = make([]*model.SysRolePermission, permLen)
		for i := range permIds {
			rolePerm[i] = &model.SysRolePermission{
				SysRoleId:       sysRole.Id,
				SysPermissionId: permIds[i],
			}
		}

		_, err = session.Insert(rolePerm)
		if err != nil {
			_ = session.Rollback()
			return err
		}
	}

	// ----------------------- end session -----------------------

	_ = session.Commit()

	return nil
}

func (r *Repository) FindOne(m *model.SysRole) (*model.SysRole, error) {
	has, err := r.orm.Get(m)
	if err != nil {
		return nil, err
	}

	if !has {
		return nil, nil
	}

	return m, nil
}

func (r *Repository) Find(sysId int, offset, limit int) ([]*model.SysRole, error) {
	var err error
	list := make([]*model.SysRole, 0)

	if sysId == 0 {
		err = r.orm.Table("sys_role").Limit(limit, offset).Find(&list)
	} else {
		err = r.orm.Table("sys_role").Where("system_id = ? ", sysId).Limit(limit, offset).Find(&list)
	}

	return list, err
}

func (r *Repository) Count(sysId int) (int, error) {
	var err error
	var count int64

	if sysId == 0 {
		count, err = r.orm.Table("sys_role").Count(&model.SysRole{})
	} else {
		count, err = r.orm.Table("sys_role").Where("system_id = ? ", sysId).Count(&model.SysRole{})
	}

	if err != nil {
		return 0, err
	}

	return int(count), nil
}

func (r *Repository) FindBySysId(sysId int) ([]*model.SysRole, error) {
	var err error
	list := make([]*model.SysRole, 0)

	if sysId == 0 {
		err = r.orm.Find(&list)
	} else {
		err = r.orm.Where("system_id = ?", sysId).Find(&list)
	}

	return list, err
}

func (r *Repository) FindByPermId(permId int) ([]*model.SysRole, error) {
	list := make([]*model.SysRole, 0)

	err := r.orm.Table("sys_role").Select("sys_role.*").Join("INNER", "sys_role_permission", "sys_role.id = sys_role_permission.sys_role_id").Where(" sys_permission_id = ? ", permId).Find(&list)
	if err != nil {
		return nil, err
	}

	return list, nil
}
