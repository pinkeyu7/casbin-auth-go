package service

import (
	"casbin-auth-go/dto/apireq"
	"casbin-auth-go/dto/apires"
	"casbin-auth-go/dto/model"
	"casbin-auth-go/internal/system/system"
	"casbin-auth-go/pkg/er"
	"casbin-auth-go/pkg/helper"
	"gopkg.in/guregu/null.v4"
	"net/http"
	"strings"
)

type Service struct {
	sysRepo system.Repository
}

func NewService(sr system.Repository) system.Service {
	return &Service{sysRepo: sr}
}

func (s *Service) ListSystem(req *apireq.ListSystem) (*apires.ListSystem, error) {
	listType := req.ListType
	page := req.Page
	perPage := req.PerPage

	if page <= 1 {
		page = 1
	}

	if perPage <= 1 {
		perPage = 1
	}

	offset := (page - 1) * perPage

	listTypes := []string{system.ListTypeEnable, system.ListTypeDisable, system.ListTypeAll}
	if contains := helper.StringContains(listTypes, listType); !contains {
		filterErr := er.NewAppErr(http.StatusBadRequest, er.ErrorParamInvalid, "list type error.", nil)
		return nil, filterErr
	}

	total, err := s.sysRepo.Count(listType)
	if err != nil {
		countErr := er.NewAppErr(http.StatusInternalServerError, er.UnknownError, "count system error.", err)
		return nil, countErr
	}

	data, err := s.sysRepo.Find(listType, offset, perPage)
	if err != nil {
		selectErr := er.NewAppErr(http.StatusInternalServerError, er.UnknownError, "find system error.", err)
		return nil, selectErr
	}

	res := &apires.ListSystem{
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

func (s *Service) GetSystem(sysId int) (*model.System, error) {
	sys, err := s.sysRepo.FindOne(&model.System{Id: sysId})
	if err != nil {
		findErr := er.NewAppErr(http.StatusInternalServerError, er.UnknownError, "find system error.", err)
		return nil, findErr
	}
	if sys == nil {
		notFoundErr := er.NewAppErr(http.StatusBadRequest, er.ResourceNotFoundError, "system not found.", nil)
		return nil, notFoundErr
	}

	return sys, nil
}

func (s *Service) AddSystem(req *apireq.AddSystem) error {
	err := s.sysRepo.Insert(req)
	if err != nil {
		insertErr := er.NewAppErr(http.StatusInternalServerError, er.DBInsertError, "insert system error.", err)
		return insertErr
	}
	return nil
}

func (s *Service) EditSystem(sysId int, req *apireq.EditSystem) error {
	// Check System Exist
	sys, err := s.sysRepo.FindOne(&model.System{Id: sysId})
	if err != nil {
		findErr := er.NewAppErr(http.StatusInternalServerError, er.UnknownError, "find system error.", err)
		return findErr
	}
	if sys == nil {
		notFoundErr := er.NewAppErr(http.StatusBadRequest, er.ResourceNotFoundError, "system not found.", nil)
		return notFoundErr
	}

	if req.Name != "" {
		sys.Name = req.Name
	}
	if req.Address != "" {
		sys.Address = req.Address
	}
	if req.Tel != "" {
		sys.Tel = req.Tel
	}
	if req.IsDisable != nil {
		sys.IsDisable = *req.IsDisable
	}
	if len(req.IpAddress) > 0 {
		sys.IpAddress = strings.Join(req.IpAddress, ",")
	}
	if len(req.MacAddress) > 0 {
		sys.MacAddress = strings.Join(req.MacAddress, ",")
	}

	err = s.sysRepo.Update(sys)
	if err != nil {
		updateErr := er.NewAppErr(http.StatusInternalServerError, er.UnknownError, "update system error.", err)
		return updateErr
	}
	return nil
}
