package repository

import (
	"casbin-auth-go/config"
	"casbin-auth-go/driver"
	"casbin-auth-go/dto/apireq"
	"casbin-auth-go/dto/model"
	"casbin-auth-go/internal/system/system"
	"casbin-auth-go/pkg/valider"
	"fmt"
	"github.com/brianvoe/gofakeit/v4"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
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

func TestRepository_Insert(t *testing.T) {
	// Arrange
	orm, _ := driver.NewXorm()
	re, _ := driver.NewRedis()

	sc := NewCache(re)
	sr := NewRepository(orm, sc)

	m := apireq.AddSystem{
		AccountId:        1,
		Name:             "test admin permission",
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

	// Act
	err := sr.Insert(&m)

	// Assert
	assert.Nil(t, err)

	// Teardown
	sys := model.System{Name: m.Name}
	_, _ = orm.Get(&sys)
	_, _ = orm.Where("name = ? ", m.Name).Delete(&model.System{})
	_, _ = orm.Where("system_id = ? ", sys.Id).Delete(&model.SysRole{})
}

func TestRepository_Find(t *testing.T) {
	// Arrange
	orm, _ := driver.NewXorm()
	re, _ := driver.NewRedis()

	sc := NewCache(re)
	sr := NewRepository(orm, sc)

	// Act
	testCases := []struct {
		ListType  string
		Limit     int
		Offset    int
		WantCount int
	}{
		{
			system.ListTypeAll,
			1,
			0,
			1,
		},
		{
			system.ListTypeAll,
			10,
			0,
			2,
		},
		{
			system.ListTypeEnable,
			10,
			0,
			2,
		},
		{
			system.ListTypeDisable,
			10,
			0,
			0,
		},
	}
	// Act
	for _, tc := range testCases {
		tc := tc
		t.Run(fmt.Sprintf("Find System,ListType:%s,Offset:%d,Limit:%d", tc.ListType, tc.Offset, tc.Limit), func(t *testing.T) {
			data, err := sr.Find(tc.ListType, tc.Offset, tc.Limit)
			assert.Nil(t, err)
			assert.Len(t, data, tc.WantCount)
		})
	}
}

func TestRepository_FindOne(t *testing.T) {
	// Arrange
	orm, _ := driver.NewXorm()
	re, _ := driver.NewRedis()

	sc := NewCache(re)
	sr := NewRepository(orm, sc)

	// No data
	// Act
	res, err := sr.FindOne(&model.System{Id: 100})

	// Assert
	assert.Nil(t, err)
	assert.Nil(t, res)

	// Has data
	// Act
	res, err = sr.FindOne(&model.System{Id: 1})

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, 1, res.Id)
}

func TestRepository_Count(t *testing.T) {
	// Arrange
	orm, _ := driver.NewXorm()
	re, _ := driver.NewRedis()

	sc := NewCache(re)
	sr := NewRepository(orm, sc)

	// Act
	testCases := []struct {
		ListType  string
		WantCount int
	}{
		{
			system.ListTypeAll,
			2,
		},
		{
			system.ListTypeEnable,
			2,
		},
		{
			system.ListTypeDisable,
			0,
		},
	}
	// Act
	for _, tc := range testCases {
		tc := tc
		t.Run(fmt.Sprintf("Count System,ListType:%s", tc.ListType), func(t *testing.T) {
			count, err := sr.Count(tc.ListType)
			assert.Nil(t, err)
			assert.Equal(t, count, tc.WantCount)
		})
	}
}

func TestRepository_Update(t *testing.T) {
	// Arrange
	orm, _ := driver.NewXorm()
	re, _ := driver.NewRedis()

	sc := NewCache(re)
	sr := NewRepository(orm, sc)

	sys := model.System{Id: 1}
	_, _ = orm.Get(&sys)

	m := model.System{
		Id:         sys.Id,
		Name:       "test admin permission",
		Address:    gofakeit.Address().Address,
		Tel:        gofakeit.Phone(),
		IpAddress:  "*",
		MacAddress: "*",
		IsDisable:  false,
	}

	// Act
	err := sr.Update(&m)

	// Assert
	assert.Nil(t, err)

	// Teardown
	_, _ = orm.ID(sys.Id).Update(&sys)
}
