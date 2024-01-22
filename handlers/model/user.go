package model

import (
	"time"
)

type User struct {
	Id        string    `json:"username" validate:"required"`
	Email     string    `json:"email" validate:"required"`
	CreatedAt time.Time `json:"createdAt" validate:"required"`
	Name      string    `json:"name"`
	MetaData  MetaData  `json:"metadata"`
}

type MetaData struct {
	IsTest    bool     `json:"isTest"`
	Lifespan  Lifespan `json:"lifespan"`
	ExpiresAt *int64   `json:"expiresAt"`
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
