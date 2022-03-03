package route

import (
	"casbin-auth-go/config"
	_ "casbin-auth-go/docs"
	"casbin-auth-go/middleware"
	"casbin-auth-go/pkg/request_cache"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"os"
	"time"

	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

func Init() *gin.Engine {
	r := gin.New()

	// gin 檔案上傳body限制
	r.MaxMultipartMemory = 64 << 20 // 8 MiB

	// Request cache
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	addr := fmt.Sprintf("%s:%s", redisHost, redisPort)
	store := request_cache.NewRedisCache(addr, "", time.Second)

	// Middleware
	r.Use(middleware.LogRequest())
	r.Use(middleware.ErrorResponse())

	// Swagger
	if mode := gin.Mode(); mode == gin.DebugMode {
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	corsConf := cors.DefaultConfig()
	corsConf.AllowCredentials = true
	corsConf.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"}
	corsConf.AllowHeaders = []string{"Origin", "X-Requested-With", "Content-Type", "Accept", "Authorization", "Bearer", "Accept-Language"}
	corsConf.AllowOriginFunc = config.GetCorsRule
	r.Use(cors.New(corsConf))

	TokenV1(r, store)
	SystemV1(r)
	SysPermissionV1(r)
	SysRoleV1(r)

	return r
}
