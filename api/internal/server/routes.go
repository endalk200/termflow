package server

import (
	"encoding/json"
	"net/http"

	"github.com/endalk200/termflow-api/internal/database"
	"github.com/endalk200/termflow-api/internal/helpers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"

	"github.com/jackc/pgx/v5/pgtype"
)

type User struct {
	FirstName    string `validate:"required"`
	LastName     string `validate:"required"`
	Email        string `validate:"required,email"`
	GitHubHandle string `validate:"required"`
}

// type Address struct {
// 	Street string `validate:"required"`
// 	City   string `validate:"required"`
// 	Planet string `validate:"required"`
// 	Phone  string `validate:"required"`
// }

var validate *validator.Validate

func (s *Server) RegisterRoutes() http.Handler {
	validate = validator.New()

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/users", s.GetAuthorsHandler)
	r.Post("/users", s.CreateAuthorHandler)

	return r
}

func (s *Server) GetAuthorsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	authors, err := s.db.ListUsers(ctx)
	if err != nil {
		http.Error(w, "Failed to fetch authors", http.StatusInternalServerError)
		return
	}

	helpers.Response(w, http.StatusOK, authors)
}

func (s *Server) CreateAuthorHandler(w http.ResponseWriter, r *http.Request) {
	var user User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		helpers.ResponseError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := validate.Struct(user); err != nil {
		// Extract validation errors
		validationErrors := err.(validator.ValidationErrors)
		helpers.ResponseError(w, http.StatusBadRequest, validationErrors.Error())
		return
	}

	ctx := r.Context()
	newUser, err := s.db.CreateUser(ctx, database.CreateUserParams{
		FirstName:       user.FirstName,
		LastName:        user.LastName,
		Email:           user.Email,
		IsEmailVerified: pgtype.Bool{Bool: false, Valid: true},
		IsActive:        pgtype.Bool{Bool: false, Valid: true},
		Password:        "password",
	})
	if err != nil {
		helpers.ResponseError(w, http.StatusInternalServerError, "Failed to create user")
		return
	}

	helpers.Response(w, http.StatusCreated, newUser)
}
