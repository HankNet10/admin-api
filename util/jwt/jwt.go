package jwt

import (
	"errors"
	"fmt"
	"myadmin/config"

	"github.com/golang-jwt/jwt/v4"
)

type AccessToken struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int64  `json:"expires_in"`
}

var TokenType = "x-token"

func EncodeToken(sub string, exp int64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": sub,
		"exp": exp,
	})
	return token.SignedString([]byte(config.JwtSecret))
}

func DecondeToken(tokenString string) (map[string]interface{}, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok { //jwt验证的方法不对。
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.JwtSecret), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid { // 通常都是验证通过的。
		return nil, errors.New("token valid false")
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return claims, nil
	}
	return nil, errors.New("unexpected token  Claims")
}
