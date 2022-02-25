package main

import (
	"casbin-auth-go/api"
	"casbin-auth-go/config"
	"casbin-auth-go/pkg/logr"
	"casbin-auth-go/pkg/valider"
	"casbin-auth-go/route"
	"flag"
	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
)

var port string

// @title Casbin Auth API
// @version 1.0
// @description Casbin Auth API
// @termsOfService https://github.com/pinkeyu7/casbin-auth-go
// @license.name MIT
// @license.url
func main() {
	// init http port
	flag.StringVar(&port, "port", "8080", "Initial port number")
	flag.Parse()

	// init config
	config.InitEnv()

	// init logger
	logr.InitLogger()

	// init validation
	valider.Init()

	// init driver
	_ = api.InitXorm()
	_ = api.InitRedis()
	_ = api.InitRedisCluster()

	// init gin router
	r := route.Init()

	// start server
	err := r.Run(":" + port)
	if err != nil {
		log.Println(err)
	}
}
