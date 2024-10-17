package router

import (
	"net/http"
	"passwordgenerator/internal/middleware"
	"passwordgenerator/internal/transport/http/handler"

	"github.com/go-chi/chi/v5"
)

type Router interface {
	RegisterNewUser(w http.ResponseWriter, r *http.Request)
	GenNewPassword(w http.ResponseWriter, r *http.Request)
	GetAllPasswords(w http.ResponseWriter, r *http.Request)
	DeleteUserPassword(w http.ResponseWriter, r *http.Request)
}

func New(h *handler.Handler) http.Handler {
	r := chi.NewRouter()

	r.Route("/api/user", func(r chi.Router) {
		r.Post("/registration", h.RegisterNewUser)
		r.With(middleware.JWT).Post("/password/new", h.GenNewPassword)
		r.With(middleware.JWT).Get("/password/{username}", h.GetAllPasswords)
		r.With(middleware.JWT).Delete("/password/{username}/{password}", h.DeleteUserPassword)
	})

	return r
}
