package dbHandler

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"golang.org/x/crypto/bcrypt"
)

type DbHandler struct {
	conn *bun.DB
}

var Instance *DbHandler
var singletonLock sync.Once

func GetInstance() *DbHandler {
	singletonLock.Do(func() {
		Instance = new()
	})
	return Instance
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func new() *DbHandler {
	_ = godotenv.Load()
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(os.Getenv("DATABASE_URL"))))
	db := bun.NewDB(sqldb, pgdialect.New())

	handler := &DbHandler{conn: db}

	return handler
}

func (db *DbHandler) GetUserById(id int) (*User, error) {
	user := User{}
	err := db.conn.NewSelect().Model(&user).Where("id = ?", id).Limit(1).Scan(context.Background())
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (db *DbHandler) GetSession(token string) (*UserSession, error) {
	session := UserSession{}

	err := db.conn.NewSelect().Model(&session).Where("token = ?", token).Order("expires desc").Limit(1).Scan(context.Background())
	if err != nil || session.Expires.Before(time.Now()) || session.Token == "" {
		return nil, err
	}

	return &session, nil
}

func (db *DbHandler) AttemptLogin(email string, password string) (*UserSession, error) {

	user := User{}

	err := db.conn.NewSelect().Model(&user).Where("email = ?", email).Limit(1).Scan(context.Background())

	if err != nil {
		return nil, err
	}

	if !checkPasswordHash(password, user.Password) {
		return nil, fmt.Errorf("invalid password")
	}

	session := UserSession{
		UserId:  user.ID,
		Token:   uuid.New().String(),
		Expires: time.Now().Add(time.Hour * 24),
	}
	_, err = db.conn.NewInsert().Model(&session).Exec(context.Background())

	if err != nil {
		return nil, err
	}
	return &session, nil
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
