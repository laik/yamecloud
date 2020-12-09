package token

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/micro/go-micro/config"
	"github.com/micro/go-micro/config/source"
	"github.com/micro/go-micro/util/log"
	"sync"
	"time"
)

// CustomClaims
type CustomClaims struct {
	UserName string `json:"user_name"`
	jwt.StandardClaims
}

// Token jwt service
type Token struct {
	config.Config
	rock       sync.RWMutex
	privateKey []byte
}

func (t *Token) get() []byte {
	t.rock.RLock()
	defer func() {
		t.rock.RUnlock()
	}()
	return t.privateKey
}

func (t *Token) put(newKey []byte) {
	t.rock.Lock()
	defer func() {
		t.rock.Unlock()
	}()
	t.privateKey = newKey
}

func NewToken(source source.Source, path ...string) (*Token, error) {
	token := &Token{
		Config: config.NewConfig(),
	}
	err := token.Load(source)
	if err != nil {
		return nil, err
	}
	value := token.Get(path...).Bytes()
	if len(value) == 0 {
		return nil, fmt.Errorf("jwt key acquisition failed")
	}
	token.put(value)
	token.enableAutoUpdate(path...)

	return token, nil
}

func (t *Token) enableAutoUpdate(path ...string) error {
	watcher, err := t.Watch(path...)
	if err != nil {
		return err
	}
	go func() {
		for {
			value, err := watcher.Next()
			if err != nil {
				log.Error(err)
				continue
			}
			t.put(value.Bytes())
		}
	}()
	return nil
}

//Decode
func (t *Token) Decode(tokenStr string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenStr,
		&CustomClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return t.get(), nil
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
	return jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		CustomClaims{
			UserName: userName,
			StandardClaims: jwt.StandardClaims{
				Issuer:    issuer,
				IssuedAt:  time.Now().Unix(),
				ExpiresAt: expireTime,
			},
		}).SignedString(t.get())
}
