package casbin

import (
	"casbin-auth-go/config"
	"casbin-auth-go/pkg/valider"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"testing"
)

func TestInit(t *testing.T) {
	e := Init()

	e.AddPolicy("admin", "/v1/test", "post")
	e.AddPolicy("admin", "/v1/test", "get")
	e.AddPolicy("admin", "/v1/test", "put")
	e.AddPolicy("admin", "/v1/test", "delete")
	e.AddRoleForUser("eric", "admin")
	ss, _ := e.GetRolesForUser("eric")
	gg, _ := e.GetRolesForUser("admin")
	fmt.Println(ss)
	fmt.Println(gg)

	sub := "eric"     // the user that wants to access a resource.
	obj := "/v1/test" // the resource that is going to be accessed.
	act := "post"     // the operation that the user performs on the resource.

	ok, _ := e.Enforce(sub, obj, act)
	fmt.Println(ok)
	sub = "eric"     // the user that wants to access a resource.
	obj = "/v2/test" // the resource that is going to be accessed.
	act = "post"     // the operation that the user performs on the resource.

	ok, _ = e.Enforce(sub, obj, act)

	e.AddPolicy("admin", "/v1/test", "post")
	e.AddPolicy("admin", "/v1/test", "get")
	e.AddPolicy("admin", "/v1/test", "put")
	e.AddPolicy("admin", "/v1/test", "delete")
	e.AddRoleForUser("eric", "admin")
}

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
