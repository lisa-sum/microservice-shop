package auth

import (
	"errors"
	"github.com/golang-jwt/jwt/v4"
)

type CustomClaims struct {
	ID          int64
	Nickname    string
	AuthorityId int
	// jwt.StandardClaims
	jwt.RegisteredClaims
}

func CreateToken(claims CustomClaims, key string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(key))
	if err != nil {
		return "", errors.New("generate token fail" + err.Error())
	}
	return signedToken, err
}
