package service

import (
	"casbin-auth-go/config"
	"casbin-auth-go/driver"
	"casbin-auth-go/dto/apireq"
	"casbin-auth-go/dto/model"
	sysPermRepo "casbin-auth-go/internal/system/sys_permission/repository"
	sysRoleRepo "casbin-auth-go/internal/system/sys_role/repository"
	sysRepo "casbin-auth-go/internal/system/system/repository"
	"casbin-auth-go/pkg/valider"
	"fmt"
	"github.com/stretchr/testify/assert"
	_ "go/types"
	"log"
	"os"
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

func TestService_ListSysPermission(t *testing.T) {
	// Arrange
	redis, _ := driver.NewRedis()
	orm, _ := driver.NewXorm()

	sc := sysRepo.NewCache(redis)
	sr := sysRepo.NewRepository(orm, sc)
	spr := sysPermRepo.NewRepository(orm)
	sps := NewService(spr, sr)

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
			6,
		},
		{
			1,
			10,
			1,
			2,
		},
		{
			2,
			10,
			1,
			4,
		},
	}
	// Act
	for _, tc := range testCases {
		tc := tc
		t.Run(fmt.Sprintf("List Sys Permission,System Id:%d,Page:%d,PerPage:%d", tc.SystemId, tc.Page, tc.PerPage), func(t *testing.T) {
			req := apireq.ListSysPermission{
				SystemId: tc.SystemId,
				Page:     tc.Page,
				PerPage:  tc.PerPage,
			}

			data, err := sps.ListSysPermission(&req)
			assert.Nil(t, err)
			assert.Len(t, data.List, tc.WantCount)
			assert.Equal(t, tc.WantCount, data.Total)
		})
	}
}

func TestService_AddSysPermission(t *testing.T) {
	// Arrange
	redis, _ := driver.NewRedis()
	orm, _ := driver.NewXorm()

	sc := sysRepo.NewCache(redis)
	sr := sysRepo.NewRepository(orm, sc)
	spr := sysPermRepo.NewRepository(orm)
	sps := NewService(spr, sr)

	req := apireq.AddSysPermission{
		AccountId:    1,
		SystemId:     1,
		AllowApiPath: "/test_path/",
		Action:       "get",
		Slug:         "test_slug",
		Description:  "test_description",
	}

	// Act
	err := sps.AddSysPermission(&req)

	// Assert
	assert.Nil(t, err)

	// Teardown
	_, _ = orm.Where("slug = ? ", req.Slug).Delete(&model.SysPermission{})
}

func TestService_EditSysPermission(t *testing.T) {
	// Arrange
	redis, _ := driver.NewRedis()
	orm, _ := driver.NewXorm()

	sc := sysRepo.NewCache(redis)
	sr := sysRepo.NewRepository(orm, sc)
	spr := sysPermRepo.NewRepository(orm)
	sps := NewService(spr, sr)

	permId := 1
	perm := model.SysPermission{Id: permId}
	_, _ = orm.Get(&perm)

	req := apireq.EditSysPermission{
		AccountId:   3,
		Slug:        "test_slug",
		Description: "test_description",
	}

	// Act
	err := sps.EditSysPermission(permId, &req)

	// Assert
	assert.Nil(t, err)

	// Teardown
	_, _ = orm.ID(permId).Update(&perm)
}

func TestService_DeleteSysPermission(t *testing.T) {
	// Arrange
	redis, _ := driver.NewRedis()
	orm, _ := driver.NewXorm()

	sc := sysRepo.NewCache(redis)
	sr := sysRepo.NewRepository(orm, sc)
	srr := sysRoleRepo.NewRepository(orm)
	spr := sysPermRepo.NewRepository(orm)
	sps := NewService(spr, sr)

	m := model.SysPermission{
		SystemId:     1,
		AllowApiPath: "/test_path/",
		Action:       "get",
		Slug:         "test_slug",
		Description:  "test_description",
	}
	_, _ = orm.Insert(&m)

	// Act
	err := sps.DeleteSysPermission(m.Id, srr)

	// Assert
	assert.Nil(t, err)
}
