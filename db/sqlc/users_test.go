package db

import (
	"context"
	"testing"

	"github.com/minhdang2803/simple_bank/utils"
	"github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {
	createUser, err := CreateUser()
	require.NoError(t, err)

	user, err := testQueries.CreateUser(context.Background(), *createUser)

	require.NoError(t, err)
	require.NotEmpty(t, user)

	//
	require.Equal(t, createUser.Username, user.Username)
	require.Equal(t, createUser.HashedPassword, user.HashedPassword)
	require.Equal(t, createUser.FullName, user.FullName)
	require.Equal(t, createUser.Email, user.Email)

	//
	require.True(t, user.PasswordChangedAt.IsZero())
	require.NotZero(t, user.CreatedAt)
}

func CreateUser() (*CreateUserParams, error) {
	hashedPassword, err := utils.HashedPassword(utils.RandomString(6))
	if err != nil {
		return nil, err
	}
	return &CreateUserParams{
		Username:       utils.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName:       utils.RandomOwner(),
		Email:          utils.RandomEmail(),
	}, nil
}

func CreateRandomUser() (*User, error) {

	createUser, _ := CreateUser()
	user, err := testQueries.CreateUser(context.Background(), *createUser)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func TestGetUser(t *testing.T) {
	user, err := CreateRandomUser()
	require.NoError(t, err)

	retreivedUser, err := testQueries.GetUser(context.Background(), user.Username)

	require.NoError(t, err)
	require.NotEmpty(t, retreivedUser)
	require.Equal(t, user.Username, retreivedUser.Username)
	require.Equal(t, user.FullName, retreivedUser.FullName)
	require.Equal(t, user.Email, retreivedUser.Email)
	require.True(t, user.PasswordChangedAt.IsZero())
	require.NotZero(t, user.CreatedAt)
}

func TestUpdateUser(t *testing.T) {
	user, err := CreateRandomUser()
	require.NoError(t, err)

	newEmail := utils.RandomEmail()
	newFullname := "Le Minh Dang"

	updateParam := UpdateUserParams{
		Username: user.Username,
		Email:    newEmail,
		FullName: newFullname,
	}

	// Update User
	err = testQueries.UpdateUser(context.Background(), updateParam)
	require.NoError(t, err)

	// Get User from DB
	userFromDB, err := testQueries.GetUser(context.Background(), user.Username)
	require.NoError(t, err)

	//Check the data consistency
	require.NotEmpty(t, userFromDB)
	require.Equal(t, user.Username, userFromDB.Username)
	require.Equal(t, userFromDB.Email, newEmail)
	require.Equal(t, userFromDB.FullName, newFullname)
}
