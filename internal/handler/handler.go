package handler

import (
	dbHandler "app/internal/database"
	"app/internal/middleware"
	"fmt"
	"net/http"
)

func Login(ctx *middleware.CustomContext, w http.ResponseWriter, r *http.Request) error {
	r.ParseForm()
	email := r.Form.Get("email")
	password := r.Form.Get("password")
	session, err := dbHandler.GetInstance().AttemptLogin(email, password)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return err
	} else {
		w.Header().Set("Set-Cookie", fmt.Sprintf("session=%s; HttpOnly; SameSite=Lax", session.Token))
		w.Header().Set("HX-Redirect", "/")
	}
	return nil
}
