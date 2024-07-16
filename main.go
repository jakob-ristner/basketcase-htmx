package main

import (
	"app/internal/database"
	"app/internal/handlers"
	"app/internal/middleware"
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {

	_ = godotenv.Load()
	mux := http.NewServeMux()

	database := database.CreateConn()

	mux.HandleFunc("GET /favicon.ico", middleware.Stack(
		handlers.ServeFavicon))

	mux.HandleFunc("GET /static/", middleware.Stack(
		handlers.ServeStaticFiles))

	mux.HandleFunc("GET /", middleware.Stack(
		middleware.AuthenticateSession(database),
		handlers.GetHome,
		middleware.Log,
	))

	mux.HandleFunc("GET /login", middleware.Stack(
		handlers.GetLogin,
		middleware.Log,
	))

	mux.HandleFunc("GET /login2", middleware.Stack(
		handlers.GetLogin2,
		middleware.Log,
	))

	mux.HandleFunc("POST /login", middleware.Stack(
		handlers.PostLogin(database),
		middleware.Log,
	))

	fmt.Println(fmt.Sprintf("server is running on port %s", os.Getenv("PORT")))
	err := http.ListenAndServe(":"+os.Getenv("PORT"), mux)
	if err != nil {
		fmt.Println(err)
	}
}
