package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestPasswor(t *testing.T) {
	password := RandomString(6)

	hashedPassword, err := HashedPassword(password)

	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword)

	err = CheckPassword(password, hashedPassword)
	require.NoError(t, err)

	wrongPassword := RandomString(6)

	err = CheckPassword(wrongPassword, hashedPassword)
	require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())
}
