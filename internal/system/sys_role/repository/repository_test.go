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

func TestRepository_InsertWithPermission(t *testing.T) {
	// Arrange
	orm, _ := driver.NewXorm()
	srr := NewRepository(orm)

	m := model.SysRole{
		Sort:        5,
		SystemId:    1,
		Name:        "test_name",
		DisplayName: "test_display_name",
	}

	// Act
	err := srr.InsertWithPermission(&m, nil)

	// Assert
	assert.Nil(t, err)

	// Teardown
	_, _ = orm.Where("name = ? ", m.Name).Delete(&model.SysRole{})
}

func TestRepository_UpdateWithPermission(t *testing.T) {
	// Arrange
	orm, _ := driver.NewXorm()
	srr := NewRepository(orm)

	role := model.SysRole{Id: 1}
	_, _ = orm.Get(&role)
	m := model.SysRole{
		Id:          role.Id,
		Name:        "test_name",
		DisplayName: "test_display_name",
		IsDisable:   false,
	}

	// Act
	err := srr.UpdateWithPermission(&m, nil)

	// Assert
	assert.Nil(t, err)

	// Teardown
	_, _ = orm.ID(role.Id).Update(&role)
}

func TestRepository_FindOne(t *testing.T) {
	// Arrange
	orm, _ := driver.NewXorm()
	srr := NewRepository(orm)

	// No data
	// Act
	sysRoleId := 100
	perm, err := srr.FindOne(&model.SysRole{Id: sysRoleId})

	// Assert
	assert.Nil(t, err)
	assert.Nil(t, perm)

	// Has data
	// Act
	sysRoleId = 1
	perm, err = srr.FindOne(&model.SysRole{Id: sysRoleId})

	// Assert
	assert.Nil(t, err)
	assert.NotNil(t, perm)
	assert.Equal(t, sysRoleId, perm.Id)
}

func TestRepository_Find(t *testing.T) {
	// Arrange
	orm, _ := driver.NewXorm()
	srr := NewRepository(orm)

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
			4,
		},
		{
			1,
			10,
			0,
			1,
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
		t.Run(fmt.Sprintf("Find Sys Role,SystemId:%d,Offset:%d,Limit:%d", tc.SystemId, tc.Offset, tc.Limit), func(t *testing.T) {
			list, err := srr.Find(tc.SystemId, tc.Offset, tc.Limit)
			assert.Nil(t, err)
			assert.Len(t, list, tc.WantCount)
		})
	}
}

func TestRepository_Count(t *testing.T) {
	// Arrange
	orm, _ := driver.NewXorm()
	srr := NewRepository(orm)

	// Act
	testCases := []struct {
		SystemId  int
		WantCount int
	}{
		{
			0,
			4,
		},
		{
			1,
			1,
		},
		{
			2,
			3,
		},
		{
			100,
			0,
		},
	}
	// Act
	for _, tc := range testCases {
		tc := tc
		t.Run(fmt.Sprintf("Count Sys Role,SystemId:%d", tc.SystemId), func(t *testing.T) {
			count, err := srr.Count(tc.SystemId)
			assert.Nil(t, err)
			assert.Equal(t, tc.WantCount, count)
		})
	}
}

func TestRepository_FindBySysId(t *testing.T) {
	// Arrange
	orm, _ := driver.NewXorm()
	srr := NewRepository(orm)

	// Act
	testCases := []struct {
		SystemId  int
		WantCount int
	}{
		{
			100,
			0,
		},
		{
			1,
			1,
		},
		{
			2,
			3,
		},
	}
	// Act
	for _, tc := range testCases {
		tc := tc
		t.Run(fmt.Sprintf("Find Sys Role By Sys Id,SystemId:%d", tc.SystemId), func(t *testing.T) {
			list, err := srr.FindBySysId(tc.SystemId)
			assert.Nil(t, err)
			assert.Len(t, list, tc.WantCount)
		})
	}
}

func TestRepository_FindByPermId(t *testing.T) {
	// Arrange
	orm, _ := driver.NewXorm()
	srr := NewRepository(orm)

	// Act
	testCases := []struct {
		SysPermId int
		WantCount int
	}{
		{
			100,
			0,
		},
		{
			1,
			1,
		},
	}
	// Act
	for _, tc := range testCases {
		tc := tc
		t.Run(fmt.Sprintf("Find Sys Role By Perm Id,SysPermId:%d", tc.SysPermId), func(t *testing.T) {
			list, err := srr.FindBySysId(tc.SysPermId)
			assert.Nil(t, err)
			assert.Len(t, list, tc.WantCount)
		})
	}
}
