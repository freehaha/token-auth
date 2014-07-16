package jwtstore

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/freehaha/token-auth"
	"time"
)

type JwtStore struct {
	tokenKey []byte
}

type JwtToken struct {
	tokenKey []byte
	jwt.Token
}

func (t *JwtToken) Claims(key string) interface{} {
	return t.Token.Claims[key]
}

func (t *JwtToken) IsExpired() bool {
	exp := time.Unix(t.Claims("exp").(int64), 0)
	return time.Now().After(exp)
}

func (t *JwtToken) String() string {
	tokenStr, _ := t.Token.SignedString(t.tokenKey)
	return tokenStr
}

func (s *JwtStore) NewToken(id interface{}) tauth.Token {
	token := jwt.New(jwt.GetSigningMethod("HS256"))
	token.Claims["id"] = id.(string)
	token.Claims["exp"] = time.Now().Add(time.Minute)
	t := &JwtToken{
		tokenKey: s.tokenKey,
		Token:    *token,
	}
	return t
}

func (s *JwtStore) CheckToken(token string) (tauth.Token, error) {
	t, err := jwt.Parse(token, func(token *jwt.Token) ([]byte, error) {
		return s.tokenKey, nil
	})
	if err != nil {
		return nil, err
	}
	return &JwtToken{s.tokenKey, *t}, nil
}

func New(tokenKey string) *JwtStore {
	return &JwtStore{
		[]byte(tokenKey),
	}
}
