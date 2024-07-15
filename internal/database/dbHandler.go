package database

import (
	"context"
	"database/sql"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

type bunWrapper struct {
	conn *bun.DB
}

func CreateConn() Connection {
	return new()
}

type Connection interface {
	GetUserById(id int) (*User, error)
	GetUserByEmail(email string) (*User, error)
	GetSessionByToken(token string) (*UserSession, error)

	InsertUserSession(session *UserSession) error
}

func new() *bunWrapper {
	_ = godotenv.Load()
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(os.Getenv("DATABASE_URL"))))
	db := bun.NewDB(sqldb, pgdialect.New())

	handler := &bunWrapper{conn: db}

	return handler
}

func (db *bunWrapper) GetUserById(id int) (*User, error) {
	user := User{}
	err := db.conn.NewSelect().Model(&user).Where("id = ?", id).Limit(1).Scan(context.Background())
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (db *bunWrapper) GetUserByEmail(email string) (*User, error) {
	user := User{}
	err := db.conn.NewSelect().Model(&user).Where("email = ?", email).Limit(1).Scan(context.Background())
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (db *bunWrapper) GetSessionByToken(token string) (*UserSession, error) {
	session := UserSession{}

	err := db.conn.NewSelect().Model(&session).Where("token = ?", token).Order("expires desc").Limit(1).Scan(context.Background())
	if err != nil || session.Expires.Before(time.Now()) {
		return nil, err
	}

	return &session, nil
}

func (db *bunWrapper) InsertUserSession(session *UserSession) error {
	_, err := db.conn.NewInsert().Model(session).Exec(context.Background())
	return err
}

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
