package database

import (
	"time"

	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel `bun:"users"`
	ID            int
	Date          time.Time
	Name          string
	Email         string
	Password      string
	Admin         bool
}

type UserSession struct {
	bun.BaseModel `bun:"usersessions"`
	UserId        int
	Token         string
	Expires       time.Time
}
