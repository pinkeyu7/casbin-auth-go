package apires

import "time"

type SysAccountToken struct {
	Token     string                 `json:"token"`
	ExpiredAt time.Time              `json:"expired_at"`
	Data      map[string]interface{} `json:"data"`
}
