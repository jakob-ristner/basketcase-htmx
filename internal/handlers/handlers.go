package handlers

import (
	"app/internal/database"
	"app/internal/middleware"
	"app/internal/template/admin"
	"app/internal/template/ingredients"
	"app/internal/template/lists"
	"app/internal/template/login2"
	"app/internal/template/recipes"
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

func GetLists(w http.ResponseWriter, r *http.Request) {
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
	lists.Lists(user).Render(r.Context(), w)
}

func GetRecipes(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(middleware.UserKeyContext).(*database.User)

	if !ok {
		ctx, cancel := context.WithCancel(r.Context())
		cancel()
		*r = *r.WithContext(ctx)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	recipes.Recipes(user).Render(r.Context(), w)
}

func GetIngredients(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(middleware.UserKeyContext).(*database.User)

	if !ok {
		ctx, cancel := context.WithCancel(r.Context())
		cancel()
		*r = *r.WithContext(ctx)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	ingredients.Ingredients(user).Render(r.Context(), w)
}

func GetAdmin(w http.ResponseWriter, r *http.Request) {

	user, ok := r.Context().Value(middleware.UserKeyContext).(*database.User)

	if !ok {
		ctx, cancel := context.WithCancel(r.Context())
		cancel()
		*r = *r.WithContext(ctx)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	admin.Admin(user).Render(r.Context(), w)
}

func GetNav(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(middleware.UserKeyContext).(*database.User)
	r.ParseForm()
	route := r.Form.Get("route")

	if !ok {
		ctx, cancel := context.WithCancel(r.Context())
		cancel()
		*r = *r.WithContext(ctx)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if route == "lists" {
		w.Header().Set("HX-Push-Url", "/")
		lists.Lists(user).Render(r.Context(), w)
		return
	} else if route == "recipes" {
		w.Header().Set("HX-Push-Url", "/recipes")
		recipes.Recipes(user).Render(r.Context(), w)
		return
	} else if route == "ingredients" {
		w.Header().Set("HX-Push-Url", "/ingredients")
		ingredients.Ingredients(user).Render(r.Context(), w)
		return
	} else if route == "admin" {

		if !user.Admin {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		w.Header().Set("HX-Push-Url", "/admin")
		admin.Admin(user).Render(r.Context(), w)
		return
	}
	w.WriteHeader(http.StatusBadRequest)
}
