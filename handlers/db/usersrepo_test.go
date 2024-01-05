package db

import (
	"context"
	"github.com/joho/godotenv"
	"github.com/projects/cmyk-tools/handlers/util"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestStoreAndRetrieveUser(t *testing.T) {

	ctx := context.TODO()

	err := godotenv.Load("../../.env")
	if err != nil {
		t.Fatal("Error loading .env file")
	}

	region := os.Getenv("AWS_REGION")

	t.Log(region)

	repo := NewUsersTableRepo(ctx, region)
	u := util.RandomUser()
	savedUser, err := repo.AddUser(ctx, u.Username, u.Email)

	assert.NoError(t, err, "nope")

	got, err := repo.GetUser(ctx, u.Email)
	assert.NoError(t, err)

	assert.Equal(t, savedUser.Id, got.Id)
	assert.Equal(t, savedUser.Email, got.Email)
	assert.Equal(t, savedUser.Username, got.Username)
	assert.Equal(t, savedUser.CreatedAt, got.CreatedAt)
}
