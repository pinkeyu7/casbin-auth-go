package service

import (
	"casbin-auth-go/config"
	"casbin-auth-go/driver"
	"casbin-auth-go/dto/apireq"
	"casbin-auth-go/dto/model"
	sysAccRepo "casbin-auth-go/internal/system/sys_account/repository"
	sysRoleRepo "casbin-auth-go/internal/system/sys_role/repository"
	sysRepo "casbin-auth-go/internal/system/system/repository"
	tokenRepo "casbin-auth-go/internal/token/repository"
	"casbin-auth-go/pkg/er"
	"net/http"
	"strconv"

	"casbin-auth-go/pkg/valider"
	"fmt"
	_ "go/types"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	_ "github.com/brianvoe/gofakeit/v4"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func TestMain(m *testing.M) {
	setUp()
	code := m.Run()
	os.Exit(code)
}

func setUp() {
	remoteBranch := os.Getenv("REMOTE_BRANCH")
	if remoteBranch == "" {
		// load env
		err := godotenv.Load(config.GetBasePath() + "/.env")
		if err != nil {
			log.Panicln(err)
		}
	}

	valider.Init()
}

func TestService_ListSysAccount(t *testing.T) {
	// Arrange
	redis, _ := driver.NewRedis()
	orm, _ := driver.NewXorm()

	tc := tokenRepo.NewCache(redis)
	sc := sysRepo.NewCache(redis)
	sr := sysRepo.NewRepository(orm, sc)
	srr := sysRoleRepo.NewRepository(orm)
	sar := sysAccRepo.NewRepository(orm)
	sas := NewService(sar, sr, srr, tc)

	// Act
	testCases := []struct {
		SystemId  int
		PerPage   int
		Page      int
		WantCount int
	}{
		{
			0,
			10,
			1,
			4,
		},
		{
			1,
			10,
			1,
			1,
		},
		{
			2,
			10,
			1,
			3,
		},
	}
	// Act
	for _, tc := range testCases {
		tc := tc
		t.Run(fmt.Sprintf("List Sys Account,System Id:%d,Page:%d,PerPage:%d", tc.SystemId, tc.Page, tc.PerPage), func(t *testing.T) {
			req := apireq.ListSysAccount{
				SystemId: tc.SystemId,
				Page:     tc.Page,
				PerPage:  tc.PerPage,
			}

			data, err := sas.ListSysAccount(&req)
			assert.Nil(t, err)
			assert.Len(t, data.List, tc.WantCount)
			assert.Equal(t, tc.WantCount, data.Total)
		})
	}
}

func TestService_GetSysAccount(t *testing.T) {
	// Arrange
	redis, _ := driver.NewRedis()
	orm, _ := driver.NewXorm()

	tc := tokenRepo.NewCache(redis)
	sc := sysRepo.NewCache(redis)
	sr := sysRepo.NewRepository(orm, sc)
	srr := sysRoleRepo.NewRepository(orm)
	sar := sysAccRepo.NewRepository(orm)
	sas := NewService(sar, sr, srr, tc)

	// No data
	// Act
	sysAccId := 100
	acc, err := sas.GetSysAccount(sysAccId)

	// Assert
	assert.NotNil(t, err)
	assert.Nil(t, acc)
	notFoundErr := err.(*er.AppError)
	assert.Equal(t, http.StatusBadRequest, notFoundErr.StatusCode)
	assert.Equal(t, strconv.Itoa(er.ResourceNotFoundError), notFoundErr.Code)

	// Has data
	// Act
	sysAccId = 1
	acc, err = sas.GetSysAccount(sysAccId)

	// Assert
	assert.Nil(t, err)
	assert.NotNil(t, acc)
	assert.Equal(t, sysAccId, acc.Id)
}

func TestService_AddSysAccount(t *testing.T) {
	// Arrange
	redis, _ := driver.NewRedis()
	orm, _ := driver.NewXorm()

	tc := tokenRepo.NewCache(redis)
	sc := sysRepo.NewCache(redis)
	sr := sysRepo.NewRepository(orm, sc)
	srr := sysRoleRepo.NewRepository(orm)
	sar := sysAccRepo.NewRepository(orm)
	sas := NewService(sar, sr, srr, tc)

	req := apireq.AddSysAccount{
		AccountId: 1,
		SystemId:  1,
		Account:   "testing_account",
		Phone:     "+886-905-888-888",
		Email:     "testing_account@testmail.com",
		Name:      "test_name",
		RoleId:    0,
	}

	// Act
	err := sas.AddSysAccount(&req)

	// Assert
	assert.Nil(t, err)

	// Teardown
	acc := model.SysAccount{}
	_, _ = orm.Where("email = ?", req.Email).Get(&acc)
	_, _ = orm.ID(acc.Id).Delete(&model.SysAccount{})
	_, _ = orm.Where("sys_account_id = ?", acc.Id).Delete(&model.SysAccountRole{})
}

func TestService_EditSysAccount(t *testing.T) {
	// Arrange
	redis, _ := driver.NewRedis()
	orm, _ := driver.NewXorm()

	tc := tokenRepo.NewCache(redis)
	sc := sysRepo.NewCache(redis)
	sr := sysRepo.NewRepository(orm, sc)
	srr := sysRoleRepo.NewRepository(orm)
	sar := sysAccRepo.NewRepository(orm)
	sas := NewService(sar, sr, srr, tc)

	sysAccId := 3
	acc := model.SysAccount{Id: sysAccId}
	_, _ = orm.Get(&acc)
	accRole := model.SysAccountRole{SysAccountId: sysAccId}
	_, _ = orm.Get(&accRole)

	status := false
	req := apireq.EditSysAccount{
		AccountId: 1,
		Phone:     "+886-905-888-888",
		Email:     "testing_email",
		Name:      "test_name",
		IsDisable: &status,
		RoleId:    4,
	}

	// Act
	err := sas.EditSysAccount(sysAccId, &req)

	// Assert
	assert.Nil(t, err)

	// Teardown
	_, _ = orm.ID(sysAccId).Update(&acc)
	_, _ = orm.Where("sys_account_id = ?", sysAccId).Update(&accRole)
}

func TestService_ChangePassword(t *testing.T) {
	// Arrange
	redis, _ := driver.NewRedis()
	orm, _ := driver.NewXorm()

	tc := tokenRepo.NewCache(redis)
	sc := sysRepo.NewCache(redis)
	sr := sysRepo.NewRepository(orm, sc)
	srr := sysRoleRepo.NewRepository(orm)
	sar := sysAccRepo.NewRepository(orm)
	sas := NewService(sar, sr, srr, tc)

	// No data
	// Act
	sysAccId := 100
	oldPwd := "123456"
	newPwd := "12345678"
	err := sas.ChangePassword(sysAccId, oldPwd, newPwd)

	// Assert
	assert.NotNil(t, err)
	notFoundErr := err.(*er.AppError)
	assert.Equal(t, http.StatusBadRequest, notFoundErr.StatusCode)
	assert.Equal(t, strconv.Itoa(er.ResourceNotFoundError), notFoundErr.Code)

	// Wrong password
	// Act
	sysAccId = 1
	oldPwd = "wrong_password"
	newPwd = "12345678"
	err = sas.ChangePassword(sysAccId, oldPwd, newPwd)

	// Assert
	assert.NotNil(t, err)
	notMatchErr := err.(*er.AppError)
	assert.Equal(t, http.StatusBadRequest, notMatchErr.StatusCode)
	assert.Equal(t, strconv.Itoa(er.ErrorParamInvalid), notMatchErr.Code)

	// Wrong password
	// Act
	sysAccId = 1
	acc := model.SysAccount{Id: sysAccId}
	_, _ = orm.Get(&acc)
	oldPwd = "123456"
	newPwd = "12345678"
	err = sas.ChangePassword(sysAccId, oldPwd, newPwd)

	// Assert
	assert.Nil(t, err)

	// Teardown
	_, _ = orm.ID(sysAccId).Update(&acc)
}

func TestService_ForgotPasswordByAdmin(t *testing.T) {
	// Arrange
	redis, _ := driver.NewRedis()
	orm, _ := driver.NewXorm()

	tc := tokenRepo.NewCache(redis)
	sc := sysRepo.NewCache(redis)
	sr := sysRepo.NewRepository(orm, sc)
	srr := sysRoleRepo.NewRepository(orm)
	sar := sysAccRepo.NewRepository(orm)
	sas := NewService(sar, sr, srr, tc)

	// No data
	// Act
	sysAccId := 100
	req := apireq.ForgotSysAccountPassword{AccountId: 1}
	err := sas.ForgotPasswordByAdmin(sysAccId, &req)

	// Assert
	assert.NotNil(t, err)
	notFoundErr := err.(*er.AppError)
	assert.Equal(t, http.StatusBadRequest, notFoundErr.StatusCode)
	assert.Equal(t, strconv.Itoa(er.ResourceNotFoundError), notFoundErr.Code)

	// Wrong password
	// Act
	sysAccId = 1
	acc := model.SysAccount{Id: sysAccId}
	_, _ = orm.Get(&acc)
	err = sas.ForgotPasswordByAdmin(sysAccId, &req)

	// Assert
	assert.Nil(t, err)

	// TeardownÏ
	_, _ = orm.ID(sysAccId).Update(&acc)
}
