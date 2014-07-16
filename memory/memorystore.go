package memstore

import (
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"github.com/freehaha/token-auth"
	"time"
)

type MemoryTokenStore struct {
	tokens   map[string]*MemoryToken
	idTokens map[string]*MemoryToken
	salt     string
}

type MemoryToken struct {
	ExpireAt time.Time
	Token    string
	Id       string
}

func (t *MemoryToken) IsExpired() bool {
	return time.Now().After(t.ExpireAt)
}

func (t *MemoryToken) String() string {
	return t.Token
}

/* lookup 'exp' or 'id' */
func (t *MemoryToken) Claims(key string) interface{} {
	switch key {
	case "exp":
		return t.ExpireAt
	case "id":
		return t.Id
	default:
		return nil
	}
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

/* returns a new token with specific id */
func (s *MemoryTokenStore) NewToken(id interface{}) tauth.Token {
	strId := id.(string)
	bToken := s.generateToken(strId)
	strToken := base64.URLEncoding.EncodeToString(bToken)
	t := &MemoryToken{
		ExpireAt: time.Now().Add(time.Minute * 30),
		Token:    strToken,
		Id:       strId,
	}
	oldT, ok := s.idTokens[strId]
	if ok {
		delete(s.tokens, oldT.Token)
	}
	s.tokens[strToken] = t
	s.idTokens[strId] = t
	return t
}

/* Create a new memory store */
func New(salt string) *MemoryTokenStore {
	return &MemoryTokenStore{
		salt:     salt,
		tokens:   make(map[string]*MemoryToken),
		idTokens: make(map[string]*MemoryToken),
	}

}

func (s *MemoryTokenStore) CheckToken(strToken string) (tauth.Token, error) {
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
