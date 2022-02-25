package helper

import (
	"crypto/md5"
	"crypto/sha1"
	"fmt"
	"golang.org/x/crypto/scrypt"
	"io"
	"os"
)

func Sha1Str(str string) string {
	h := sha1.New()
	_, _ = h.Write([]byte(str))
	bs := h.Sum(nil)
	return fmt.Sprintf("%x", bs)
}

func ScryptStr(str string) string {
	salt := os.Getenv("SCRYPT_SALT")
	secret := []byte(salt)
	dk, _ := scrypt.Key([]byte(str), secret, 16384, 8, 1, 32)
	return fmt.Sprintf("%x", dk)
}

func Md5Str(str string) (string, error) {
	m := md5.New()
	_, err := io.WriteString(m, str)
	hash := fmt.Sprintf("%x", m.Sum(nil))
	return hash, err
}
