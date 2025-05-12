package token

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/minhdang2803/simple_bank/utils"
	"github.com/stretchr/testify/require"
)

func TestJWTMaker(t *testing.T) {
	maker, err := NewJwtMaker(utils.RandomString(32))
	require.NoError(t, err)

	username := utils.RandomOwner()
	duration := time.Minute

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	token, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.NotZero(t, payload.ID)
	require.Equal(t, username, payload.Username)
	require.WithinDuration(t, issuedAt, payload.IssueAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)
}

func TestExpiredJwtToken(t *testing.T) {
	maker, err := NewJwtMaker(utils.RandomString(32))
	require.NoError(t, err)

	token, err := maker.CreateToken(utils.RandomOwner(), -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, err.Error(), jwt.ErrTokenInvalidClaims)
	require.Nil(t, payload)
}

// / Steps to remember:
// /
func TestInValidJwtTokenAlgNone(t *testing.T) {
	payload, err := NewPayload(utils.RandomOwner(), time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	jwtMaker := jwt.NewWithClaims(jwt.SigningMethodNone, payload)
	require.NoError(t, err)
	token, err := jwtMaker.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)

	maker, err := NewJwtMaker(utils.RandomString(32))

	payload, err = maker.VerifyToken(token)
	require.Error(t, err)
	require.Nil(t, payload)
}
