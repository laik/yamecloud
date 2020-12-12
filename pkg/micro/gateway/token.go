package gateway

import (
	"github.com/dgrijalva/jwt-go"
	flag "github.com/spf13/pflag"
	"time"
)

var salt string

func init() {
	flag.StringVar(&salt, "slat", "gw-private-se", "-slat secret_key")
}

// CustomClaims
type CustomClaims struct {
	UserName string `json:"user_name"`
	jwt.StandardClaims
}

// Token jwt service
type Token struct {
	privateKey []byte
}

//Decode
func (t *Token) Decode(tokenStr string) (*CustomClaims, error) {
	if len(t.privateKey) == 0 {
		t.privateKey = []byte(salt)
	}
	token, err := jwt.ParseWithClaims(
		tokenStr,
		&CustomClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return t.privateKey, nil
		},
	)
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, err
}

// Encode
func (t *Token) Encode(issuer, userName string, expireTime int64) (string, error) {
	if len(t.privateKey) == 0 {
		t.privateKey = []byte(salt)
	}
	return jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		CustomClaims{
			UserName: userName,
			StandardClaims: jwt.StandardClaims{
				Issuer:    issuer,
				IssuedAt:  time.Now().Unix(),
				ExpiresAt: expireTime,
			},
		}).SignedString(t.privateKey)
}
