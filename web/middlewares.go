package main

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	jwtutils "github.com/endalk200/http-server/utils"
)

type Middleware func(http.Handler) http.Handler

func CreateStack(xs ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for i := len(xs) - 1; i >= 0; i-- {
			x := xs[i]
			next = x(next)
		}

		return next
	}
}

const AuthUserId = "middleware.auth.userID"

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rawAuthorizationHeader := r.Header.Get("Authorization")
		if rawAuthorizationHeader == "" && len(strings.Split(rawAuthorizationHeader, " ")) != 2 {
			log.Printf("No Bearer token provided")
			http.Error(w, "Unauthorized API call", http.StatusUnauthorized)
			return
		}

		parsedAuthToken := strings.Split(rawAuthorizationHeader, " ")[1]
		if parsedAuthToken == "" {
			log.Printf("Unauthorized API call. Authorization token is not provided")
			http.Error(w, "Unauthorized API call", http.StatusUnauthorized)
			return
		}

		token, err := jwtutils.VerifyJWT(parsedAuthToken)
		if err != nil {
			log.Printf("Invalid JWT is provided: %s", err)
			http.Error(w, "Unauthorized API call", http.StatusUnauthorized)
			return
		}

		userId, _ := token.Claims.GetSubject()
		ctx := context.WithValue(r.Context(), AuthUserId, userId)
		req := r.WithContext(ctx)

		next.ServeHTTP(w, req)
	})
}

type WrappedWritter struct {
	http.ResponseWriter
	statusCode int
}

func (w *WrappedWritter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.statusCode = statusCode
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		wrapped := &WrappedWritter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(wrapped, r)

		log.Println(wrapped.statusCode, r.Method, r.URL.Path, time.Since(start))
	})
}
