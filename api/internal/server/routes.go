package server

import (
	"net/http"

	"github.com/endalk200/termflow-api/pkgs/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

func (s *Server) RegisterRoutes() http.Handler {
	Validate = validator.New(validator.WithRequiredStructEnabled())

	r := chi.NewRouter()
	r.Use(middleware.Logging(s.logger))

	r.Route("/api", func(r chi.Router) {
		r.Post("/auth/signup", s.CreateUser)
		r.Post("/auth/signin", s.SignIn)

		r.Group(func(r chi.Router) {
			r.Use(middleware.Authentication(s.logger))
			r.Get("/auth/me", s.Me)

			r.Route("/tags", func(r chi.Router) {
				r.Get("/", s.GetTags)
				r.Post("/", s.CreateTag)
				r.Put("/{id}", s.UpdateTag)
				r.Delete("/{id}", s.DeleteTag)

				r.Get("/{id}/commands", s.GetCommandsWithTag)
			})

			r.Route("/commands", func(r chi.Router) {
				r.Get("/", s.GetCommands)
				r.Post("/", s.CreateCommand)
				r.Put("/{id}", s.UpdateCommand)
				r.Delete("/{id}", s.DeleteCommand)
			})
		})
	})

	return r
}
