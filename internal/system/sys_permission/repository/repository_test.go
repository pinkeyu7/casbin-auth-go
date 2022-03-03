package repository

import (
	"casbin-auth-go/config"
	"casbin-auth-go/driver"
	"casbin-auth-go/dto/model"
	"casbin-auth-go/pkg/valider"
	"fmt"
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
	spr := NewRepository(orm)

	m := model.SysPermission{
		SystemId:     1,
		AllowApiPath: "/test_path/",
		Action:       "get",
		Slug:         "test_slug",
		Description:  "test_description",
	}

	// Act
	err := spr.Insert(&m)

	// Assert
	assert.Nil(t, err)

	// Teardown
	_, _ = orm.Where("slug = ? ", m.Slug).Delete(&model.SysPermission{})
}

func TestRepository_Update(t *testing.T) {
	// Arrange
	orm, _ := driver.NewXorm()
	spr := NewRepository(orm)

	perm := model.SysPermission{Id: 1}
	_, _ = orm.Get(&perm)
	m := model.SysPermission{
		Id:          perm.Id,
		Slug:        "test_slug",
		Description: "test_description",
	}

	// Act
	err := spr.Update(&m)

	// Assert
	assert.Nil(t, err)

	// Teardown
	_, _ = orm.ID(perm.Id).Update(&perm)
}

func TestRepository_Delete(t *testing.T) {
	// Arrange
	orm, _ := driver.NewXorm()
	spr := NewRepository(orm)

	m := model.SysPermission{
		SystemId:     1,
		AllowApiPath: "/test_path/",
		Action:       "get",
		Slug:         "test_slug",
		Description:  "test_description",
	}
	_, _ = orm.Insert(&m)

	// Act
	err := spr.Delete(m.Id)

	// Assert
	assert.Nil(t, err)
}

func TestRepository_FindOne(t *testing.T) {
	// Arrange
	orm, _ := driver.NewXorm()
	spr := NewRepository(orm)

	// No data
	// Act
	sysPermId := 100
	perm, err := spr.FindOne(&model.SysPermission{Id: sysPermId})

	// Assert
	assert.Nil(t, err)
	assert.Nil(t, perm)

	// Has data
	// Act
	sysPermId = 1
	perm, err = spr.FindOne(&model.SysPermission{Id: sysPermId})

	// Assert
	assert.Nil(t, err)
	assert.NotNil(t, perm)
	assert.Equal(t, sysPermId, perm.Id)
}

func TestRepository_FindByIds(t *testing.T) {
	// Arrange
	orm, _ := driver.NewXorm()
	spr := NewRepository(orm)

	// No data
	// Act
	ids := []int{100, 101, 102}
	list, err := spr.FindByIds(ids)

	// Assert
	assert.Nil(t, err)
	assert.Len(t, list, 0)

	// Has data
	// Act
	ids = []int{1, 2, 3, 100, 101, 102}
	list, err = spr.FindByIds(ids)

	// Assert
	assert.Nil(t, err)
	assert.NotNil(t, list)
	assert.Len(t, list, 3)
}

func TestRepository_Find(t *testing.T) {
	// Arrange
	orm, _ := driver.NewXorm()
	spr := NewRepository(orm)

	// Act
	testCases := []struct {
		SystemId  int
		Limit     int
		Offset    int
		WantCount int
	}{
		{
			0,
			1,
			0,
			1,
		},
		{
			0,
			10,
			0,
			6,
		},
		{
			1,
			10,
			0,
			2,
		},
		{
			100,
			10,
			0,
			0,
		},
	}
	// Act
	for _, tc := range testCases {
		tc := tc
		t.Run(fmt.Sprintf("Find Sys Permission,SystemId:%d,Offset:%d,Limit:%d", tc.SystemId, tc.Offset, tc.Limit), func(t *testing.T) {
			list, err := spr.Find(tc.SystemId, tc.Offset, tc.Limit)
			assert.Nil(t, err)
			assert.Len(t, list, tc.WantCount)
		})
	}
}

func TestRepository_Count(t *testing.T) {
	// Arrange
	orm, _ := driver.NewXorm()
	spr := NewRepository(orm)

	// Act
	testCases := []struct {
		SystemId  int
		WantCount int
	}{
		{
			0,
			6,
		},
		{
			1,
			2,
		},
		{
			100,
			0,
		},
	}
	// Act
	for _, tc := range testCases {
		tc := tc
		t.Run(fmt.Sprintf("Count Sys Permission,SystemId:%d", tc.SystemId), func(t *testing.T) {
			count, err := spr.Count(tc.SystemId)
			assert.Nil(t, err)
			assert.Equal(t, tc.WantCount, count)
		})
	}
}

func TestRepository_Exist(t *testing.T) {
	// Arrange
	orm, _ := driver.NewXorm()
	spr := NewRepository(orm)

	// No data
	// Act
	sysPermId := 100
	exist, err := spr.Exist(&model.SysPermission{Id: sysPermId})

	// Assert
	assert.Nil(t, err)
	assert.False(t, exist)

	// Has data
	// Act
	sysPermId = 1
	exist, err = spr.Exist(&model.SysPermission{Id: sysPermId})

	// Assert
	assert.Nil(t, err)
	assert.True(t, exist)
}

func TestRepository_FindIdsByRole(t *testing.T) {
	// Arrange
	orm, _ := driver.NewXorm()
	spr := NewRepository(orm)

	// No data
	// Act
	sysRoleId := 100
	list, err := spr.FindIdsByRole(sysRoleId)

	// Assert
	assert.Nil(t, err)
	assert.Len(t, list, 0)

	// Has data
	// Act
	sysRoleId = 1
	list, err = spr.FindIdsByRole(sysRoleId)

	// Assert
	assert.Nil(t, err)
	assert.Len(t, list, 2)
}
