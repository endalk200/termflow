package server

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/endalk200/termflow-api/pkgs/utils"
	"github.com/go-playground/validator/v10"
)

func (s *Server) DecodeAndValidate(w http.ResponseWriter, r *http.Request, payload interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(payload); err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "Invalid request payload")
		return err
	}

	if customErrors, err := utils.ValidateAndFormatErrors(payload); err != nil {
		s.logger.Error("Error during request payload validation", slog.String("ERROR", err.Error()))

		utils.Response(w, http.StatusBadRequest, customErrors)
		return err
	}

	return nil
}

var Validate *validator.Validate
