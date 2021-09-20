package data

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"time"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

const (
	ScopeActivation = "activation"
)

type Token struct {
	Hash      []byte `gorm:"primaryKey"`
	Plaintext string
	UserID    uint
	User      User
	Expiry    time.Time
	Scope     string `gorm:"not null"`
}

func (t *Token) Validate() error {
	validate := validator.New()
	return validate.Struct(t)
}

func generateToken(userID uint, ttl time.Duration, scope string) (*Token, error) {
	token := &Token{
		UserID: userID,
		Expiry: time.Now().Add(ttl),
		Scope:  scope,
	}

	randomBytes := make([]byte, 16)

	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, err
	}

	token.Plaintext = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)
	hash := sha256.Sum256([]byte(token.Plaintext))
	token.Hash = hash[:]

	return token, nil
}

type TokenModel struct {
	DB *gorm.DB
}

func (m *TokenModel) New(userID uint, ttl time.Duration, scope string) (*Token, error) {
	token, err := generateToken(userID, ttl, scope)
	if err != nil {
		return nil, err
	}

	err = m.Insert(token)
	return token, err
}

func (m *TokenModel) Insert(token *Token) error {
	return m.DB.Create(token).Error
}

func (m *TokenModel) DeleteAllForUser(scope string, userID uint) error {
	return m.DB.Delete(&Token{UserID: userID, Scope: scope}).Error
}
