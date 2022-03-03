package repository

import (
	"casbin-auth-go/dto/apires"
	"casbin-auth-go/dto/model"
	"casbin-auth-go/internal/system/sys_account"
	"xorm.io/xorm"
)

type Repository struct {
	orm *xorm.EngineGroup
}

func NewRepository(orm *xorm.EngineGroup) sys_account.Repository {
	return &Repository{orm: orm}
}

func (r *Repository) InsertWithRole(m *model.SysAccount, sysRoleId int) (*model.SysAccount, error) {
	session := r.orm.NewSession()
	defer session.Close()

	err := session.Begin()
	if err != nil {
		_ = session.Rollback()
		return nil, err
	}

	// ----------------------- start session -----------------------

	_, err = session.Insert(m)
	if err != nil {
		_ = session.Rollback()
		return nil, err
	}

	_, err = session.Insert(&model.SysAccountRole{
		SysAccountId: m.Id,
		SysRoleId:    sysRoleId,
	})
	if err != nil {
		_ = session.Rollback()
		return nil, err
	}

	// ----------------------- end session -----------------------

	_ = session.Commit()

	return m, nil
}

func (r *Repository) UpdateWithRole(m *model.SysAccount, sysRoleId int) error {
	session := r.orm.NewSession()
	defer session.Close()

	err := session.Begin()
	if err != nil {
		_ = session.Rollback()
		return err
	}

	// ----------------------- start session -----------------------
	_, err = session.ID(m.Id).Cols("name", "email", "phone", "is_disable").Update(m)
	if err != nil {
		_ = session.Rollback()
		return err
	}

	_, err = session.Where("sys_account_id = ?", m.Id).Cols("sys_role_id").Update(&model.SysAccountRole{SysRoleId: sysRoleId})
	if err != nil {
		_ = session.Rollback()
		return err
	}

	// ----------------------- end session -----------------------

	_ = session.Commit()

	return nil
}

func (r *Repository) UpdatePassword(sysAccId int, newPassword string) error {
	_, err := r.orm.ID(sysAccId).Cols("password").Update(&model.SysAccount{Password: newPassword})
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) FindOne(m *model.SysAccount) (*model.SysAccount, error) {
	has, err := r.orm.Get(m)
	if err != nil {
		return nil, err
	}

	if !has {
		return nil, nil
	}

	return m, nil
}

func (r *Repository) FindOneWithRole(sysAccId int) (*apires.SysAccountWithRole, error) {
	acc := apires.SysAccountWithRole{}
	has, err := r.orm.Table("sys_account").Join("INNER", "sys_account_role",
		"sys_account_role.sys_account_id = sys_account.id").Where("sys_account.`id` = ?", sysAccId).Get(&acc)

	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil
	}

	return &acc, nil
}

func (r *Repository) Find(sysId, offset, limit int) ([]*apires.ListSysAccountItem, error) {
	var args []interface{}
	list := make([]*apires.ListSysAccountItem, 0)

	stmt := `
SELECT 
    sa.id,
    sa.system_id,
    sa.account,
    sa.phone,
    sa.email,
    sa.name,
    sa.is_disable,
    sa.verify_at,
    sa.created_at,
    sa.updated_at,
    sar.sys_role_id
FROM
    sys_account sa
        INNER JOIN
    sys_account_role sar ON sar.sys_account_id = sa.id
WHERE
    1 = 1 
`

	if sysId != 0 {
		stmt += ` AND sa.system_id = ? `
		args = append(args, sysId)
	}

	stmt += ` LIMIT ?, ? `
	args = append(args, offset, limit)

	err := r.orm.SQL(stmt, args...).Find(&list)

	return list, err
}

func (r *Repository) Count(sysId int) (int, error) {
	var total int64

	var args []interface{}

	stmt := `
SELECT 
    COUNT(1)
FROM
    sys_account sa
        INNER JOIN
    sys_account_role sar ON sar.sys_account_id = sa.id
WHERE
    1 = 1
`

	if sysId != 0 {
		stmt += ` AND sa.system_id = ? `
		args = append(args, sysId)
	}

	has, err := r.orm.SQL(stmt, args...).Get(&total)
	if err != nil {
		return 0, err
	}

	if !has {
		return 0, nil
	}

	return int(total), nil
}
