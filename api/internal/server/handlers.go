package server

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/endalk200/termflow-api/internal/repository"
	"github.com/endalk200/termflow-api/pkgs/auth"
	"github.com/endalk200/termflow-api/pkgs/middleware"
	"github.com/endalk200/termflow-api/pkgs/utils"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

func (s *Server) DecodeAndValidate(w http.ResponseWriter, r *http.Request, payload interface{}) error {
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
	if err := s.DecodeAndValidate(w, r, &requestPayload); err != nil {
		return
	}

	hash, err := auth.HashPassword(requestPayload.Password, auth.Bcrypt)
	if err != nil {
		s.logger.Error("Error while hashing user password", slog.String("ERROR", err.Error()))
		utils.ResponseError(w, http.StatusInternalServerError, "Something went wrong while creating a user account")
		return
	}

	ctx := r.Context()
	user, err := s.db.InsertUser(ctx, repository.InsertUserParams{
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
	user, err := s.db.GetUser(ctx, repository.GetUserParams{
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
		Subject:   user.ID.String(),
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

	err = s.db.CreateRefreshToken(ctx, repository.CreateRefreshTokenParams{
		UserID:    pgtype.UUID{Bytes: user.ID, Valid: true},
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

	// TODO: Handle uuid parse error
	_userId, _ := uuid.Parse(userID)

	ctx := r.Context()
	user, err := s.db.GetUser(ctx, repository.GetUserParams{
		ID: _userId,
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
		"created_at":        user.CreatedAt,
	}

	utils.Response(w, http.StatusOK, responsePayload)
}

type createTagRequestPayloadSchema struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:""`
}

func (s *Server) CreateTag(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserFromContext(r)
	if !ok {
		s.logger.Error("userId not found in request context")
		utils.ResponseError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// TODO: Handle uuid parse error
	_userId, _ := uuid.Parse(userID)

	var requestPayload createTagRequestPayloadSchema
	if err := s.DecodeAndValidate(w, r, &requestPayload); err != nil {
		return
	}

	ctx := r.Context()
	tag, err := s.db.InsertTag(ctx, repository.InsertTagParams{
		UserID:      _userId,
		Name:        requestPayload.Name,
		Description: pgtype.Text{String: requestPayload.Description, Valid: true},
	})
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == pgerrcode.UniqueViolation {
				if pgErr.ConstraintName == "tags_name_key" {
					s.logger.Error(pgErr.Message)
					utils.ResponseError(w, http.StatusConflict, "Tag name already exists")
					return
				}
			} else if pgErr.Code == pgerrcode.ForeignKeyViolation {
				s.logger.Error(pgErr.Message)
				utils.ResponseError(w, http.StatusConflict, "User does not exist")
				return
			}
		}

		s.logger.Error(err.Error())
		utils.ResponseError(w, http.StatusInternalServerError, "Failed to create user")
		return
	}

	// responsePayload := map[string]interface{
	//    Name: tag.Name,
	//  }

	utils.Response(w, http.StatusCreated, tag)
}

func (s *Server) GetTags(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserFromContext(r)
	if !ok {
		s.logger.Error("userId not found in request context")
		utils.ResponseError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// TODO: Handle uuid parse error
	_userId, _ := uuid.Parse(userID)

	ctx := r.Context()
	tags, err := s.db.FindTags(ctx, _userId)
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

	// responsePayload := map[string]interface{}

	utils.Response(w, http.StatusOK, tags)
}

type updateTagRequestPayloadSchema struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:""`
}

func (s *Server) UpdateTag(w http.ResponseWriter, r *http.Request) {
	// userID, ok := middleware.GetUserFromContext(r)
	// if !ok {
	// 	s.logger.Error("userId not found in request context")
	// 	utils.ResponseError(w, http.StatusUnauthorized, "Unauthorized")
	// 	return
	// }
	//
	// // TODO: Handle uuid parse error
	// _userId, _ := uuid.Parse(userID)

	var requestPayload createTagRequestPayloadSchema
	if err := s.DecodeAndValidate(w, r, &requestPayload); err != nil {
		return
	}

	tagID := chi.URLParam(r, "id")
	_tagID, err := uuid.Parse(tagID)
	if err != nil {
		s.logger.Error("Invalid tag id format", slog.String("ERROR", err.Error()))
		utils.ResponseError(w, http.StatusBadRequest, "Invalid tag ID")
		return
	}

	ctx := r.Context()
	tag, err := s.db.UpdateTag(ctx, repository.UpdateTagParams{
		ID:      _tagID,
		Column2: requestPayload.Name,
		Column3: pgtype.Text{String: requestPayload.Description, Valid: true},
	})
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == pgerrcode.UniqueViolation {
				if pgErr.ConstraintName == "tags_name_key" {
					s.logger.Error(pgErr.Message)
					utils.ResponseError(w, http.StatusConflict, "Tag name already exists")
					return
				}
			} else if pgErr.Code == pgerrcode.ForeignKeyViolation {
				s.logger.Error(pgErr.Message)
				utils.ResponseError(w, http.StatusConflict, "User does not exist")
				return
			}
		}

		s.logger.Error(err.Error())
		utils.ResponseError(w, http.StatusInternalServerError, "Failed to update tag")
		return
	}

	// responsePayload := map[string]interface{
	//    Name: tag.Name,
	//  }

	utils.Response(w, http.StatusCreated, tag)
}

func (s *Server) DeleteTag(w http.ResponseWriter, r *http.Request) {
	// userID, ok := middleware.GetUserFromContext(r)
	// if !ok {
	// 	s.logger.Error("userId not found in request context")
	// 	utils.ResponseError(w, http.StatusUnauthorized, "Unauthorized")
	// 	return
	// }
	//
	// // TODO: Handle uuid parse error
	// _userId, _ := uuid.Parse(userID)

	tagID := chi.URLParam(r, "id")
	_tagID, err := uuid.Parse(tagID)
	if err != nil {
		s.logger.Error("Invalid tag id format", slog.String("ERROR", err.Error()))
		utils.ResponseError(w, http.StatusBadRequest, "Invalid tag ID")
		return
	}

	ctx := r.Context()
	err = s.db.DeleteTag(ctx, _tagID)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == pgerrcode.UniqueViolation {
				if pgErr.ConstraintName == "tags_name_key" {
					s.logger.Error(pgErr.Message)
					utils.ResponseError(w, http.StatusConflict, "Tag name already exists")
					return
				}
			} else if pgErr.Code == pgerrcode.ForeignKeyViolation {
				s.logger.Error(pgErr.Message)
				utils.ResponseError(w, http.StatusConflict, "User does not exist")
				return
			}
		}

		s.logger.Error(err.Error())
		utils.ResponseError(w, http.StatusInternalServerError, "Failed to update tag")
		return
	}

	responsePayload := map[string]interface{}{
		"message": "Tag deleted successfully",
	}

	utils.Response(w, http.StatusCreated, responsePayload)
}

type createCommandsRequestPayloadSchema struct {
	Command     string `json:"command" validate:"required"`
	Description string `json:"description" validate:""`
	TagId       string `json:"tag_id" validate:"required"`
}

func (s *Server) CreateCommand(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserFromContext(r)
	if !ok {
		s.logger.Error("userId not found in request context")
		utils.ResponseError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// TODO: Handle uuid parse error
	_userId, _ := uuid.Parse(userID)

	var requestPayload createCommandsRequestPayloadSchema
	if err := s.DecodeAndValidate(w, r, &requestPayload); err != nil {
		return
	}

	_tagId, _ := uuid.Parse(requestPayload.TagId)

	ctx := r.Context()
	tx, err := s.conn.Begin(ctx)
	if err != nil {
		s.logger.Error("Failed to start transaction: " + err.Error())
		utils.ResponseError(w, http.StatusInternalServerError, "Failed to create command")
		return
	}

	// Rollback the transaction in case of failure
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	queriesWithTx := s.db.WithTx(tx)

	command, err := queriesWithTx.InsertCommands(ctx, repository.InsertCommandsParams{
		UserID:      _userId,
		Command:     requestPayload.Command,
		Description: pgtype.Text{String: requestPayload.Description, Valid: true},
	})
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == pgerrcode.UniqueViolation {
				if pgErr.ConstraintName == "tags_name_key" {
					s.logger.Error(pgErr.Message)
					utils.ResponseError(w, http.StatusConflict, "Tag name already exists")
					return
				}
			} else if pgErr.Code == pgerrcode.ForeignKeyViolation {
				s.logger.Error(pgErr.Message)
				utils.ResponseError(w, http.StatusConflict, "User does not exist")
				return
			}
		}

		s.logger.Error(err.Error())
		utils.ResponseError(w, http.StatusInternalServerError, "Failed to create user")
		return
	}

	tag, _ := queriesWithTx.FindTagById(ctx, _tagId)

	commandTagRelation, err := queriesWithTx.AttachCommandToTag(ctx, repository.AttachCommandToTagParams{
		CommandID: command.ID,
		TagID:     tag.ID,
	})
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == pgerrcode.UniqueViolation {
				if pgErr.ConstraintName == "tags_name_key" {
					s.logger.Error(pgErr.Message)
					utils.ResponseError(w, http.StatusConflict, "Tag name already exists")
					return
				}
			} else if pgErr.Code == pgerrcode.ForeignKeyViolation {
				if pgErr.ConstraintName == "command_tags_tag_id_fkey" {
					s.logger.Error(pgErr.Message)
					utils.ResponseError(w, http.StatusConflict, "Foreign key constraint error")
					return
				}

				s.logger.Error(pgErr.Message)
				utils.ResponseError(w, http.StatusConflict, "Foreign key constraint")
				return
			}
		}

		s.logger.Error(err.Error())
		utils.ResponseError(w, http.StatusInternalServerError, "Failed to create user")
		return
	}

	if err := tx.Commit(ctx); err != nil {
		s.logger.Error("Failed to commit transaction: " + err.Error())
		utils.ResponseError(w, http.StatusInternalServerError, "Failed to create command")
		return
	}

	s.logger.Info("Command tag relation", slog.String("command", commandTagRelation.TagID.String()))

	responsePayload := map[string]interface{}{
		"command": map[string]interface{}{
			"id":          command.ID,
			"command":     command.Command,
			"description": command.Description,
		},
		"tag": map[string]interface{}{
			"id":          tag.ID,
			"name":        tag.Name,
			"description": tag.Description,
		},
	}

	utils.Response(w, http.StatusCreated, responsePayload)
}

