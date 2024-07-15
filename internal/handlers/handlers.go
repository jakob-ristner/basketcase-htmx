package handlers

import (
	dbHandler "app/internal/database"
	"app/internal/middleware"
	template "app/internal/template/home"
	"app/internal/template/login"
	"context"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
)

func ServeFavicon(w http.ResponseWriter, r *http.Request) {
	filePath := "svg/logo.svg"
	fullPath := filepath.Join(".", "static", filePath)
	http.ServeFile(w, r, fullPath)
}

func ServeStaticFiles(w http.ResponseWriter, r *http.Request) {
	filePath := r.URL.Path[len("/static/"):]
	fullPath := filepath.Join(".", "static", filePath)
	http.ServeFile(w, r, fullPath)
}

func GetHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return

	}
	user, ok := r.Context().Value(middleware.UserKeyContext).(*dbHandler.User)

	if !ok {
		ctx, cancel := context.WithCancel(r.Context())
		cancel()
		*r = *r.WithContext(ctx)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Cache-Control", "no-cache")
	template.Home(user).Render(r.Context(), w)
}

func GetLogin(w http.ResponseWriter, r *http.Request) {
	login.Login().Render(r.Context(), w)
}

func PostLogin(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	email := r.Form.Get("email")
	password := r.Form.Get("password")
	session, err := dbHandler.GetInstance().AttemptLogin(email, password)
	if err != nil {
		w.Header().Set("Status-Code", strconv.Itoa(http.StatusUnauthorized))
		w.Header().Set("HX-Retarget", "#password-container")
		login.IncorrectLoginPwContainer().Render(r.Context(), w)
	} else {
		w.Header().Set("Set-Cookie", fmt.Sprintf("token=%s; HttpOnly; SameSite=Lax", session.Token))
		w.Header().Set("HX-Redirect", "/")
		w.WriteHeader(http.StatusOK)
	}
}
