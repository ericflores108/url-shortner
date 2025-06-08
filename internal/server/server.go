package server

import (
	"fmt"
	"net/http"
	"os"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"urlshort/internal/database"
)

type Server struct {
	baseURL string

	db database.Service
}

func NewServer() *http.Server {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = fmt.Sprintf("http://localhost:%d", port)
	}
	server := &Server{
		baseURL: baseURL,

		db: database.New(),
	}

	// Declare Server config
	httpServer := &http.Server{
		Addr:         ":" + port,
		Handler:      server.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return httpServer
}
