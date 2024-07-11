package view

import (
	"app/internal/middleware"
	"app/internal/template"
	shared "app/internal/template/shared"
	"net/http"
	"path/filepath"
	"strconv"
)

func ServeFavicon(w http.ResponseWriter, r *http.Request) {
	filePath := "favicon.ico"
	fullPath := filepath.Join(".", "static", filePath)
	http.ServeFile(w, r, fullPath)
}

func ServeStaticFiles(w http.ResponseWriter, r *http.Request) {
	filePath := r.URL.Path[len("/static/"):]
	fullPath := filepath.Join(".", "static", filePath)
	http.ServeFile(w, r, fullPath)
}

var count int = 0

func Home(ctx *middleware.CustomContext, w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return

	}
	w.Header().Set("Cache-Control", "no-cache")
	template.Home("Templ Quickstart", strconv.Itoa(count)).Render(ctx, w)
}

func Recipes(ctx *middleware.CustomContext, w http.ResponseWriter, r *http.Request) {
	template.Recipes().Render(ctx, w)
}

func Increment(ctx *middleware.CustomContext, w http.ResponseWriter, r *http.Request) {
	count++
	shared.Count(strconv.Itoa(count)).Render(ctx, w)
}
