package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/endalk200/termflow-api/internal/server"
)

func main() {
	server := server.NewServer(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	err := server.ListenAndServe()
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}
}