func (s *Server) GetCommands(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserFromContext(r)
	if !ok {
		s.logger.Error("userId not found in request context")
		utils.ResponseError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// TODO: Handle uuid parse error
	_userId, _ := uuid.Parse(userID)

	ctx := r.Context()
	commands, err := s.db.FindCommandsWithTags(ctx, _userId)
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

	responsePayload := []map[string]interface{}{}
	commandMap := make(map[uuid.UUID]map[string]interface{})

	for _, row := range commands {
		// Check if the command is already in the map
		command, exists := commandMap[row.CommandID]
		if !exists {
			command = map[string]interface{}{
				"id":          row.CommandID,
				"command":     row.CommandName,
				"description": row.CommandDescription,
				"created_at":  row.CommandCreatedAt,
				"updated_at":  row.CommandUpdatedAt,
			}
			commandMap[row.CommandID] = command
		}

		// Create the tag response
		tag := map[string]interface{}{
			"id":          row.TagID,
			"name":        row.TagName,
			"description": row.TagDescription,
			"created_at":  row.TagCreatedAt,
			"updated_at":  row.TagUpdatedAt,
		}

		// Add the command and associated tag to the response
		responsePayload = append(responsePayload, map[string]interface{}{
			"command": command,
			"tag":     tag,
		})
	}

	utils.Response(w, http.StatusOK, responsePayload)
}

func (s *Server) GetCommandsWithTag(w http.ResponseWriter, r *http.Request) {
	// userID, ok := middleware.GetUserFromContext(r)
	// if !ok {
	// 	s.logger.Error("userId not found in request context")
	// 	utils.ResponseError(w, http.StatusUnauthorized, "Unauthorized")
	// 	return
	// }

	// TODO: Handle uuid parse error
	// _userId, _ := uuid.Parse(userID)

	tagID := chi.URLParam(r, "id")
	_tagID, err := uuid.Parse(tagID)
	if err != nil {
		s.logger.Error("Invalid tag id format", slog.String("ERROR", err.Error()))
		utils.ResponseError(w, http.StatusBadRequest, "Invalid tag ID")
		return
	}

	ctx := r.Context()
	commands, err := s.db.FindCommandsByTagId(ctx, _tagID)
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

	responsePayload := []map[string]interface{}{}
	commandMap := make(map[uuid.UUID]map[string]interface{})

	for _, row := range commands {
		// Check if the command is already in the map
		command, exists := commandMap[row.CommandID]
		if !exists {
			command = map[string]interface{}{
				"id":          row.CommandID,
				"command":     row.CommandName,
				"description": row.CommandDescription,
				"created_at":  row.CommandCreatedAt,
				"updated_at":  row.CommandUpdatedAt,
			}
			commandMap[row.CommandID] = command
		}

		// Create the tag response
		tag := map[string]interface{}{
			"id":          row.TagID,
			"name":        row.TagName,
			"description": row.TagDescription,
			"created_at":  row.TagCreatedAt,
			"updated_at":  row.TagUpdatedAt,
		}

		// Add the command and associated tag to the response
		responsePayload = append(responsePayload, map[string]interface{}{
			"command": command,
			"tag":     tag,
		})
	}

	utils.Response(w, http.StatusOK, responsePayload)
}

type updateCommandRequestPayloadSchema struct {
	Command     string `json:"command" validate:"required"`
	Description string `json:"description" validate:""`
}

func (s *Server) UpdateCommand(w http.ResponseWriter, r *http.Request) {
	// userID, ok := middleware.GetUserFromContext(r)
	// if !ok {
	// 	s.logger.Error("userId not found in request context")
	// 	utils.ResponseError(w, http.StatusUnauthorized, "Unauthorized")
	// 	return
	// }
	//
	// // TODO: Handle uuid parse error
	// _userId, _ := uuid.Parse(userID)

	var requestPayload updateCommandRequestPayloadSchema
	if err := s.DecodeAndValidate(w, r, &requestPayload); err != nil {
		return
	}

	tagID := chi.URLParam(r, "id")
	_tagID, err := uuid.Parse(tagID)
	if err != nil {
		s.logger.Error("Invalid tag id format", slog.String("ERROR", err.Error()))
		utils.ResponseError(w, http.StatusBadRequest, "Invalid tag ID")
		return
	}

	ctx := r.Context()
	tag, err := s.db.UpdateCommand(ctx, repository.UpdateCommandParams{
		ID:      _tagID,
		Column2: requestPayload.Command,
		Column3: pgtype.Text{String: requestPayload.Description, Valid: true},
	})
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == pgerrcode.UniqueViolation {
				if pgErr.ConstraintName == "tags_name_key" {
					s.logger.Error(pgErr.Message)
					utils.ResponseError(w, http.StatusConflict, "Tag name already exists")
					return
				}
			} else if pgErr.Code == pgerrcode.ForeignKeyViolation {
				s.logger.Error(pgErr.Message)
				utils.ResponseError(w, http.StatusConflict, "User does not exist")
				return
			}
		}

		s.logger.Error(err.Error())
		utils.ResponseError(w, http.StatusInternalServerError, "Failed to update tag")
		return
	}

	// responsePayload := map[string]interface{
	//    Name: tag.Name,
	//  }

	utils.Response(w, http.StatusCreated, tag)
}

