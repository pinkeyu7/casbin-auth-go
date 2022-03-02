package service

import (
	"casbin-auth-go/config"
	"casbin-auth-go/driver"
	"casbin-auth-go/dto/apireq"
	"casbin-auth-go/dto/model"
	"casbin-auth-go/internal/system/system"
	sysRepo "casbin-auth-go/internal/system/system/repository"
	"casbin-auth-go/pkg/valider"
	"fmt"
	"github.com/brianvoe/gofakeit/v4"
	_ "go/types"
	"log"
	"os"
	"testing"

	_ "github.com/brianvoe/gofakeit/v4"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
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

func TestService_ListSystem(t *testing.T) {
	// Arrange
	redis, _ := driver.NewRedis()
	orm, _ := driver.NewXorm()

	sc := sysRepo.NewCache(redis)
	sr := sysRepo.NewRepository(orm, sc)
	ss := NewService(sr)

	// Act
	testCases := []struct {
		ListType  string
		PerPage   int
		Page      int
		WantCount int
	}{
		{
			system.ListTypeAll,
			10,
			1,
			2,
		},
		{
			system.ListTypeEnable,
			10,
			1,
			2,
		},
		{
			system.ListTypeDisable,
			10,
			1,
			0,
		},
	}
	// Act
	for _, tc := range testCases {
		tc := tc
		t.Run(fmt.Sprintf("List System,ListType:%s,Page:%d,PerPage:%d", tc.ListType, tc.Page, tc.PerPage), func(t *testing.T) {
			req := apireq.ListSystem{
				ListType: tc.ListType,
				Page:     tc.Page,
				PerPage:  tc.PerPage,
			}

			data, err := ss.ListSystem(&req)
			assert.Nil(t, err)
			assert.Len(t, data.List, tc.WantCount)
			assert.Equal(t, tc.WantCount, data.Total)
		})
	}
}

func TestService_AddSystem(t *testing.T) {
	// Arrange
	redis, _ := driver.NewRedis()
	orm, _ := driver.NewXorm()

	sc := sysRepo.NewCache(redis)
	sr := sysRepo.NewRepository(orm, sc)
	ss := NewService(sr)

	// Act
	req := apireq.AddSystem{
		AccountId:        3,
		Name:             "test system insert name",
		SystemType:       system.TypeAdminBackend,
		Tag:              system.TypeAdminBackend + "_test",
		Email:            gofakeit.Email(),
		Address:          gofakeit.Address().Address,
		Tel:              gofakeit.Phone(),
		Uuid:             gofakeit.SSN(),
		Quota:            0,
		IpAddress:        nil,
		MacAddress:       nil,
		Principal:        "",
		Salesman:         "",
		SalesmanPhone:    "",
		CopyFromSystemId: 0,
	}

	err := ss.AddSystem(&req)

	// Assert
	assert.Nil(t, err)

	// Teardown
	sys := model.System{Name: req.Name}
	_, _ = orm.Get(&sys)
	_, _ = orm.Where("name = ? ", req.Name).Delete(&model.System{})
	_, _ = orm.Where("system_id = ? ", sys.Id).Delete(&model.SysRole{})
}

func TestService_EditSystem(t *testing.T) {
	// Arrange
	redis, _ := driver.NewRedis()
	orm, _ := driver.NewXorm()

	sc := sysRepo.NewCache(redis)
	sr := sysRepo.NewRepository(orm, sc)
	ss := NewService(sr)

	sysId := 1
	status := false
	req := apireq.EditSystem{
		AccountId:  3,
		Name:       "test update name",
		Address:    gofakeit.Address().Address,
		Tel:        gofakeit.Phone(),
		IsDisable:  &status,
		IpAddress:  nil,
		MacAddress: nil,
	}
	sys := model.System{Id: sysId}
	_, _ = orm.Get(&sys)

	// Act
	err := ss.EditSystem(sysId, &req)

	// Assert
	assert.Nil(t, err)

	// Teardown
	_, _ = orm.ID(sysId).Update(&sys)
}
