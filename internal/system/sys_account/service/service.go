package service

import (
	"casbin-auth-go/dto/apireq"
	"casbin-auth-go/dto/apires"
	"casbin-auth-go/dto/model"
	"casbin-auth-go/internal/system/sys_account"
	"casbin-auth-go/internal/system/sys_role"
	"casbin-auth-go/internal/system/system"
	"casbin-auth-go/internal/token"
	"casbin-auth-go/pkg/casbin"
	"casbin-auth-go/pkg/er"
	"casbin-auth-go/pkg/helper"
	"net/http"
	"time"

	"gopkg.in/guregu/null.v4"
)

type Service struct {
	sysAccRepo  sys_account.Repository
	sysRepo     system.Repository
	sysRoleRepo sys_role.Repository
	tokenCache  token.Cache
}

func NewService(sar sys_account.Repository, sr system.Repository, srr sys_role.Repository, tc token.Cache) sys_account.Service {
	return &Service{
		sysAccRepo:  sar,
		sysRepo:     sr,
		sysRoleRepo: srr,
		tokenCache:  tc,
	}
}

func (s *Service) ListSysAccount(req *apireq.ListSysAccount) (*apires.ListSysAccount, error) {
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

	total, err := s.sysAccRepo.Count(sysId)
	if err != nil {
		countErr := er.NewAppErr(http.StatusInternalServerError, er.UnknownError, "count sys account error.", err)
		return nil, countErr
	}

	data, err := s.sysAccRepo.Find(sysId, offset, perPage)
	if err != nil {
		findErr := er.NewAppErr(http.StatusInternalServerError, er.UnknownError, "find sys account error.", err)
		return nil, findErr
	}

	res := &apires.ListSysAccount{
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

func (s *Service) GetSysAccount(sysAccId int) (*apires.SysAccount, error) {
	// Check Search account id exist
	acc, err := s.sysAccRepo.FindOne(&model.SysAccount{Id: sysAccId})
	if err != nil {
		findErr := er.NewAppErr(http.StatusInternalServerError, er.UnknownError, "find sys account error.", err)
		return nil, findErr
	}
	if acc == nil {
		notFoundErr := er.NewAppErr(http.StatusBadRequest, er.ResourceNotFoundError, "sys account not found.", nil)
		return nil, notFoundErr
	}

	res := apires.SysAccount{
		Id:        acc.Id,
		SystemId:  acc.SystemId,
		Account:   acc.Account,
		Phone:     acc.Phone,
		Email:     acc.Email,
		Name:      acc.Name,
		IsDisable: acc.IsDisable,
		VerifyAt:  acc.VerifyAt,
		CreatedAt: acc.CreatedAt,
		UpdatedAt: acc.UpdatedAt,
	}

	return &res, nil
}

func (s *Service) AddSysAccount(req *apireq.AddSysAccount) error {
	// Check account id exist
	acc, err := s.sysAccRepo.FindOne(&model.SysAccount{Id: req.AccountId})
	if err != nil {
		findErr := er.NewAppErr(http.StatusInternalServerError, er.UnknownError, "find sys_account error.", err)
		return findErr
	}
	if acc == nil || acc.IsDisable {
		notFoundErr := er.NewAppErr(http.StatusBadRequest, er.ResourceNotFoundError, "account not found.", nil)
		return notFoundErr
	}

	// Get system
	sys, err := s.sysRepo.FindOne(&model.System{Id: req.SystemId})
	if err != nil {
		findErr := er.NewAppErr(http.StatusInternalServerError, er.UnknownError, "find system error.", err)
		return findErr
	}
	if sys == nil || sys.IsDisable {
		notFoundErr := er.NewAppErr(http.StatusBadRequest, er.ResourceNotFoundError, "system not found.", nil)
		return notFoundErr
	}

	// Check sys_role_id
	role, err := s.sysRoleRepo.FindOne(&model.SysRole{Id: req.RoleId})
	if err != nil {
		findErr := er.NewAppErr(http.StatusInternalServerError, er.UnknownError, "find sys role error.", err)
		return findErr
	}
	if role == nil || role.IsDisable {
		notFoundErr := er.NewAppErr(http.StatusBadRequest, er.ResourceNotFoundError, "sys role not found.", nil)
		return notFoundErr
	}
	if role.SystemId != sys.Id {
		notMatchErr := er.NewAppErr(http.StatusBadRequest, er.ResourceNotFoundError, "sys role not match system.", nil)
		return notMatchErr
	}

	// Check account exist
	ac, err := s.sysAccRepo.FindOne(&model.SysAccount{
		SystemId: req.SystemId,
		Account:  req.Account,
	})
	if err != nil {
		findErr := er.NewAppErr(http.StatusInternalServerError, er.UnknownError, "find sys account error.", err)
		return findErr
	}
	if ac != nil {
		existErr := er.NewAppErr(http.StatusBadRequest, er.ErrorParamInvalid, "sys account exist.", nil)
		return existErr
	}

	// Check email exist
	emAc, err := s.sysAccRepo.FindOne(&model.SysAccount{
		SystemId: req.SystemId,
		Email:    req.Email,
	})
	if err != nil {
		findErr := er.NewAppErr(http.StatusInternalServerError, er.UnknownError, "find sys account error.", err)
		return findErr
	}
	if emAc != nil {
		existErr := er.NewAppErr(http.StatusBadRequest, er.ErrorParamInvalid, "email exist.", nil)
		return existErr
	}

	// Insert account and role
	// Gen password
	password := helper.RandString(8)
	sysAccountModel := model.SysAccount{
		SystemId: req.SystemId,
		Account:  req.Account,
		Phone:    req.Phone,
		Email:    req.Email,
		Password: helper.ScryptStr(password),
		Name:     req.Name,
	}

	sysAccount, err := s.sysAccRepo.InsertWithRole(&sysAccountModel, req.RoleId)
	if err != nil {
		insertErr := er.NewAppErr(http.StatusInternalServerError, er.UnknownError, "insert sys account error.", err)
		return insertErr
	}

	// 從樣板產生信件內容，並寄送密碼信
	// TODO - Send mail

	// Casbin
	_, _ = casbin.AddUserRole(sys.Tag, sysAccount.Id, 0, req.RoleId)

	return nil
}

func (s *Service) EditSysAccount(sysAccId int, req *apireq.EditSysAccount) error {
	// Check account id exist
	acc, err := s.sysAccRepo.FindOne(&model.SysAccount{Id: req.AccountId})
	if err != nil {
		findErr := er.NewAppErr(http.StatusInternalServerError, er.UnknownError, "find sys account error.", err)
		return findErr
	}
	if acc == nil || acc.IsDisable {
		notFoundErr := er.NewAppErr(http.StatusBadRequest, er.ResourceNotFoundError, "sys account not found.", nil)
		return notFoundErr
	}

	// check account id
	ac, err := s.sysAccRepo.FindOneWithRole(sysAccId)
	if err != nil {
		findErr := er.NewAppErr(http.StatusInternalServerError, er.UnknownError, "find sys account error.", err)
		return findErr
	}
	if ac == nil {
		notFoundErr := er.NewAppErr(http.StatusBadRequest, er.ResourceNotFoundError, "sys account not found.", nil)
		return notFoundErr
	}

	// get system data
	sys, err := s.sysRepo.FindOne(&model.System{Id: ac.SystemId})
	if err != nil {
		findErr := er.NewAppErr(http.StatusInternalServerError, er.UnknownError, "find system error.", err)
		return findErr
	}
	if sys == nil || sys.IsDisable {
		notFoundErr := er.NewAppErr(http.StatusBadRequest, er.ResourceNotFoundError, "system not found.", nil)
		return notFoundErr
	}

	// Check sys_role_id
	role, err := s.sysRoleRepo.FindOne(&model.SysRole{Id: req.RoleId})
	if err != nil {
		findErr := er.NewAppErr(http.StatusInternalServerError, er.UnknownError, "find sys role error.", err)
		return findErr
	}
	if role == nil || role.IsDisable {
		notFoundErr := er.NewAppErr(http.StatusBadRequest, er.ResourceNotFoundError, "sys role not found.", nil)
		return notFoundErr
	}
	if role.SystemId != sys.Id {
		notMatchErr := er.NewAppErr(http.StatusBadRequest, er.ResourceNotFoundError, "sys role not match system.", nil)
		return notMatchErr
	}

	updateAcc := model.SysAccount{
		Id:        sysAccId,
		Phone:     ac.Phone,
		Email:     ac.Email,
		Name:      ac.Name,
		IsDisable: ac.IsDisable,
	}
	if req.Phone != "" {
		updateAcc.Phone = req.Phone
	}
	if req.Email != "" {
		// Check email exist
		emAc, err := s.sysAccRepo.FindOne(&model.SysAccount{
			SystemId: ac.SystemId,
			Email:    req.Email,
		})
		if err != nil {
			findErr := er.NewAppErr(http.StatusInternalServerError, er.UnknownError, "find sys account error.", err)
			return findErr
		}
		if emAc != nil && emAc.Id != ac.Id {
			existErr := er.NewAppErr(http.StatusBadRequest, er.ErrorParamInvalid, "email exist.", nil)
			return existErr
		}
		updateAcc.Email = req.Email
	}
	if req.Name != "" {
		updateAcc.Name = req.Name
	}
	if req.IsDisable != nil {
		updateAcc.IsDisable = *req.IsDisable
	}

	err = s.sysAccRepo.UpdateWithRole(&updateAcc, req.RoleId)
	if err != nil {
		updateErr := er.NewAppErr(http.StatusInternalServerError, er.UnknownError, "update sys account and sys account role error.", err)
		return updateErr
	}

	if acc.IsDisable {
		iat := time.Now().UTC().Unix()
		_ = s.tokenCache.SetTokenIat(acc.Id, iat)
	}

	// Casbin
	_, _ = casbin.AddUserRole(sys.Tag, sysAccId, ac.SysRoleId, req.RoleId)

	return nil
}

func (s *Service) ChangePassword(sysAccId int, oldPw, newPw string) error {
	acc, err := s.sysAccRepo.FindOne(&model.SysAccount{Id: sysAccId})
	if err != nil {
		findErr := er.NewAppErr(http.StatusInternalServerError, er.UnknownError, "find sys account error.", err)
		return findErr
	}
	if acc == nil {
		notFoundErr := er.NewAppErr(http.StatusBadRequest, er.ResourceNotFoundError, "sys account not found.", nil)
		return notFoundErr
	}
	if acc.Password != helper.ScryptStr(oldPw) {
		notMatchErr := er.NewAppErr(http.StatusBadRequest, er.ErrorParamInvalid, "incorrect password.", nil)
		return notMatchErr
	}

	newPw = helper.ScryptStr(newPw)

	err = s.sysAccRepo.UpdatePassword(sysAccId, newPw)
	if err != nil {
		updateErr := er.NewAppErr(http.StatusInternalServerError, er.DBUpdateError, "update sys account error.", err)
		return updateErr
	}

	return nil
}

func (s *Service) ForgotPasswordByAdmin(sysAccId int, req *apireq.ForgotSysAccountPassword) error {
	// Check 執行者 account
	account, err := s.sysAccRepo.FindOne(&model.SysAccount{Id: req.AccountId})
	if err != nil {
		findErr := er.NewAppErr(http.StatusInternalServerError, er.UnknownError, "find sys account error.", err)
		return findErr
	}
	if account == nil || account.IsDisable {
		notFoundErr := er.NewAppErr(http.StatusBadRequest, er.ResourceNotFoundError, "sys account not found.", nil)
		return notFoundErr
	}

	whereAcc := model.SysAccount{Id: sysAccId}
	acc, err := s.sysAccRepo.FindOne(&whereAcc)
	if err != nil {
		findErr := er.NewAppErr(http.StatusInternalServerError, er.UnknownError, "find sys account error.", err)
		return findErr
	}
	if acc == nil || acc.IsDisable {
		notFoundErr := er.NewAppErr(http.StatusBadRequest, er.ResourceNotFoundError, "sys account not found.", nil)
		return notFoundErr
	}

	// get system data
	sys, err := s.sysRepo.FindOne(&model.System{Id: acc.SystemId})
	if err != nil {
		findErr := er.NewAppErr(http.StatusInternalServerError, er.UnknownError, "find system error.", err)
		return findErr
	}
	if sys == nil || sys.IsDisable {
		notFoundErr := er.NewAppErr(http.StatusBadRequest, er.ResourceNotFoundError, "system not found.", nil)
		return notFoundErr
	}

	// gen password
	password := helper.RandString(8)
	err = s.sysAccRepo.UpdatePassword(sysAccId, helper.ScryptStr(password))
	if err != nil {
		updateErr := er.NewAppErr(http.StatusInternalServerError, er.UnknownError, "update sys account error.", err)
		return updateErr
	}

	// TODO - Send email

	iat := time.Now().UTC().Unix()
	_ = s.tokenCache.SetTokenIat(sysAccId, iat)

	return nil
}
