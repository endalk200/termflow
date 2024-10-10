package server

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/endalk200/termflow-api/models"
	"github.com/endalk200/termflow-api/pkgs/auth"
	"github.com/endalk200/termflow-api/pkgs/middleware"
	"github.com/endalk200/termflow-api/pkgs/utils"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

func (s *Server) decodeAndValidate(w http.ResponseWriter, r *http.Request, payload interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(payload); err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "Invalid request payload")
		return err
	}

	if err := validate.Struct(payload); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		s.logger.Error("Error during request payload validation", slog.String("ERROR", validationErrors.Error()))
		utils.ResponseError(w, http.StatusBadRequest, validationErrors.Error())
		return err
	}

	return nil
}

var validate *validator.Validate

type signUpRequestPayloadSchema struct {
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required"`
}

func (s *Server) CreateUser(w http.ResponseWriter, r *http.Request) {
	var requestPayload signUpRequestPayloadSchema
	if err := s.decodeAndValidate(w, r, &requestPayload); err != nil {
		return
	}

	hash, err := auth.HashPassword(requestPayload.Password, auth.Bcrypt)
	if err != nil {
		s.logger.Error("Error while hashing user password", slog.String("ERROR", err.Error()))
		utils.ResponseError(w, http.StatusInternalServerError, "Something went wrong while creating a user account")
		return
	}

	ctx := r.Context()
	user, err := s.db.CreateUser(ctx, models.CreateUserParams{
		FirstName:       requestPayload.FirstName,
		LastName:        requestPayload.LastName,
		Email:           requestPayload.Email,
		IsEmailVerified: pgtype.Bool{Bool: false, Valid: true},
		IsActive:        pgtype.Bool{Bool: false, Valid: true},
		Password:        string(hash),
	})
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == pgerrcode.UniqueViolation {
				if pgErr.ConstraintName == "users_email_key" {
					s.logger.Error(pgErr.Message)
					utils.ResponseError(w, http.StatusConflict, "Email already exists")
					return
				}
				if pgErr.ConstraintName == "users_github_handle_key" {
					s.logger.Error(pgErr.Message)
					utils.ResponseError(w, http.StatusConflict, "GitHub handle already exists")
					return
				}
			}
		}

		s.logger.Error(err.Error())
		utils.ResponseError(w, http.StatusInternalServerError, "Failed to create user")
		return
	}

	responsePayload := map[string]interface{}{
		"id":                user.ID,
		"first_name":        user.FirstName,
		"last_name":         user.LastName,
		"email":             user.Email,
		"is_email_verified": user.IsEmailVerified,
		"isActive":          user.IsActive,
	}

	utils.Response(w, http.StatusCreated, responsePayload)
}

