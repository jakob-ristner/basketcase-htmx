package handlers

import (
	"app/internal/database"
	"app/internal/template/login"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func GetLogin(w http.ResponseWriter, r *http.Request) {
	login.Login().Render(r.Context(), w)
}

func PostLogin(database database.Connection) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		email := r.Form.Get("email")
		password := r.Form.Get("password")

		session, err := tryLogin(email, password, database)

		if err != nil {
			w.Header().Set("Status-Code", strconv.Itoa(http.StatusUnauthorized))
			w.Header().Set("HX-Retarget", "input[name='password']")
			login.PasswordError().Render(r.Context(), w)
		} else {
			w.Header().Set("Set-Cookie", fmt.Sprintf("token=%s; HttpOnly; SameSite=Lax", session.Token))
			w.Header().Set("HX-Redirect", "/")

			w.WriteHeader(http.StatusOK)
		}
	}
}

const dummyhash string = "$2a$14$6wxoUBumzfhGb6EMG3/w7ejIh9hJDbeOedwLtYm0sqt7Wl3JCdM4q"

func checkPasswordHash(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func spoofPasswordCheck(password string) {
	checkPasswordHash(password, dummyhash)
}
func tryLogin(email string, password string, conn database.Connection) (*database.UserSession, error) {

	user, err := conn.GetUserByEmail(email)

	if err != nil {
		spoofPasswordCheck(password)
		return nil, err
	}

	if !checkPasswordHash(password, user.Password) {
		return nil, fmt.Errorf("invalid password")
	}

	session := database.UserSession{
		UserId:  user.ID,
		Token:   uuid.New().String(),
		Expires: time.Now().Add(time.Hour * 24),
	}

	err = conn.InsertUserSession(&session)

	if err != nil {
		return nil, err
	}
	return &session, nil
}
