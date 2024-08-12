package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/endalk200/termflow-api/internal/database"
	_ "github.com/joho/godotenv/autoload"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct {
	port int

	db *database.Queries
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))

	dbConnectionString := "postgres://admin:admin@localhost:5432/my-db?sslmode=disable"

	connection, err := pgxpool.New(context.Background(), dbConnectionString)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v", err)
	}

	queries := database.New(connection)

	NewServer := &Server{
		port: port,
		db:   queries,
	}

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
