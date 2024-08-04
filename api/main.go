package main

import (
	"log"
	"net/http"

	"github.com/endalk200/termflow-api/handlers"
)

func main() {
	router := http.NewServeMux()

	coreRouter := http.NewServeMux()
	coreRouter.HandleFunc("GET /health", handlers.ServerHealthHandler)

	// authRouter := http.NewServeMux()
	// authRouter.HandleFunc("POST /login", handlers.LoginHandler)
	// authRouter.HandleFunc("POST /signup", handlers.SignupHandler)
	//
	// userRouter := http.NewServeMux()
	// userRouter.HandleFunc("GET /me", handlers.GetSessionHandler)

	// router.Handle("/auth/", http.StripPrefix("/auth", authRouter))
	// router.Handle("/user/", http.StripPrefix("/user", AuthMiddleware(userRouter)))

	router.Handle("/", coreRouter)
	stack := CreateStack(LoggingMiddleware)

	server := http.Server{
		Addr:    ":3000",
		Handler: stack(router),
	}

	log.Printf("Server listening on port 3000...")
	err := server.ListenAndServe()
	if err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
