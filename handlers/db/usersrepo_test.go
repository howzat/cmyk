package db

import (
	"context"
	"github.com/brianvoe/gofakeit"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestStoreAndRetrieveUser(t *testing.T) {

	ctx := context.TODO()
	err := godotenv.Load("../.env.local")
	if err != nil {
		t.Fatal("Error loading .env file")
	}

	ok := os.Getenv("LOADED")
	assert.Equal(t, "OK", ok)
	repo := NewUsersTableRepo(ctx, os.Getenv("AWS_REGION"))
	repo.AddUser(ctx, gofakeit.Username(), gofakeit.Email())
}
