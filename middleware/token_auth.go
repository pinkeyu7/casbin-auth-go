package middleware

import (
	"casbin-auth-go/api"
	tokenLibrary "casbin-auth-go/internal/token/library"
	tokenRepo "casbin-auth-go/internal/token/repository"
	"casbin-auth-go/pkg/er"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func TokenAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("Bearer")
		if token == "" {
			authErr := er.NewAppErr(http.StatusUnauthorized, er.UnauthorizedError, "token is required.", nil)
			c.AbortWithStatusJSON(authErr.GetStatus(), authErr.GetMsg())
			return
		}

		claims, err := tokenLibrary.ParseToken(token)
		if err != nil {
			authErr := er.NewAppErr(http.StatusUnauthorized, er.UnauthorizedError, "token is not valid.", err)
			c.AbortWithStatusJSON(authErr.GetStatus(), authErr.GetMsg())
			return
		}

		var jwtAccountId, parseAccountIdOk = claims["account_id"].(string)
		var jwtIat, jwtIatOk = claims["iat"].(float64)

		if !parseAccountIdOk || !jwtIatOk {
			parseJwtInfoErr := er.NewAppErr(http.StatusUnauthorized, er.UnauthorizedError, "token is not valid.", nil)
			c.AbortWithStatusJSON(parseJwtInfoErr.GetStatus(), parseJwtInfoErr.GetMsg())
			return
		}

		accId, _ := strconv.Atoi(jwtAccountId)

		// Jwt token state management
		env := api.GetEnv()
		tc := tokenRepo.NewCache(env.RedisCluster)
		serverIat, _ := tc.GetTokenIat(accId)
		if jwtIat < serverIat {
			iatErr := er.NewAppErr(http.StatusUnauthorized, er.UnauthorizedError, "token is expired.", nil)
			c.AbortWithStatusJSON(iatErr.GetStatus(), iatErr.GetMsg())
			return
		}

		// Check casbin permission
		_ = env.Casbin.LoadPolicy()
		access, err := env.Casbin.Enforce(jwtAccountId, c.Request.URL.Path, c.Request.Method)
		if err != nil {
			permissionErr := er.NewAppErr(http.StatusInternalServerError, er.UnknownError, "permission error.", err)
			c.AbortWithStatusJSON(permissionErr.GetStatus(), permissionErr.GetMsg())
			return
		}
		if !access {
			accessDeniedErr := er.NewAppErr(http.StatusForbidden, er.ForbiddenError, "access denied.", nil)
			c.AbortWithStatusJSON(accessDeniedErr.GetStatus(), accessDeniedErr.GetMsg())
			return
		}

		// Set claims
		c.Set("claims", claims)
		c.Set("account_id", accId)

		c.Next()
	}
}
