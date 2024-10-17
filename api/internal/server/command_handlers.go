package server

import (
	"log/slog"
	"net/http"

	"github.com/endalk200/termflow-api/internal/repository"
	"github.com/endalk200/termflow-api/pkgs/middleware"
	"github.com/endalk200/termflow-api/pkgs/utils"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

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
