package model

import (
	"time"
)

type User struct {
	Username  string    `json:"username" validate:"required"`
	Email     string    `json:"email" validate:"required"`
	Id        string    `json:"id" validate:"required"`
	CreatedAt time.Time `json:"createdAt" validate:"required"`
}
