package handler

import (
	"encoding/json"
	"net/http"
	"password_generator/internal/dto"
	"password_generator/internal/models"
	"password_generator/internal/utils/token"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt"
	"go.uber.org/zap"
)

type Handler struct {
	service Handlerer
	logger  *zap.Logger
}

type Handlerer interface {
	RegisterNewUser(body dto.User) (*models.User, error)
	GenNewPassword(username string) (*models.User, error)
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

	token, err := token.GenerateJWT(us.ID, us.Username, time.Now().Add(15*time.Minute))
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
}

func (h *Handler) GenNewPassword(w http.ResponseWriter, r *http.Request) {
	// Извлекаем JWT-токен из cookie
	cookie, err := r.Cookie("token")
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Декодируем токен
	tokenString := cookie.Value
	claims := &jwt.MapClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Здесь вы должны вернуть ключ для верификации токена
		return []byte("default"), nil
	})

	if err != nil || !token.Valid {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Получаем username из claims
	username, ok := (*claims)["username"].(string)
	if !ok {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	user, err := h.service.GenNewPassword(username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *Handler) GetAllPasswords(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")

	passwords, err := h.service.GetAllPasswords(username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(passwords)
}

func (h *Handler) DeleteUserPassword(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	password := chi.URLParam(r, "password")

	err := h.service.DeleteUserPassword(username, password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}
