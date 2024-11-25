package router

import (
	"database/sql"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/template/html/v2"

	"01-Login/platform/authenticator"
	"01-Login/platform/middleware"
	"01-Login/platform/router/handlers"
)

func New(auth *authenticator.Authenticator, dbConn *sql.DB) *fiber.App {
	if dbConn == nil {
		log.Fatal("Database connection is nil")
	}

	// Initialize template engine
	engine := html.New("./web/template", ".html")

	// Create a new Fiber app
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	// Initialize session store
	store := session.New(session.Config{
		KeyLookup:  "cookie:auth-session",
		Expiration: 24 * time.Hour,
	})

	// Serve static files
	app.Static("/public", "./web/static")

	// Routes
	app.Get("/", handlers.HandleHome)
	app.Get("/login", handlers.HandleLogin(auth))
	app.Get("/callback", handlers.HandleCallback(auth, dbConn, store))
	app.Get("/user", middleware.IsAuthenticated(store), handlers.HandleUser(store))
	app.Get("/logout", handlers.HandleLogout(store))

	return app
}
