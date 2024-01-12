package model

import (
	"time"
)

type User struct {
	Id        string    `json:"username" validate:"required"`
	Email     string    `json:"email" validate:"required"`
	CreatedAt time.Time `json:"createdAt" validate:"required"`
	Name      string    `json:"name"`
	Ttl       *int64    `json:"ttl"`
}
