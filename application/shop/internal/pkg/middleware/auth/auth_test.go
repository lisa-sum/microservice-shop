package auth

import (
	jwt2 "github.com/golang-jwt/jwt/v4"
	"testing"
	"time"
)

var jwtToken = "hqFr%3ddt32DGlSTOI5cO6@TH#fFwYnP$S"

func TestCreateToken(t *testing.T) {
	// Create claims with multiple fields populated
	claims := CustomClaims{
		ID:          1,
		Nickname:    "213",
		AuthorityId: 1,
		RegisteredClaims: jwt2.RegisteredClaims{
			ExpiresAt: jwt2.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt2.NewNumericDate(time.Now()),
			NotBefore: jwt2.NewNumericDate(time.Now()),
			Subject:   "somebody",
			ID:        "1",
			Audience:  []string{"somebody_else"},
			Issuer:    "Gyl",
		},
		// StandardClaims: jwt2.StandardClaims{
		// 	NotBefore: time.Now().Unix() * 30, // 签名生效时间
		// 	ExpiresAt: time.Now().Unix() * 30, // 过期时间, 30天
		// 	Issuer:    "Gyl",
		// },
	}

	t.Logf("claims.ID:%d", claims.ID)
	t.Logf("claims.NickName:%s", claims.Nickname)
	t.Logf("claims.AuthorityId:%#v", claims.AuthorityId)

	token := jwt2.NewWithClaims(jwt2.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(jwtToken))
	if err != nil {
		t.Errorf("generate token err:%v", err)
	}
	t.Log(ss)
}
