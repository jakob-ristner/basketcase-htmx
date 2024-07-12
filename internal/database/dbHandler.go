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
)

type Notification struct {
	Channel string
	Payload string
}

type DbHandler struct {
	db        *bun.DB
	ln        *pgdriver.Listener
	listeners []chan Notification
}

var Instance *DbHandler
var singletonLock sync.Once

func GetInstance() *DbHandler {
	singletonLock.Do(func() {
		Instance = new()
	})
	return Instance
}

func new() *DbHandler {
	_ = godotenv.Load()
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(os.Getenv("DATABASE_URL"))))
	db := bun.NewDB(sqldb, pgdialect.New())
	ln := pgdriver.NewListener(db)
	if err := ln.Listen(context.Background(), "users:updated"); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	listeners := make([]chan Notification, 0)

	handler := &DbHandler{db: db, ln: ln, listeners: listeners}

	go func(handler *DbHandler) {
		for notif := range ln.Channel() {
			for _, ch := range handler.listeners {
				ch <- Notification{
					Channel: notif.Channel,
					Payload: notif.Payload,
				}
			}
		}

	}(handler)

	return handler
}

func (db *DbHandler) GetUsers() []User {

	users := make([]User, 0)

	err := db.db.NewRaw(
		"SELECT id, date, name FROM ? LIMIT ?",
		bun.Ident("users"), 100,
	).Scan(context.Background(), &users)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return users
}

func (db *DbHandler) GetUserSessions() []UserSession {

	sessions := make([]UserSession, 0)

	err := db.db.NewRaw(
		"SELECT user_id, token, expires FROM ? LIMIT ?",
		bun.Ident("usersessions"), 100,
	).Scan(context.Background(), &sessions)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return sessions
}

func (db *DbHandler) AuthSession(token string) (*UserSession, error) {
	session := UserSession{}
	err := db.db.NewSelect().Model(&session).Where("token = ?", token).Limit(1).Scan(context.Background())

	if err != nil || session.Expires.Before(time.Now()) || session.Token == "" {
		return nil, err
	}

	return &session, nil
}

func (db *DbHandler) AttemptLogin(email string, password string) (*UserSession, error) {
	//TODO move to authhandler

	user := userWithPassword{}

	err := db.db.NewSelect().Model(&user).Where("email = ?", email).Limit(1).Scan(context.Background())

	if err != nil {
		return nil, err
	}

	if user.Password != password {
		return nil, fmt.Errorf("invalid password")
	}

	session := UserSession{
		UserId:  user.ID,
		Token:   uuid.New().String(),
		Expires: time.Now().Add(time.Hour * 24),
	}
	fmt.Println(session)
	_, err = db.db.NewInsert().Model(&session).Exec(context.Background())

	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (db *DbHandler) AddListener(ch chan Notification) {
	db.listeners = append(db.listeners, ch)
}

type User struct {
	ID    int
	Date  time.Time
	Name  string
	Email string
}

type userWithPassword struct {
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
