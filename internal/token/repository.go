package token

import "fmt"

type Cache interface {
	GetTokenIat(accId int) (float64, error)
	SetTokenIat(accId int, iat int64) error
}

func GetSysAccountTokenRedisKey(accId int) string {
	return fmt.Sprintf("sys_account:%v:token", accId)
}
