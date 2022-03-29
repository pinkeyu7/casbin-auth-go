package token_library

import (
	"casbin-auth-go/config"
	"casbin-auth-go/pkg/er"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func parseToken(tokenStr, salt string) (jwt.MapClaims, error) {
	// Parse takes the token string and a function for looking up the key. The latter is especially
	// useful if you use multiple keys for your application.  The standard is to use 'kid' in the
	// head of the token to identify which key to use, but the parsed token (head and claims) is provided
	// to the callback, providing flexibility.
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(salt), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, err
}

func ParseClaims(claims interface{}) jwt.MapClaims {
	if c, ok := claims.(jwt.MapClaims); ok {
		return c
	}

	return nil
}

// ---------------------------------------- JWT Token Generation ----------------------------------------------

func GenToken(accId int) (string, time.Time, error) {
	salt := config.GetJwtSalt()
	secret := []byte(salt)

	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	exp := time.Now().Add(time.Hour * 24).UTC()

	accIdStr := strconv.Itoa(accId)
	token := jwt.NewWithClaims(jwt.SigningMethodHS384, jwt.MapClaims{
		"iss":        "casbin-auth-go",
		"exp":        exp.Unix(),              // Expiration Time,
		"iat":        time.Now().UTC().Unix(), // Issued At Time
		"account_id": accIdStr,
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(secret)

	return tokenString, exp, err
}

// ---------------------------------------- Middleware JWT Token 驗證 ----------------------------------------------

func ParseToken(tokenStr string) (jwt.MapClaims, error) {
	salt := config.GetJwtSalt()
	return parseToken(tokenStr, salt)
}

// ---------------------------------------- JWT Account Id 驗證 ----------------------------------------------

func CheckJWTAccountId(c *gin.Context, accId int) error {
	claims, _ := c.Get("claims")
	claimsMap := ParseClaims(claims)

	if strconv.Itoa(accId) != claimsMap["account_id"] {
		authErr := er.NewAppErr(http.StatusUnauthorized, er.UnauthorizedError, "token is not valid.", nil)
		return authErr
	}

	return nil
}
