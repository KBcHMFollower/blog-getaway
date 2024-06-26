package handlers

import (
	"github.com/unrolled/render"
	"net/http"
	"strings"
	"test-plate/internal/domain/models"
	"test-plate/internal/services"
)

type AuthHandler struct {
	authService *services.AuthService
	rnd         *render.Render
}

func NewAuthHandler(userService *services.AuthService, rnd *render.Render) *AuthHandler {
	return &AuthHandler{userService, rnd}
}

func (uh *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		uh.rnd.Text(w, http.StatusBadRequest, "Unable to parse form")
		return
	}

	email := r.Form.Get("email")
	if len(email) == 0 {
		uh.rnd.Text(w, http.StatusBadRequest, "Email is required")
		return
	}
	password := r.Form.Get("password")
	if len(password) == 0 {
		uh.rnd.Text(w, http.StatusBadRequest, "Password is required")
		return
	}
	FName := r.Form.Get("fname")
	if len(FName) == 0 {
		uh.rnd.Text(w, http.StatusBadRequest, "Fname is required")
		return
	}
	LName := r.Form.Get("lname")
	if len(LName) == 0 {
		uh.rnd.Text(w, http.StatusBadRequest, "Lname is required")
		return
	}

	token, err := uh.authService.Register(r.Context(), &models.RegisterData{
		email,
		password,
		FName,
		LName,
	})
	if err != nil {
		uh.rnd.Text(w, http.StatusInternalServerError, "Unable to register user")
		return
	}

	uh.rnd.JSON(w, http.StatusOK, token)
}
func (uh *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		uh.rnd.Text(w, http.StatusBadRequest, "Unable to parse form")
		return
	}

	email := r.Form.Get("email")
	if len(email) == 0 {
		uh.rnd.Text(w, http.StatusBadRequest, "Email is required")
	}
	password := r.Form.Get("password")
	if len(password) == 0 {
		uh.rnd.Text(w, http.StatusBadRequest, "Password is required")
	}

	token, err := uh.authService.Login(r.Context(), &models.LoginData{
		email,
		password,
	})
	if err != nil {
		uh.rnd.Text(w, http.StatusInternalServerError, "Unable to register user")
	}

	uh.rnd.JSON(w, http.StatusOK, token)
}
func (uh *AuthHandler) CheckAuth(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if len(token) == 0 {
		uh.rnd.JSON(w, http.StatusUnauthorized, map[string]string{"error": "Missing token"})
	}

	token = strings.TrimPrefix(token, "Bearer ")

	tokenRes, err := uh.authService.CheckAuth(r.Context(), &models.TokenData{
		Token: token,
	})
	if err != nil {
		uh.rnd.Text(w, http.StatusUnauthorized, "Unable to check token")
	}

	uh.rnd.JSON(w, http.StatusOK, tokenRes)
}