func (s *Server) DeleteCommand(w http.ResponseWriter, r *http.Request) {
	// userID, ok := middleware.GetUserFromContext(r)
	// if !ok {
	// 	s.logger.Error("userId not found in request context")
	// 	utils.ResponseError(w, http.StatusUnauthorized, "Unauthorized")
	// 	return
	// }
	//
	// // TODO: Handle uuid parse error
	// _userId, _ := uuid.Parse(userID)

	commandID := chi.URLParam(r, "id")
	_commandId, err := uuid.Parse(commandID)
	if err != nil {
		s.logger.Error("Invalid tag id format", slog.String("ERROR", err.Error()))
		utils.ResponseError(w, http.StatusBadRequest, "Invalid tag ID")
		return
	}

	ctx := r.Context()
	tx, err := s.conn.Begin(ctx)
	if err != nil {
		s.logger.Error("Failed to start transaction: " + err.Error())
		utils.ResponseError(w, http.StatusInternalServerError, "Failed to create command")
		return
	}

	// Rollback the transaction in case of failure
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	queriesWithTx := s.db.WithTx(tx)

	_ = queriesWithTx.DeleteCommandTagRelationByCommandId(ctx, _commandId)

	err = queriesWithTx.DeleteCommand(ctx, _commandId)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == pgerrcode.UniqueViolation {
				if pgErr.ConstraintName == "tags_name_key" {
					s.logger.Error(pgErr.Message)
					utils.ResponseError(w, http.StatusConflict, "Tag name already exists")
					return
				}
			} else if pgErr.Code == pgerrcode.ForeignKeyViolation {
				s.logger.Error(pgErr.Message)
				utils.ResponseError(w, http.StatusConflict, "User does not exist")
				return
			}
		}

		s.logger.Error(err.Error())
		utils.ResponseError(w, http.StatusInternalServerError, "Failed to update tag")
		return
	}

	if err := tx.Commit(ctx); err != nil {
		s.logger.Error("Failed to commit transaction: " + err.Error())
		utils.ResponseError(w, http.StatusInternalServerError, "Failed to create command")
		return
	}

	responsePayload := map[string]interface{}{
		"message": "Command deleted successfully",
	}

	utils.Response(w, http.StatusCreated, responsePayload)
}
