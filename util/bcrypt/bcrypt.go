package bcrypt

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

//输入明文，加密为密码
func GeneratePassword(pwd string) string {
	// bcrypt.CompareHashAndPassword()
	crypt, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.MinCost)
	if err != nil {
		log.Println(err)
		return ""
	}
	return string(crypt)
}

// pwd 明文密码
// hash 密文
func ComparePassword(pwd, hash string) bool {
	result := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pwd))
	return result == nil
}
