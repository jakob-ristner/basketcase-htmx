package dbHandler

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

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

func NewDbHandler() *DbHandler {
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

func (db *DbHandler) AddListener(ch chan Notification) {
	db.listeners = append(db.listeners, ch)
}

type User struct {
	ID   int
	Date time.Time
	Name string
}

// var _ = godotenv.Load()
// var sqldb = sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(os.Getenv("DATABASE_URL"))))
// var db = bun.NewDB(sqldb, pgdialect.New())

// func GetUsers() []User {

// 	users := make([]User, 0)

// 	err := db.NewRaw(
// 		"SELECT id, date, name FROM ? LIMIT ?",
// 		bun.Ident("users"), 100,
// 	).Scan(context.Background(), &users)

// 	if err != nil {
// 		fmt.Println(err)
// 		os.Exit(1)
// 	}

// 	for _, user := range users {
// 		fmt.Println(user)
// 	}

// 	return users

// }
