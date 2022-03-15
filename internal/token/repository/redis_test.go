package repository

import (
	"casbin-auth-go/config"
	"casbin-auth-go/driver"
	"casbin-auth-go/pkg/valider"
	"log"
	"os"
	"testing"
	"time"

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

func TestCache_GetTokenIat(t *testing.T) {
	// Arrange
	rc, _ := driver.NewRedis()
	tc := NewCache(rc)

	accId := 1
	newIat := int64(1594698129)
	_ = tc.SetTokenIat(accId, newIat)

	// Act
	iat, err := tc.GetTokenIat(accId)

	// Assert
	assert.Equal(t, nil, err)
	assert.Equal(t, float64(newIat), iat)
	if err != nil {
		t.Error(err)
	}
}

func TestCache_SetTokenIat(t *testing.T) {
	// Arrange
	rc, _ := driver.NewRedis()
	tc := NewCache(rc)

	accId := 1
	iat := time.Now().UTC().Unix()

	// Act
	err := tc.SetTokenIat(accId, iat)

	// Assert
	assert.Equal(t, nil, err)

	if err != nil {
		t.Error(err)
	}
}
