package auth

import (
	"errors"
	"github.com/golang-jwt/jwt/v4"
)

type CustomClaims struct {
	ID          int64
	NickName    string
	AuthorityId int
	// jwt.StandardClaims
	jwt.RegisteredClaims
}

func CreateToken(c CustomClaims, key string) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodES512, c)
	signedToken, err := claims.SignedString([]byte(key))
	if err != nil {
		return "", errors.New("generate token fail" + err.Error())
	}
	return signedToken, err
}
