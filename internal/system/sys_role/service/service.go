package service

import (
	"casbin-auth-go/dto/apireq"
	"casbin-auth-go/dto/apires"
	"casbin-auth-go/dto/model"
	"casbin-auth-go/internal/system/sys_permission"
	"casbin-auth-go/internal/system/sys_role"
	"casbin-auth-go/internal/system/system"
	"casbin-auth-go/pkg/casbin"
	"casbin-auth-go/pkg/er"
	"gopkg.in/guregu/null.v4"
	"net/http"
)

type Service struct {
	sysRepo     system.Repository
	sysRoleRepo sys_role.Repository
	sysPermRepo sys_permission.Repository
}

func NewService(srr sys_role.Repository, sr system.Repository, spr sys_permission.Repository) sys_role.Service {
	return &Service{
		sysRepo:     sr,
		sysRoleRepo: srr,
		sysPermRepo: spr,
	}
}

func (s *Service) ListSysRole(req *apireq.ListSysRole) (*apires.ListSysRole, error) {
	// Check system exist
	sysId := req.SystemId
	if sysId != 0 {
		sys, err := s.sysRepo.FindOne(&model.System{Id: sysId})
		if err != nil {
			findErr := er.NewAppErr(http.StatusInternalServerError, er.UnknownError, "find system error.", err)
			return nil, findErr
		}
		if sys == nil || sys.IsDisable {
			notFoundErr := er.NewAppErr(http.StatusBadRequest, er.ResourceNotFoundError, "system not found.", nil)
			return nil, notFoundErr
		}
	}

	page := req.Page
	perPage := req.PerPage

	if page <= 1 {
		page = 1
	}

	if perPage <= 1 {
		perPage = 1
	}

	offset := (page - 1) * perPage

	total, err := s.sysRoleRepo.Count(sysId)
	if err != nil {
		countErr := er.NewAppErr(http.StatusInternalServerError, er.UnknownError, "count sys role error.", err)
		return nil, countErr
	}

	data, err := s.sysRoleRepo.Find(sysId, offset, perPage)
	if err != nil {
		findErr := er.NewAppErr(http.StatusInternalServerError, er.UnknownError, "find sys role error.", err)
		return nil, findErr
	}

	res := &apires.ListSysRole{
		List:        data,
		Total:       total,
		CurrentPage: page,
		PerPage:     perPage,
	}

	// 判斷 offset 加上資料筆數，是否仍小於總筆數,是的話回傳下一頁頁數
	dataCount := len(data)
	if (offset + dataCount) < total {
		res.NextPage = null.IntFrom(int64(page) + int64(1))
	}

	return res, nil
}

func (s *Service) AddSysRoleWithPermission(req *apireq.AddSysRole) error {
	// Check system id
	sys, err := s.sysRepo.FindOne(&model.System{Id: req.SystemId})
	if err != nil {
		findErr := er.NewAppErr(http.StatusInternalServerError, er.UnknownError, "find system error.", err)
		return findErr
	}
	if sys == nil || sys.IsDisable {
		notFoundErr := er.NewAppErr(http.StatusBadRequest, er.ResourceNotFoundError, "system not found.", err)
		return notFoundErr
	}

	// Check permission ids
	perms, err := s.sysPermRepo.FindByIds(req.PermissionIds)
	if err != nil {
		findErr := er.NewAppErr(http.StatusInternalServerError, er.UnknownError, "find permission ids error.", err)
		return findErr
	}
	if len(perms) != len(req.PermissionIds) {
		notMatchErr := er.NewAppErr(http.StatusBadRequest, er.ErrorParamInvalid, "permission ids illegal.", nil)
		return notMatchErr
	}

	m := model.SysRole{
		SystemId:    sys.Id,
		Name:        req.Name,
		DisplayName: req.DisplayName,
	}
	sysRole, err := s.sysRoleRepo.InsertWithPermission(&m, req.PermissionIds)
	if err != nil {
		insertErr := er.NewAppErr(http.StatusInternalServerError, er.DBInsertError, "insert sys role and sys role permission error.", err)
		return insertErr
	}

	// Casbin
	_, _ = casbin.AddRolePolicy(sys.Tag, sysRole.Id, perms)

	return nil
}

func (s *Service) EditSysRoleWithPermission(sysRoleId int, req *apireq.EditSysRole) error {
	role, err := s.sysRoleRepo.FindOne(&model.SysRole{Id: sysRoleId})
	if err != nil {
		findErr := er.NewAppErr(http.StatusInternalServerError, er.UnknownError, "find sys role error.", err)
		return findErr
	}
	if role == nil {
		notFoundErr := er.NewAppErr(http.StatusBadRequest, er.ResourceNotFoundError, "sys role not found.", nil)
		return notFoundErr
	}

	if req.Name != "" {
		role.Name = req.Name
	}
	if req.DisplayName != "" {
		role.DisplayName = req.DisplayName
	}
	if req.IsDisable != nil {
		role.IsDisable = *req.IsDisable
	}

	// Check permission ids
	perms, err := s.sysPermRepo.FindByIds(req.PermissionIds)
	if err != nil {
		findErr := er.NewAppErr(http.StatusInternalServerError, er.UnknownError, "check permission ids error.", err)
		return findErr
	}
	if len(perms) != len(req.PermissionIds) {
		notMatchErr := er.NewAppErr(http.StatusBadRequest, er.ErrorParamInvalid, "permission ids illegal.", nil)
		return notMatchErr
	}

	err = s.sysRoleRepo.UpdateWithPermission(role, req.PermissionIds)
	if err != nil {
		updateErr := er.NewAppErr(http.StatusInternalServerError, er.UnknownError, "update sys role and sys role permission error.", err)
		return updateErr
	}

	// Casbin
	sys, _ := s.sysRepo.FindOne(&model.System{Id: role.SystemId})
	if sys != nil {
		_, err = casbin.RemoveRolePolicy(sysRoleId)
		if err != nil {
			return nil
		}
		_, _ = casbin.AddRolePolicy(sys.Tag, sysRoleId, perms)
	}

	return nil
}

func (s *Service) GetSysRoleWithPermission(sysRoleId int) (*apires.SysRoleWithPermissionIds, error) {
	role, err := s.sysRoleRepo.FindOne(&model.SysRole{Id: sysRoleId})
	if err != nil {
		findErr := er.NewAppErr(http.StatusInternalServerError, er.UnknownError, "find sys role error.", err)
		return nil, findErr
	}
	if role == nil {
		notFoundErr := er.NewAppErr(http.StatusBadRequest, er.ResourceNotFoundError, "sys role not found.", nil)
		return nil, notFoundErr
	}

	perms, err := s.sysPermRepo.FindIdsByRole(sysRoleId)
	if err != nil {
		findErr := er.NewAppErr(http.StatusInternalServerError, er.UnknownError, "Get Role Permission Ids error", nil)
		return nil, findErr
	}

	result := apires.SysRoleWithPermissionIds{
		SysRole:       role,
		PermissionIds: perms,
	}

	return &result, nil
}
