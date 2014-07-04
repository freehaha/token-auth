package tauth

import (
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/gorilla/context"
	"net/http"
	"time"
)

type TokenAuth struct {
	handler             http.Handler
	store               TokenStore
	UnauthorizedHandler http.HandlerFunc
}

type TokenAuthNegroni struct {
	TokenAuth
}

type TokenStore interface {
	NewToken(id string) *Token
	CheckToken(token string) (*Token, error)
}

type MemoryTokenStore struct {
	tokens map[string]*Token
	salt   string
}

type Token struct {
	ExpireAt time.Time
	Token    string
	Id       string
}

func DefaultUnauthorizedHandler(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(401)
	fmt.Fprint(w, "unauthorized")
}

func (s *MemoryTokenStore) generateToken(id string) []byte {
	hash := sha1.New()
	now := time.Now()
	timeStr := now.Format(time.ANSIC)
	hash.Write([]byte(timeStr))
	hash.Write([]byte(id))
	hash.Write([]byte("salt"))
	return hash.Sum(nil)
}

func (s *MemoryTokenStore) NewToken(id string) *Token {
	bToken := s.generateToken(id)
	strToken := base64.URLEncoding.EncodeToString(bToken)
	t := &Token{
		ExpireAt: time.Now().Add(time.Minute * 30),
		Token:    strToken,
		Id:       id,
	}
	s.tokens[strToken] = t
	return t
}

func NewMemoryTokenStore(salt string) *MemoryTokenStore {
	return &MemoryTokenStore{
		salt:   salt,
		tokens: make(map[string]*Token),
	}

}

func (s *MemoryTokenStore) CheckToken(strToken string) (*Token, error) {
	t, ok := s.tokens[strToken]
	if !ok {
		return nil, errors.New("Failed to authenticate")
	}
	if t.ExpireAt.Before(time.Now()) {
		delete(s.tokens, strToken)
		return nil, errors.New("Token expired")
	}
	return t, nil
}

/*
	Returns a TokenAuth object implemting Handler interface

	if a handler is given it proxies the request to the handler

	if a unauthorizedHandler is provided, unauthorized requests will be handled by this HandlerFunc,
	otherwise a default unauthorized handler is used.

	store is the TokenStore that stores and verify the tokens
*/
func NewTokenAuth(handler http.Handler, unauthorizedHandler http.HandlerFunc, store TokenStore) *TokenAuth {
	t := &TokenAuth{
		handler:             handler,
		store:               store,
		UnauthorizedHandler: unauthorizedHandler,
	}
	if t.UnauthorizedHandler == nil {
		t.UnauthorizedHandler = DefaultUnauthorizedHandler
	}
	return t
}

/* wrap a HandlerFunc to be authenticated */
func (t *TokenAuth) HandleFunc(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		token, err := t.authenticate(req)
		if err != nil {
			t.UnauthorizedHandler.ServeHTTP(w, req)
			return
		}
		context.Set(req, "token", token)
		handlerFunc.ServeHTTP(w, req)
	}
}

func (t *TokenAuth) authenticate(req *http.Request) (*Token, error) {
	strToken := req.URL.Query().Get("token")
	if strToken == "" {
		return nil, errors.New("token required")
	}
	token, err := t.store.CheckToken(strToken)
	if err != nil {
		return nil, errors.New("Invalid token")
	}
	return token, nil
}

/* implement Handler */
func (t *TokenAuth) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	token, err := t.authenticate(req)
	if err != nil {
		t.UnauthorizedHandler.ServeHTTP(w, req)
		return
	}
	context.Set(req, "token", token)
	t.handler.ServeHTTP(w, req)
}

func Get(req *http.Request) *Token {
	return context.Get(req, "token").(*Token)
}
