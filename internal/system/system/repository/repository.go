package repository

import (
	"casbin-auth-go/dto/apireq"
	"casbin-auth-go/dto/model"
	"casbin-auth-go/internal/system/system"
	"strings"
	"xorm.io/xorm"
)

type Repository struct {
	orm   *xorm.EngineGroup
	cache system.Cache
}

func NewRepository(orm *xorm.EngineGroup, sc system.Cache) system.Repository {
	return &Repository{
		orm:   orm,
		cache: sc,
	}
}

func (r *Repository) Insert(req *apireq.AddSystem) error {
	sys := model.System{
		Name:          req.Name,
		SystemType:    req.SystemType,
		Tag:           req.Tag,
		Email:         req.Email,
		Address:       req.Address,
		Tel:           req.Tel,
		Quota:         req.Quota,
		IpAddress:     strings.Join(req.IpAddress, ","),
		MacAddress:    strings.Join(req.MacAddress, ","),
		Principal:     req.Principal,
		Salesman:      req.Salesman,
		SalesmanPhone: req.SalesmanPhone,
	}

	session := r.orm.NewSession()
	defer session.Close()

	err := session.Begin()
	if err != nil {
		_ = session.Rollback()
		return err
	}

	// ----------------------- start session -----------------------

	// Insert system
	_, err = session.Insert(&sys)
	if err != nil {
		_ = session.Rollback()
		return err
	}

	// Generate uniq uuid
	uuid := ""
	uuidNotExist := true

	for uuidNotExist {
		// gen inc_no
		uuid = system.GenIncNo(req.Uuid)

		// check inc_no unique
		uuidNotExist, err = session.Where("uuid = ?", uuid).Exist(&model.System{})
		if err != nil {
			_ = session.Rollback()
			return err
		}
	}

	// Update system.uuid
	sys.Uuid = uuid
	_, err = session.ID(sys.Id).Cols("uuid").Update(sys)
	if err != nil {
		_ = session.Rollback()
		return err
	}

	// Copy permissions from system_id
	if req.CopyFromSystemId != 0 {
		permissions := make([]*model.SysPermission, 0)
		err = session.Where("system_id = ?", req.CopyFromSystemId).Find(&permissions)
		if err != nil {
			_ = session.Rollback()
			return err
		}
		for _, perm := range permissions {
			permissionRes := model.SysPermission{
				SystemId:     sys.Id,
				AllowApiPath: perm.AllowApiPath,
				Action:       perm.Action,
				Slug:         perm.Slug,
				Description:  perm.Description,
			}
			_, err = session.Insert(&permissionRes)
			if err != nil {
				_ = session.Rollback()
				return err
			}
		}
	}

	// set Default Role Admin1 Admin2 Operator
	defaultRoles := []string{"Admin"}
	batchRole := make([]model.SysRole, len(defaultRoles))
	for index, role := range defaultRoles {
		roleRes := model.SysRole{
			Sort:        index,
			SystemId:    sys.Id,
			Name:        role,
			DisplayName: role,
		}
		batchRole[index] = roleRes
	}
	_, err = session.Insert(&batchRole)
	if err != nil {
		_ = session.Rollback()
		return err
	}

	// ----------------------- end session -----------------------

	_ = session.Commit()

	return err
}

func (r *Repository) Find(listType string, offset, limit int) ([]*model.System, error) {
	var err error
	list := make([]*model.System, 0)
	switch listType {
	case system.ListTypeEnable:
		err = r.orm.Where(" is_disable = ? ", 0).Limit(limit, offset).Find(&list)
	case system.ListTypeDisable:
		err = r.orm.Where(" is_disable = ? ", 1).Limit(limit, offset).Find(&list)
	case system.ListTypeAll:
		err = r.orm.Limit(limit, offset).Find(&list)
	}

	if err != nil {
		return nil, err
	}

	return list, nil
}

func (r *Repository) FindOne(m *model.System) (*model.System, error) {
	has, err := r.orm.Get(m)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil
	}

	return m, nil
}

func (r *Repository) Count(listType string) (int, error) {
	var count int64
	var err error
	switch listType {
	case system.ListTypeEnable:
		count, err = r.orm.Where(" is_disable = ? ", 0).Count(&model.System{})
	case system.ListTypeDisable:
		count, err = r.orm.Where(" is_disable = ? ", 1).Count(&model.System{})
	case system.ListTypeAll:
		count, err = r.orm.Count(&model.System{})
	}

	if err != nil {
		return 0, err
	}

	return int(count), nil
}

func (r *Repository) Update(m *model.System) error {
	_, err := r.orm.ID(m.Id).Cols("is_disable", "name", "address", "tel", "ip_address", "mac_address").Update(m)
	if err != nil {
		return err
	}
	return nil
}
