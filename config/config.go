package config

import (
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const (
	EnvProduction          = "production"
	EnvStaging             = "staging"
	EnvDevelopment         = "development"
	EnvLocalhost           = "localhost"
	AdminUserId            = int64(-1)
	RedisDefaultExpireTime = time.Second * 60 * 60 * 24 * 30 // 預設一個月
)

var EnvShortName = map[string]string{
	EnvProduction:  "prod",
	EnvStaging:     "stag",
	EnvDevelopment: "dev",
	EnvLocalhost:   "local",
}

// Environment
func GetEnvironment() string {
	return os.Getenv("ENVIRONMENT")
}

func GetShortEnv() string {
	return EnvShortName[GetEnvironment()]
}

// Jwt Salt
func GetJwtSalt() string {
	return os.Getenv("JWT_SALT")
}

// Base path
var (
	_, b, _, _ = runtime.Caller(0)
	basePath   = filepath.Dir(b)
)

func GetBasePath() string {
	d := path.Join(path.Dir(b))
	return filepath.Dir(d)
}

// Cors
func GetCorsRule(origin string) bool {
	switch GetEnvironment() {
	case EnvLocalhost:
		return true
	case EnvDevelopment:
		return origin == "https://sample-development.website.com" || strings.Contains(origin, "http://localhost")
	case EnvStaging:
		return origin == "https://sample-staging.website.com"
	case EnvProduction:
		return origin == "https://sample.website.com"
	default:
		return true
	}
}

func InitEnv() {
	remoteBranch := os.Getenv("REMOTE_BRANCH")

	if remoteBranch == "" {
		// load env from .env file
		err := godotenv.Load()
		if err != nil {
			log.Panicln(err)
		}
	}
}
