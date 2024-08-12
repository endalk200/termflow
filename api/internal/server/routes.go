package server

import (
	"encoding/json"
	"net/http"

	"github.com/endalk200/termflow-api/internal/database"
	"github.com/endalk200/termflow-api/internal/helpers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/jackc/pgx/v5/pgtype"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/authors", s.GetAuthorsHandler)
	r.Post("/authors", s.CreateAuthorHandler)

	return r
}

func (s *Server) GetAuthorsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	authors, err := s.db.ListAuthors(ctx)
	if err != nil {
		http.Error(w, "Failed to fetch authors", http.StatusInternalServerError)
		return
	}

	helpers.Response(w, http.StatusOK, authors)
}

func (s *Server) CreateAuthorHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name string `json:"name"`
		Bio  string `json:"bio"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helpers.ResponseError(w, http.StatusBadRequest, "Something went wrong")
		return
	}

	if req.Bio == "" {
		helpers.ResponseError(w, http.StatusBadRequest, "Bio is required field")
		return
	}

	ctx := r.Context()
	author, err := s.db.CreateAuthor(ctx, database.CreateAuthorParams{
		Name: req.Name,
		Bio:  pgtype.Text{String: req.Bio, Valid: true},
	})
	if err != nil {
		http.Error(w, "Failed to create author", http.StatusInternalServerError)
		return
	}

	helpers.Response(w, http.StatusCreated, author)
}
