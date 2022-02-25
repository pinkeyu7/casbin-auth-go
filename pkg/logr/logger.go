package logr

import (
	"casbin-auth-go/config"
	"go.uber.org/zap"
	"log"
)

var L *zap.Logger

func InitLogger() {
	var err error

	if config.GetEnvironment() == config.EnvProduction {
		L, err = zap.NewProduction()
	} else if config.GetEnvironment() == config.EnvLocalhost {
		L, err = zap.NewDevelopment()
	} else {
		config := zap.NewProductionConfig()
		config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
		L, err = config.Build()
	}

	defer L.Sync()

	if err != nil {
		log.Panicln("Init zap log failed...")
	}
}

func NewLogger() *zap.Logger {
	var err error
	var logger *zap.Logger
	if config.GetEnvironment() != config.EnvLocalhost {
		logger, err = zap.NewProduction()
	} else {
		logger, err = zap.NewDevelopment()
	}

	defer logger.Sync()

	if err != nil {
		log.Println("Init zap log failed...")
	}

	return logger
}
