package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/B33Boy/Judgement/internal/app"

	_ "github.com/joho/godotenv/autoload"
)

func NewServer(app *app.App) *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))

	// Declare Server config
	return &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      app.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
}
