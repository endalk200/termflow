package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/endalk200/termflow-api/internal/repository"
	"github.com/endalk200/termflow-api/pkgs/config"
	_ "github.com/joho/godotenv/autoload"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct {
	port   int
	cfg    config.AppConfig
	logger *slog.Logger
	db     *repository.Queries
}

func NewServer(logger *slog.Logger) *http.Server {
	var cfg config.AppConfig
	err := config.LoadConfig(&cfg)

	databaseConnectionUri, err := config.ConstructDatabaseUrl(cfg)
	if err != nil {
		logger.Error("Error while constructing DATABASE_URI", slog.String("ERROR", err.Error()))
	}

	connection, err := pgxpool.New(context.Background(), databaseConnectionUri)
	if err != nil {
		logger.Error("Unable to create connection pool", slog.String("ERROR", err.Error()))
	}

	queries := repository.New(connection)

	NewServer := &Server{
		port:   cfg.ApplicationPort,
		logger: logger,
		cfg:    cfg,
		db:     queries,
	}

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	logger.Info("Server listening", slog.String("addr", ":8080"))

	return server
}
