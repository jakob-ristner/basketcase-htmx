package handlers

import (
	"app/internal/database"
	"app/internal/middleware"
	template "app/internal/template/home"
	"app/internal/template/login2"
	"context"
	"net/http"
	"path/filepath"
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

func GetLogin2(w http.ResponseWriter, r *http.Request) {
	login2.Login().Render(r.Context(), w)
}

func GetHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return

	}
	user, ok := r.Context().Value(middleware.UserKeyContext).(*database.User)

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
