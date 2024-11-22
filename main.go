package main

import (
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"

	"01-Login/platform/authenticator"
	"01-Login/platform/router"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Failed to load the env vars: %v", err)
	}

	// Initialize authenticator
	auth, err := authenticator.New()
	if err != nil {
		log.Fatalf("Failed to initialize the authenticator: %v", err)
	}

	// Initialize database connection
	dbConn := router.ConnectDB()
	if dbConn == nil {
		log.Fatalf("Failed to initialize database connection")
	}

	// Initialize Fiber app with router
	app := router.New(auth, dbConn)

	// Configure Fiber settings
	//app.Config().DisableStartupMessage = false

	// Start server
	log.Print("Server listening on http://localhost:3000/")
	if err := app.Listen("0.0.0.0:3000"); err != nil {
		log.Fatalf("There was an error with the http server: %v", err)
	}
}