type signInRequestPayloadSchema struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func (s *Server) SignIn(w http.ResponseWriter, r *http.Request) {
	var requestPayload signInRequestPayloadSchema

	if err := json.NewDecoder(r.Body).Decode(&requestPayload); err != nil {
		s.logger.Error("Error during request payload decoding", slog.String("ERROR", err.Error()))
		utils.ResponseError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := validate.Struct(requestPayload); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		s.logger.Error("Error during request payload validation", slog.String("ERROR", validationErrors.Error()))
		utils.ResponseError(w, http.StatusBadRequest, validationErrors.Error())
		return
	}

	ctx := r.Context()
	user, err := s.db.GetUser(ctx, models.GetUserParams{
		Email: requestPayload.Email,
	})
	if err != nil {
		if err == pgx.ErrNoRows {
			s.logger.Error("No user with specified email", slog.String("ERROR", err.Error()))
			utils.ResponseError(w, http.StatusUnauthorized, "Invalid credentials")
		} else {
			s.logger.Error("Failed to fetch user due to", slog.String("ERROR", err.Error()))
			utils.ResponseError(w, http.StatusInternalServerError, "Invalid credentials")
		}

		return
	}

	isMatch, err := auth.CompareHash(requestPayload.Password, user.Password, auth.Bcrypt)
	if !isMatch || err != nil {
		if err != nil {
			s.logger.Error("Error during password hash comparison", slog.String("ERROR", err.Error()))
			utils.ResponseError(w, http.StatusInternalServerError, "Something went wrong while trying to log you in")
		} else {
			s.logger.Error("Invalid password attempt", slog.String("email", requestPayload.Email))
			utils.ResponseError(w, http.StatusUnauthorized, "Invalid credentials")
		}
		return
	}

	jwtClaims := jwt.RegisteredClaims{
		Issuer:    "Termflow",
		Subject:   fmt.Sprintf("%d", user.ID),
		Audience:  jwt.ClaimStrings{"https://example.com"},
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	token, err := auth.GenerateJWT(jwtClaims)
	if err != nil {
		s.logger.Warn("Error while generating jwt", slog.String("ERROR", err.Error()))
		utils.ResponseError(w, http.StatusInternalServerError, "Something went wrong while trying to log you in")
		return
	}

	refreshTokenClaims := jwt.RegisteredClaims{
		Issuer:    "Termflow",
		Subject:   fmt.Sprintf("%d", user.ID),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	refreshToken, err := auth.GenerateJWT(refreshTokenClaims)
	if err != nil {
		s.logger.Warn("Error while generating refresh token", slog.String("ERROR", err.Error()))
		utils.ResponseError(w, http.StatusInternalServerError, "Something went wrong while trying to log you in")
		return
	}

	hashedRefreshToken, err := auth.HashPassword(refreshToken, auth.SHA256)
	if err != nil {
		s.logger.Error("Error during refreshToken hash", slog.String("ERROR", err.Error()))
		utils.ResponseError(w, http.StatusInternalServerError, "Something went wrong while trying to singin")
		return
	}

	err = s.db.CreateRefreshToken(ctx, models.CreateRefreshTokenParams{
		UserID:    pgtype.Int8{Int64: int64(user.ID), Valid: true},
		TokenHash: hashedRefreshToken,
		ExpiresAt: pgtype.Timestamptz{Time: time.Now().Add(time.Hour), Valid: true},
	})
	if err != nil {
		s.logger.Error("Error during recording of the refresh token", slog.String("ERROR", err.Error()))
		utils.ResponseError(w, http.StatusInternalServerError, "Something went wrong while trying to singin")
		return
	}

	responsePayload := map[string]interface{}{
		"id":                user.ID,
		"first_name":        user.FirstName,
		"last_name":         user.LastName,
		"email":             user.Email,
		"is_email_verified": user.IsEmailVerified,
		"isActive":          user.IsActive,
		"token":             token,
		"refreshToken":      refreshToken,
	}

	utils.Response(w, http.StatusOK, responsePayload)
}

func (s *Server) Me(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserFromContext(r)
	if !ok {
		s.logger.Error("userId not found in request context")
		utils.ResponseError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	_userId, _ := strconv.Atoi(userID)

	ctx := r.Context()
	user, err := s.db.GetUser(ctx, models.GetUserParams{
		ID: int32(_userId),
	})
	if err != nil {
		if err == pgx.ErrNoRows {
			s.logger.Error("No user with specified email", slog.String("ERROR", err.Error()))
			utils.ResponseError(w, http.StatusUnauthorized, "Invalid credentials")
		} else {
			s.logger.Error("Failed to fetch user due to", slog.String("ERROR", err.Error()))
			utils.ResponseError(w, http.StatusInternalServerError, "Invalid credentials")
		}

		return
	}

	responsePayload := map[string]interface{}{
		"id":                user.ID,
		"first_name":        user.FirstName,
		"last_name":         user.LastName,
		"email":             user.Email,
		"is_email_verified": user.IsEmailVerified,
		"isActive":          user.IsActive,
	}

	utils.Response(w, http.StatusOK, responsePayload)
}
