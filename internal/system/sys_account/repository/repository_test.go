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

func TestRepository_InsertWithRole(t *testing.T) {
	// Arrange
	orm, _ := driver.NewXorm()
	sar := NewRepository(orm)

	sysRoleId := 1
	m := model.SysAccount{
		SystemId: 1,
		Account:  "test_insert_account",
		Phone:    "test_phone",
		Email:    "test_email",
		Password: "test_password",
		Name:     "test_name",
	}

	// Act
	acc, err := sar.InsertWithRole(&m, sysRoleId)

	// Assert
	assert.Nil(t, err)
	assert.NotNil(t, acc)

	// Teardown
	_, _ = orm.ID(acc.Id).Delete(&model.SysAccount{})
	_, _ = orm.Where("sys_account_id = ? ", acc.Id).Delete(&model.SysAccountRole{})
}

func TestRepository_UpdateWithRole(t *testing.T) {
	// Arrange
	orm, _ := driver.NewXorm()
	sar := NewRepository(orm)

	sysAccId := 3
	sysRoleId := 3
	acc := model.SysAccount{Id: sysAccId}
	_, _ = orm.Get(&acc)
	m := model.SysAccount{
		Id:        acc.Id,
		Phone:     "test_phone",
		Email:     "test_email",
		Name:      "test_name",
		IsDisable: false,
	}

	// Act
	err := sar.UpdateWithRole(&m, sysRoleId)

	// Assert
	assert.Nil(t, err)

	// Teardown
	_, _ = orm.ID(sysAccId).Update(&acc)
	_, _ = orm.Where("sys_account_id = ?", m.Id).Cols("sys_role_id").Update(&model.SysAccountRole{SysRoleId: 2})
}

func TestRepository_UpdatePassword(t *testing.T) {
	// Arrange
	orm, _ := driver.NewXorm()
	sar := NewRepository(orm)

	sysAccId := 1
	acc := model.SysAccount{Id: sysAccId}
	_, _ = orm.Get(&acc)

	// Act
	err := sar.UpdatePassword(sysAccId, "test_password")

	// Assert
	assert.Nil(t, err)

	// Teardown
	_, _ = orm.ID(sysAccId).Update(&acc)
}

func TestRepository_FindOne(t *testing.T) {
	// Arrange
	orm, _ := driver.NewXorm()
	sar := NewRepository(orm)

	// No data
	// Act
	sysAccId := 100
	acc, err := sar.FindOne(&model.SysAccount{Id: sysAccId})

	// Assert
	assert.Nil(t, err)
	assert.Nil(t, acc)

	// Has data
	// Act
	sysAccId = 1
	acc, err = sar.FindOne(&model.SysAccount{Id: sysAccId})

	// Assert
	assert.Nil(t, err)
	assert.NotNil(t, acc)
}

func TestRepository_FindOneWithRole(t *testing.T) {
	// Arrange
	orm, _ := driver.NewXorm()
	sar := NewRepository(orm)

	// No data
	// Act
	sysAccId := 100
	acc, err := sar.FindOneWithRole(sysAccId)

	// Assert
	assert.Nil(t, err)
	assert.Nil(t, acc)

	// Has data
	// Act
	sysAccId = 1
	acc, err = sar.FindOneWithRole(sysAccId)

	// Assert
	assert.Nil(t, err)
	assert.NotNil(t, acc)
}

func TestRepository_Find(t *testing.T) {
	// Arrange
	orm, _ := driver.NewXorm()
	sar := NewRepository(orm)

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
			2,
			10,
			0,
			3,
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
		t.Run(fmt.Sprintf("Find Sys Account,SystemId:%d,Offset:%d,Limit:%d", tc.SystemId, tc.Offset, tc.Limit), func(t *testing.T) {
			list, err := sar.Find(tc.SystemId, tc.Offset, tc.Limit)
			assert.Nil(t, err)
			assert.Len(t, list, tc.WantCount)
			for i, item := range list {
				fmt.Println(i, item, item.SysRoleId)
			}
		})
	}
}

func TestRepository_Count(t *testing.T) {
	// Arrange
	orm, _ := driver.NewXorm()
	sar := NewRepository(orm)

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
		t.Run(fmt.Sprintf("Count Sys Account,SystemId:%d", tc.SystemId), func(t *testing.T) {
			count, err := sar.Count(tc.SystemId)
			assert.Nil(t, err)
			assert.Equal(t, tc.WantCount, count)
		})
	}
}
