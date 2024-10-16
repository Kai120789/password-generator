package handler

import (
	"encoding/json"
	"net/http"
	"password_generator/internal/dto"
	"password_generator/internal/models"
	"password_generator/internal/utils/token"
	"time"

	"go.uber.org/zap"
)

type Handler struct {
	service Handlerer
	logger  *zap.Logger
}

type Handlerer interface {
	RegisterNewUser(body dto.User) (*models.User, error)
	GenNewPassword(body dto.User) (*models.User, error)
	GetAllPasswords(username string) (*[]models.User, error)
	DeleteUserPassword(username string, password string) error
}

func New(s Handlerer, l *zap.Logger) Handler {
	return Handler{
		service: s,
		logger:  l,
	}
}

func (h *Handler) RegisterNewUser(w http.ResponseWriter, r *http.Request) {
	var user dto.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if user.Username == "" {
		http.Error(w, "username cannot be empty", http.StatusBadRequest)
		return
	}

	us, err := h.service.RegisterNewUser(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	token, err := token.GenerateJWT(us.ID, time.Now().Add(15*time.Minute))
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	TokenCokie := http.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		Expires:  time.Now().Add(15 * time.Minute),
		HttpOnly: true,
		Secure:   false,
	}

	http.SetCookie(w, &TokenCokie)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(token)
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) GenNewPassword(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) GetAllPasswords(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) DeleteUserPassword(w http.ResponseWriter, r *http.Request) {

}
