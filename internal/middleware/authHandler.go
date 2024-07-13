package middleware

import (
	dbHandler "app/internal/database"
	"fmt"
	"net/http"
)

func Auth(ctx *CustomContext, w http.ResponseWriter, r *http.Request) error {
	token, err := r.Cookie("session")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return err
	}
	session, err := dbHandler.GetInstance().AuthSession(token.Value)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return err
	}

	ctx.token = session.Token
	return nil
}

func Login(ctx *CustomContext, w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	email := r.Form.Get("email")
	password := r.Form.Get("password")
	session, err := dbHandler.GetInstance().AttemptLogin(email, password)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	} else {
		w.Header().Set("Set-Cookie", fmt.Sprintf("session=%s; HttpOnly; SameSite=Lax", session.Token))
		w.Header().Set("HX-Redirect", "/")
		w.WriteHeader(http.StatusOK)
	}
}
