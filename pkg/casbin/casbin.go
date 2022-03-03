package casbin

import (
	"casbin-auth-go/dto/model"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/casbin/casbin/v2"
	casbinModel "github.com/casbin/casbin/v2/model"
	xormadapter "github.com/casbin/xorm-adapter/v2"
	_ "github.com/go-sql-driver/mysql"
)

func Init() *casbin.Enforcer {
	adapter, err := NewAdapter()
	if err != nil {
		log.Panicln(err)
	}
	m := casbinModel.NewModel()
	m.AddDef("r", "r", "sub, obj, act")
	m.AddDef("p", "p", "sub, obj, act")
	m.AddDef("g", "g", "_, _")
	m.AddDef("e", "e", "some(where (p.eft == allow))")
	m.AddDef("m", "m", "g(r.sub, p.sub) && keyMatch2(r.obj, p.obj) && regexMatch(r.act, p.act)")
	e, err := casbin.NewEnforcer(m, adapter)

	if err != nil {
		log.Panicln(err)
	}

	// Load the policy from DB.
	fmt.Println("Start Load Policy:", time.Now())
	err = e.LoadPolicy()

	if err != nil {
		log.Panicln(err)
	}

	return e
}

func NewAdapter() (*xormadapter.Adapter, error) {
	dsn := "%s:%s@(%s:%s)/%s?parseTime=true&charset=utf8mb4&collation=utf8mb4_unicode_ci"

	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	dbSource := fmt.Sprintf(dsn, dbUser, dbPassword, dbHost, dbPort, dbName)
	adapter, err := xormadapter.NewAdapter("mysql", dbSource, true)
	// Your driver and data source.

	return adapter, err
}

func AddUserRole(systemTag string, accountId, oldRoleId, newRoleId int64) (bool, error) {
	e := Init()

	// 先刪除後新增
	if oldRoleId != 0 {
		_, err := e.DeleteRoleForUser(strconv.FormatInt(accountId, 10), strconv.FormatInt(oldRoleId, 10))
		if err != nil {
			return false, err
		}
	}
	// e.AddRoleForUserInDomain("eric", "admin", "hc")
	added, err := e.AddRoleForUser(strconv.FormatInt(accountId, 10), strconv.FormatInt(newRoleId, 10))
	return added, err
}

func AddRolePolicy(systemTag string, roleId int, perms []*model.SysPermission) (bool, error) {
	e := Init()

	var added bool = true
	var err error

	// e.AddPolicy("admin", "/v1/test", "post", "hc")
	for i := range perms {
		if perms[i].AllowApiPath == "" || perms[i].Action == "" {
			continue
		}
		added, err = e.AddPolicy(roleId, perms[i].AllowApiPath, perms[i].Action)
	}
	return added, err
}

func RemoveRolePolicy(roleId int) (bool, error) {
	e := Init()
	removed, err := e.RemovePolicy(roleId)
	return removed, err
}
