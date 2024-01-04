package util

import (
	"fmt"
	"github.com/brianvoe/gofakeit"
	"github.com/projects/cmyk-tools/handlers/model"
)

func RandomUser() model.User {
	return model.User{
		Username: gofakeit.Username(),
		Email:    randomEmail(),
	}
}

func randomEmail() string {
	return fmt.Sprintf("%s.%s@%s.com", gofakeit.FirstName(), gofakeit.LastName(), gofakeit.WeekDay())
}
