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
	sysPermRepo sys_permission.Repository
}

func NewService(spr sys_permission.Repository, sr system.Repository) sys_permission.Service {
	return &Service{
		sysRepo:     sr,
		sysPermRepo: spr,
	}
}

func (s *Service) ListPermission(req *apireq.ListSysPermission) (*apires.ListSysPermission, error) {
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

	total, err := s.sysPermRepo.Count(sysId)
	if err != nil {
		countErr := er.NewAppErr(http.StatusInternalServerError, er.UnknownError, "count sys permission error.", err)
		return nil, countErr
	}

	data, err := s.sysPermRepo.Find(sysId, offset, perPage)
	if err != nil {
		findErr := er.NewAppErr(http.StatusInternalServerError, er.UnknownError, "find system error.", err)
		return nil, findErr
	}

	res := &apires.ListSysPermission{
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

func (s *Service) AddPermission(req *apireq.AddSysPermission) error {
	// Check slug exist
	exist, err := s.sysPermRepo.Exist(&model.SysPermission{
		SystemId: req.SystemId,
		Slug:     req.Slug,
	})
	if err != nil {
		findErr := er.NewAppErr(http.StatusInternalServerError, er.UnknownError, "find sys permission error.", err)
		return findErr
	}
	if exist {
		existErr := er.NewAppErr(http.StatusBadRequest, er.ErrorParamInvalid, "slug already exist.", nil)
		return existErr
	}

	m := model.SysPermission{
		SystemId:     req.SystemId,
		AllowApiPath: req.AllowApiPath,
		Action:       req.Action,
		Slug:         req.Slug,
		Description:  req.Description,
	}

	err = s.sysPermRepo.Insert(&m)
	if err != nil {
		insertErr := er.NewAppErr(http.StatusInternalServerError, er.UnknownError, "insert sys permission error.", err)
		return insertErr
	}

	return nil
}

func (s *Service) EditPermission(sysPermId int, req *apireq.EditSysPermission) error {
	// Check sys permission exist
	perm, err := s.sysPermRepo.FindOne(&model.SysPermission{Id: sysPermId})
	if err != nil {
		findErr := er.NewAppErr(http.StatusInternalServerError, er.UnknownError, "find permission error.", err)
		return findErr
	}
	if perm == nil {
		notFoundErr := er.NewAppErr(http.StatusBadRequest, er.ResourceNotFoundError, "permission not found.", nil)
		return notFoundErr
	}

	if req.Slug != "" {
		perm.Slug = req.Slug
	}
	if req.Description != "" {
		perm.Description = req.Description
	}

	err = s.sysPermRepo.Update(perm)
	if err != nil {
		updateErr := er.NewAppErr(http.StatusInternalServerError, er.UnknownError, "update permission error.", err)
		return updateErr
	}

	return nil
}

func (s *Service) DeletePermission(sysPermId int, sysRoleRepo sys_role.Repository) error {
	// Check sys permission exist
	perm, err := s.sysPermRepo.FindOne(&model.SysPermission{Id: sysPermId})
	if err != nil {
		findErr := er.NewAppErr(http.StatusInternalServerError, er.UnknownError, "find permission error.", err)
		return findErr
	}
	if perm == nil {
		notFoundErr := er.NewAppErr(http.StatusBadRequest, er.ResourceNotFoundError, "permission not found.", nil)
		return notFoundErr
	}

	// Check system exist
	sys, err := s.sysRepo.FindOne(&model.System{Id: perm.SystemId})
	if err != nil {
		findErr := er.NewAppErr(http.StatusInternalServerError, er.UnknownError, "find system error.", err)
		return findErr
	}
	if sys == nil {
		notFoundErr := er.NewAppErr(http.StatusBadRequest, er.ResourceNotFoundError, "system not found.", nil)
		return notFoundErr
	}

	// Get affected roles
	roles, err := sysRoleRepo.FindByPermId(sysPermId)
	if err != nil {
		findErr := er.NewAppErr(http.StatusInternalServerError, er.UnknownError, "find sys role error.", err)
		return findErr
	}

	// Delete permission and role permission
	err = s.sysPermRepo.Delete(sysPermId)
	if err != nil {
		deleteErr := er.NewAppErr(http.StatusInternalServerError, er.UnknownError, "delete sys permission error.", err)
		return deleteErr
	}

	// Casbin
	for _, role := range roles {
		rolePerms, err := s.sysPermRepo.FindIdsByRole(role.Id)
		if err != nil {
			findErr := er.NewAppErr(http.StatusInternalServerError, er.UnknownError, "find sys role permission ids error.", err)
			return findErr
		}

		perms, err := s.sysPermRepo.FindByIds(rolePerms)
		if err != nil {
			findErr := er.NewAppErr(http.StatusInternalServerError, er.UnknownError, "find sys permission error.", err)
			return findErr
		}

		_, err = casbin.RemoveRolePolicy(role.Id)
		if err != nil {
			return nil
		}
		_, _ = casbin.AddRolePolicy(sys.Tag, role.Id, perms)
	}

	return err
}
