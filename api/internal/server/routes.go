package server

import (
	"net/http"

	"github.com/endalk200/termflow-api/pkgs/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

func (s *Server) RegisterRoutes() http.Handler {
	validate = validator.New(validator.WithRequiredStructEnabled())

	r := chi.NewRouter()
	r.Use(middleware.Logging(s.logger))

	r.Post("/api/auth/signup", s.CreateUser)
	r.Post("/api/auth/signin", s.SignIn)

	r.Get("/api/auth/me", s.Me)
	// r.With(middleware.Authentication(s.logger)).Get("/api/auth/me", s.Me)

	r.Group(func(r chi.Router) {
		r.Use(middleware.Authentication(s.logger))
		r.Get("/api/auth/me", s.Me)
	})

	return r
}
