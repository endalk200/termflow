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
