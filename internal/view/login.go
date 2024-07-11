package view

import (
	"app/internal/middleware"
	"app/internal/template/login"
	"net/http"
)

func Login(ctx *middleware.CustomContext, w http.ResponseWriter, r *http.Request) {
	login.Login().Render(ctx, w)
}
