package middleware

import (
	"app/internal/database"
	"context"
	"fmt"
	"net/http"
	"time"
)

type key int

//const timeKeyContext key = iota

const (
	timeKeyContext key = iota
	UserKeyContext
)

func Stack(handlers ...http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), timeKeyContext, time.Now())
		r = r.WithContext(ctx)
		for _, handler := range handlers {
			handler.ServeHTTP(w, r)
			if r.Context().Err() != nil {
				return
			}
		}
	}
}

func AuthenticateSession(conn database.Connection) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		token, err := r.Cookie("token")

		ctx, cancel := context.WithCancel(r.Context())
		cancel()

		if err != nil {
			*r = *r.WithContext(ctx)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		session, err := conn.GetSessionByToken(token.Value)
		if err != nil || session.Expires.Before(time.Now()) {
			*r = *r.WithContext(ctx)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		user, err := conn.GetUserById(session.UserId)
		if err != nil {
			*r = *r.WithContext(ctx)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		ctx = context.WithValue(r.Context(), UserKeyContext, user)
		*r = *r.WithContext(ctx)
	}
}

func Log(w http.ResponseWriter, r *http.Request) {
	startTime, ok := r.Context().Value(timeKeyContext).(time.Time)
	if !ok {
		return
	}
	elapsedTime := time.Since(startTime)
	formattedTime := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("[%s] [%s] [%s] [%s]\n", formattedTime, r.Method, r.URL.Path, elapsedTime)
}

func EnsureAdmin(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(UserKeyContext).(*database.User)
	ctx, cancel := context.WithCancel(r.Context())
	cancel()
	if !ok {
		*r = *r.WithContext(ctx)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !user.Admin {
		*r = *r.WithContext(ctx)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
}
