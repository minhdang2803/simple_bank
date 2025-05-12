package token

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	ErrInvalidToken = errors.New("token is invalid")
	ErrExpiredToken = errors.New("token has expired")
)

func getRole(role string) string {
	if role == "admin" {
		return "admin"
	}
	return "user"
}

type Payload struct {
	ID        uuid.UUID
	Username  string           `json:"username"`
	IssueAt   time.Time        `json:"issue_at"`
	ExpiredAt time.Time        `json:"expired_at"`
	Audience  jwt.ClaimStrings `json:"audience"`
}

// GetExpirationTime implements jwt.Claims.
func (payload *Payload) GetExpirationTime() (*jwt.NumericDate, error) {
	return jwt.NewNumericDate(payload.ExpiredAt), nil
}

// GetIssuedAt implements jwt.Claims.
func (payload *Payload) GetIssuedAt() (*jwt.NumericDate, error) {
	return jwt.NewNumericDate(payload.IssueAt), nil
}

// GetIssuer implements jwt.Claims.
func (payload *Payload) GetIssuer() (string, error) {
	return "simple_bank", nil
}

// GetNotBefore implements jwt.Claims.
func (payload *Payload) GetNotBefore() (*jwt.NumericDate, error) {
	return jwt.NewNumericDate(payload.IssueAt), nil
}

// GetSubject implements jwt.Claims.
func (payload *Payload) GetSubject() (string, error) {
	return payload.Username, nil
}

func NewPayload(username string, duration time.Duration) (*Payload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	payload := &Payload{
		ID:        tokenID,
		Username:  username,
		IssueAt:   time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}
	payload.Audience, _ = payload.GetAudience()
	return payload, nil
}

func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiredAt) {
		return jwt.ErrTokenExpired
	}
	return nil
}

func (payload *Payload) GetAudience() (jwt.ClaimStrings, error) {
	audience := []string{getRole(payload.Username)}
	return jwt.ClaimStrings(audience), nil
}
