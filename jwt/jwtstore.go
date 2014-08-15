package jwtstore

import (
	"errors"
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

func (t *JwtToken) SetClaim(key string, value interface{}) tauth.ClaimSetter {
	t.Token.Claims[key] = value
	return t
}

func (t *JwtToken) IsExpired() bool {
	/* converted to float64 when parsed from JSON */
	exp := time.Unix(int64(t.Claims("exp").(float64)), 0)
	return time.Now().After(exp)
}

func (t *JwtToken) String() string {
	tokenStr, _ := t.Token.SignedString(t.tokenKey)
	return tokenStr
}

func (s *JwtStore) NewToken(id interface{}) *JwtToken {
	token := jwt.New(jwt.GetSigningMethod("HS256"))
	token.Claims["exp"] = time.Now().Add(time.Minute * 30).Unix()
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
	jtoken := &JwtToken{s.tokenKey, *t}
	if jtoken.IsExpired() {
		return nil, errors.New("Token expired")
	}
	return jtoken, nil
}

func New(tokenKey string) *JwtStore {
	return &JwtStore{
		[]byte(tokenKey),
	}
}
