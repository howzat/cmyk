package util

import (
	"fmt"
	"github.com/brianvoe/gofakeit"
	"github.com/projects/cmyk-tools/handlers/model"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"time"
)

type TestUserOptions = func(user model.User) model.User

func RandomUser(options ...TestUserOptions) model.User {
	firstName := gofakeit.FirstName()
	lastName := gofakeit.LastName()
	title := gofakeit.NamePrefix()
	user := model.User{
		Id:    gofakeit.Username(),
		Email: fmt.Sprintf("%s.%s@%s.com", firstName, lastName, gofakeit.WeekDay()),
		Name:  fmt.Sprintf("%s %s %s", title, firstName, lastName),
	}

	for _, option := range options {
		user = option(user)
	}

	return user
}

func WithCreatedAt(createdAt time.Time) TestUserOptions {
	return func(user model.User) model.User {
		user.CreatedAt = createdAt
		return user
	}
}

type Lifespan int64

const (
	None Lifespan = iota
	Short
)

func TestLifespan(l Lifespan, now time.Time) int64 {
	switch l {
	case Short:
		return now.Add(1 * (24 * time.Hour)).Unix()
	default:
		return TestLifespan(Short, now)
	}
}

func GetOSEnvOrFail(t *testing.T, key string) string {
	value := os.Getenv(key)
	require.NotEmpty(t, value, fmt.Sprintf("environment variable with key [%s] must not be empty", key))
	return value
}
