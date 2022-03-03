package service

import (
	"casbin-auth-go/config"
	"casbin-auth-go/driver"
	"casbin-auth-go/dto/apireq"
	"casbin-auth-go/dto/model"
	sysPermRepo "casbin-auth-go/internal/system/sys_permission/repository"
	sysRoleRepo "casbin-auth-go/internal/system/sys_role/repository"
	sysRepo "casbin-auth-go/internal/system/system/repository"
	"casbin-auth-go/pkg/er"
	"casbin-auth-go/pkg/valider"
	"github.com/stretchr/testify/assert"
	_ "go/types"
	"log"
	"net/http"
	"os"
	"strconv"
	"testing"

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

func TestService_AddSysRoleWithPermission(t *testing.T) {
	// Arrange
	redis, _ := driver.NewRedis()
	orm, _ := driver.NewXorm()

	sc := sysRepo.NewCache(redis)
	sr := sysRepo.NewRepository(orm, sc)
	spr := sysPermRepo.NewRepository(orm)
	srr := sysRoleRepo.NewRepository(orm)
	srs := NewService(srr, sr, spr)

	req := apireq.AddSysRole{
		AccountId:     2,
		SystemId:      1,
		Name:          "test_role",
		DisplayName:   "test_display_name",
		PermissionIds: nil,
	}

	// Act
	err := srs.AddSysRoleWithPermission(&req)

	// Assert
	assert.Nil(t, err)

	// Teardown
	_, _ = orm.Where("name = ? ", req.Name).Delete(&model.SysRole{})
}

func TestService_EditSysRoleWithPermission(t *testing.T) {
	// Arrange
	redis, _ := driver.NewRedis()
	orm, _ := driver.NewXorm()

	sc := sysRepo.NewCache(redis)
	sr := sysRepo.NewRepository(orm, sc)
	spr := sysPermRepo.NewRepository(orm)
	srr := sysRoleRepo.NewRepository(orm)
	srs := NewService(srr, sr, spr)

	sysRoleId := 1
	sysRole := model.SysRole{Id: sysRoleId}
	_, _ = orm.Get(&sysRole)

	status := false
	req := apireq.EditSysRole{
		AccountId:     1,
		Name:          "test_name",
		DisplayName:   "test_display_name",
		IsDisable:     &status,
		PermissionIds: []int{1, 2},
	}

	// Act
	err := srs.EditSysRoleWithPermission(sysRoleId, &req)

	// Assert
	assert.Nil(t, err)

	// Teardown
	_, _ = orm.ID(sysRoleId).Update(&sysRole)
}

func TestService_GetSysRoleWithPermission(t *testing.T) {
	// Arrange
	redis, _ := driver.NewRedis()
	orm, _ := driver.NewXorm()

	sc := sysRepo.NewCache(redis)
	sr := sysRepo.NewRepository(orm, sc)
	spr := sysPermRepo.NewRepository(orm)
	srr := sysRoleRepo.NewRepository(orm)
	srs := NewService(srr, sr, spr)

	// No data
	// Act
	sysRoleId := 100
	role, err := srs.GetSysRoleWithPermission(sysRoleId)

	// Assert
	assert.NotNil(t, err)
	assert.Nil(t, role)
	notFoundErr := err.(*er.AppError)
	assert.Equal(t, http.StatusBadRequest, notFoundErr.StatusCode)
	assert.Equal(t, strconv.Itoa(er.ResourceNotFoundError), notFoundErr.Code)

	// Has data
	// Act
	sysRoleId = 1
	role, err = srs.GetSysRoleWithPermission(sysRoleId)

	// Assert
	assert.Nil(t, err)
	assert.NotNil(t, role)
	assert.Equal(t, sysRoleId, role.Id)
	assert.Len(t, role.PermissionIds, 2)
}
