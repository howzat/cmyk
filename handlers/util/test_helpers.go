package util

import (
	"fmt"
	"github.com/brianvoe/gofakeit"
	"github.com/projects/cmyk-tools/handlers/model"
)

func RandomUser() model.User {
	timestamp, ulid, err := CurrentTimeAndULID(NewRealClock())
	if err != nil {
		panic(err)
	}

	return model.User{
		Username:  gofakeit.Username(),
		Email:     randomEmail(),
		Id:        ulid.String(),
		CreatedAt: timestamp,
	}
}

func randomEmail() string {
	return fmt.Sprintf("%s.%s@%s.com", gofakeit.FirstName(), gofakeit.LastName(), gofakeit.WeekDay())
}
