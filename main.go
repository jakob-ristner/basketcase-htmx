package main

import (
	dbHandler "app/internal/database"
	"app/internal/middleware"
	shared "app/internal/template/shared"
	"app/internal/view"
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/joho/godotenv"
)

type User struct {
	ID   int
	Date time.Time
	Name string
}

func main() {

	_ = godotenv.Load()
	mux := http.NewServeMux()

	db := dbHandler.NewDbHandler()
	users := db.GetUsers()

	testListener := make(chan dbHandler.Notification)

	db.AddListener(testListener)
	go func() {
		for notif := range testListener {
			fmt.Println(notif)
		}
	}()

	for _, user := range users {
		fmt.Printf("ID: %d, date: %s, Name: %s\n", user.ID, user.Date.Format(time.DateOnly), user.Name)
	}

	mux.HandleFunc("GET /favicon.ico", view.ServeFavicon)
	mux.HandleFunc("GET /static/", view.ServeStaticFiles)
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		middleware.Chain(w, r, view.Home)
	})

	mux.HandleFunc("GET /recipes", func(w http.ResponseWriter, r *http.Request) {
		middleware.Chain(w, r, view.Recipes)
	})
	mux.HandleFunc("GET /login", func(w http.ResponseWriter, r *http.Request) {
		middleware.Chain(w, r, view.Login)
	})

	mux.HandleFunc("POST /increment", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("inc")
		middleware.Chain(w, r, view.Increment)
	})

	mux.HandleFunc("GET /random", func(w http.ResponseWriter, r *http.Request) {
		middleware.Chain(w, r, func(ctx *middleware.CustomContext, w http.ResponseWriter, r *http.Request) {
			random := rand.Int()
			fmt.Fprintf(w, "Random number: %d", random)
			w.(http.Flusher).Flush()
		})
	})

	mux.HandleFunc("GET /events", func(w http.ResponseWriter, r *http.Request) {
		middleware.Chain(w, r, func(ctx *middleware.CustomContext, w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/event-stream")
			w.Header().Set("Connection", "keep-alive")
			w.Header().Set("HX-Trigger", "event-update")
			eventstream := make(chan string)

			rctx, cancel := context.WithCancel(r.Context())
			var wg sync.WaitGroup

			// Send data to the stream
			wg.Add(1)
			go func() {
				defer wg.Done()
				for {
					buf := new(bytes.Buffer)
					shared.Count(time.Now().Format(time.TimeOnly)).Render(context.Background(), buf)
					select {
					case <-rctx.Done():
						return
					case eventstream <- fmt.Sprintf("event: Event\ndata: %s\n\n", buf.String()):
						time.Sleep(1 * time.Second)
					}
				}
			}()

			// Write data to the client
			wg.Add(1)
			go func() {
				defer cancel()
				defer wg.Done()
				for {
					data := <-eventstream
					if _, err := fmt.Fprintf(w, data); err != nil {
						cancel()
						return
					}
					w.(http.Flusher).Flush() // Ensure the data is sent immediately
				}

			}()

			wg.Wait()
			close(eventstream)
		})
	})
	fmt.Println(fmt.Sprintf("server is running on port %s", os.Getenv("PORT")))
	err := http.ListenAndServe(":"+os.Getenv("PORT"), mux)
	if err != nil {
		fmt.Println(err)
	}
}
