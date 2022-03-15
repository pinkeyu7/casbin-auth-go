package helper

import (
	"casbin-auth-go/pkg/er"
	"net/http"
	"regexp"
)

func PasswordValidation(pw string) error {
	// 驗證密碼
	// \d 比對数。 MustCompile必须
	var rNum = regexp.MustCompile(`\d`)

	// [a-zA-Z] 比對大小寫字母
	var rCharacter = regexp.MustCompile("[a-zA-Z]")

	// 比對空格
	var rBlank = regexp.MustCompile(" ")
	// 取数字的个数
	num := len(rNum.FindAllStringSubmatch(pw, -1))

	// 取字母的個数
	character := len(rCharacter.FindAllStringSubmatch(pw, -1))

	// 取空格的個数
	blank := len(rBlank.FindAllStringSubmatch(pw, -1))

	if num < 1 || character < 1 || blank != 0 {
		pwErr := er.NewAppErr(http.StatusBadRequest, er.ErrorParamInvalid, "Contain at least one character ,one number and without blank", nil)
		return pwErr
	}
	return nil
}
