package api

import (
	"casbin-auth-go/driver"
	casbinPkg "casbin-auth-go/pkg/casbin"
	"log"

	"github.com/casbin/casbin/v2"

	"github.com/go-redis/redis/v7"
	"xorm.io/xorm"
)

type Env struct {
	Orm          *xorm.EngineGroup
	Redis        *redis.Client
	RedisCluster *redis.ClusterClient
	Casbin       *casbin.Enforcer
}

var env = &Env{}

func GetEnv() *Env {
	return env
}

func InitXorm() *xorm.EngineGroup {
	var err error

	env.Orm, err = driver.NewXorm()

	if err != nil {
		log.Println(err)
	}

	return env.Orm
}

func InitRedis() *redis.Client {
	var err error
	env.Redis, err = driver.NewSingleRedis()
	if err != nil {
		log.Println(err)
	}

	return env.Redis
}

func InitRedisCluster() *redis.ClusterClient {
	var err error
	env.RedisCluster, err = driver.NewRedis()
	if err != nil {
		log.Println(err)
	}

	return env.RedisCluster
}

func InitCasbin() {
	env.Casbin = casbinPkg.Init()
}
