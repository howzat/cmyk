package util

import (
	"fmt"
	"github.com/brianvoe/gofakeit"
	"github.com/projects/cmyk-api/handlers/model"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"time"
)

type TestUserOptions = func(user model.User) model.User

func RandomTestUser(options ...TestUserOptions) model.User {
	firstName := gofakeit.FirstName()
	lastName := gofakeit.LastName()
	title := gofakeit.NamePrefix()
	user := model.User{
		Id:    gofakeit.Username(),
		Email: RandomEmail(firstName, lastName),
		Name:  fmt.Sprintf("%s %s %s", title, firstName, lastName),
		MetaData: model.MetaData{
			IsTest:   true,
			Lifespan: model.Short,
		},
	}

	for _, option := range options {
		user = option(user)
	}

	return user
}

func RandomEmail(firstName string, lastName string) string {
	return fmt.Sprintf("testuser_%s.%s@%s.com", firstName, lastName, gofakeit.WeekDay())
}

func WithCreatedAt(createdAt time.Time) TestUserOptions {
	return func(user model.User) model.User {
		user.CreatedAt = createdAt
		return user
	}
}

func GetOSEnvOrFail(t *testing.T, key string) string {
	value := os.Getenv(key)
	require.NotEmpty(t, value, fmt.Sprintf("environment variable with key [%s] must not be empty", key))
	return value
}
