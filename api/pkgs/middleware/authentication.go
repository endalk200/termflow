package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"strings"

	"github.com/endalk200/termflow-api/pkgs/auth"
	"github.com/endalk200/termflow-api/pkgs/utils"
)

type contextKey string

const userContextKey contextKey = "userId"

func GetUserFromContext(r *http.Request) (string, bool) {
	userID, ok := r.Context().Value(userContextKey).(string)

	return userID, ok
}

func Authentication(logger *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				logger.Error("Missing Authorization header")
				utils.ResponseError(w, http.StatusUnauthorized, "Missing authorization header")
				return
			}

			// Split the header into type and token (e.g., "Bearer <token>")
			tokenParts := strings.Split(authHeader, " ")
			if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
				logger.Error("Invalid authorization header format", slog.String("header", authHeader))
				utils.ResponseError(w, http.StatusUnauthorized, "Invalid authorization header format")
				return
			}

			token := tokenParts[1]
			verifiedToken, err := auth.VerifyJWT(token)
			if err != nil {
				logger.Error("Invalid jwt", slog.String("ERROR", err.Error()))
				utils.ResponseError(w, http.StatusUnauthorized, "Invalid authorization header format")
				return
			}

			userId, _ := verifiedToken.Claims.GetSubject()

			ctx := context.WithValue(r.Context(), userContextKey, userId)
			next.ServeHTTP(w, r.WithContext(ctx))

			// next.ServeHTTP(w, r)
		})
	}
}
