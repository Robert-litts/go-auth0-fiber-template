// package router

// import (
// 	"crypto/rand"
// 	"database/sql"
// 	"encoding/base64"
// 	"encoding/gob"
// 	"fmt"
// 	"log"
// 	"net/url"
// 	"os"
// 	"strings"
// 	"time"

// 	"github.com/gofiber/fiber/v2"
// 	"github.com/gofiber/fiber/v2/middleware/session"
// 	"github.com/gofiber/template/html/v2"
// 	"github.com/google/uuid"
// 	_ "github.com/jackc/pgx/v5/stdlib"

// 	"01-Login/internal/db"
// 	"01-Login/platform/authenticator"
// 	"01-Login/platform/middleware"
// 	"01-Login/platform/router/handlers"
// )

// func init() {
// 	// Register necessary types with gob
// 	gob.Register(map[string]interface{}{})
// 	gob.Register(time.Time{}) // Register time.Time
// }

// // New registers the routes and returns the router.
// func New(auth *authenticator.Authenticator, dbConn *sql.DB) *fiber.App {
// 	if dbConn == nil {
// 		log.Fatal("Database connection is nil")
// 	}

// 	// Register the map type for session storage
// 	gob.Register(map[string]interface{}{})

// 	// Initialize template engine
// 	engine := html.New("./web/template", ".html")

// 	// Create a new Fiber app
// 	app := fiber.New(fiber.Config{
// 		Views: engine,
// 	})

// 	// Initialize session store
// 	store := session.New(session.Config{
// 		KeyLookup:  "cookie:auth-session",
// 		Expiration: 24 * time.Hour,
// 	})

// 	// Serve static files
// 	app.Static("/public", "./web/static")

// 	// Routes
// 	app.Get("/", handleHome)
// 	app.Get("/login", handleLogin(auth))
// 	app.Get("/callback", handleCallback(auth, dbConn, store))
// 	app.Get("/user", middleware.IsAuthenticated(store), handleUser(store))
// 	app.Get("/logout", handleLogout(store))

// 	return app
// }

// // Handlers
// func handleHome(c *fiber.Ctx) error {
// 	return c.Render("home", fiber.Map{})
// }

// func handleLogin(auth *authenticator.Authenticator) fiber.Handler {
// 	return func(c *fiber.Ctx) error {
// 		state := generateRandomState()
// 		return c.Redirect(auth.AuthCodeURL(state))
// 	}
// }

// // func handleCallback(auth *authenticator.Authenticator, store *session.Store) fiber.Handler {
// // 	return func(c *fiber.Ctx) error {
// // 		sess, err := store.Get(c)
// // 		if err != nil {
// // 			return c.Status(fiber.StatusInternalServerError).SendString("Failed to get session")
// // 		}

// // 		code := c.Query("code")
// // 		token, err := auth.Exchange(c.Context(), code)
// // 		if err != nil {
// // 			return c.Status(fiber.StatusUnauthorized).SendString("Failed to exchange token")
// // 		}

// // 		idToken, err := auth.VerifyIDToken(c.Context(), token)
// // 		if err != nil {
// // 			return c.Status(fiber.StatusUnauthorized).SendString("Failed to verify ID token")
// // 		}

// // 		var profile map[string]interface{}
// // 		if err := idToken.Claims(&profile); err != nil {
// // 			return c.Status(fiber.StatusInternalServerError).SendString("Failed to get user info")
// // 		}

// // 		sess.Set("profile", profile)
// // 		if err := sess.Save(); err != nil {
// // 			return c.Status(fiber.StatusInternalServerError).SendString("Failed to save session")
// // 		}

// // 		return c.Redirect("/user")
// // 	}
// // }

// func handleCallback(auth *authenticator.Authenticator, dbConn *sql.DB, store *session.Store) fiber.Handler {
// 	q := db.New(dbConn)

// 	return func(c *fiber.Ctx) error {
// 		// Get auth code and exchange for token
// 		code := c.Query("code")
// 		if code == "" {
// 			log.Println("No authorization code received")
// 			return c.Status(fiber.StatusBadRequest).SendString("Missing authorization code")
// 		}

// 		// Exchange the auth code for a token
// 		token, err := auth.Exchange(c.Context(), code)
// 		if err != nil {
// 			log.Printf("Token exchange error: %v", err)
// 			return c.Status(fiber.StatusInternalServerError).SendString("Failed to exchange token")
// 		}

// 		// Verify ID token
// 		idToken, err := auth.VerifyIDToken(c.Context(), token)
// 		if err != nil {
// 			log.Printf("Token verification error: %v", err)
// 			return c.Status(fiber.StatusUnauthorized).SendString("Failed to verify ID token")
// 		}

// 		// Extract user profile with detailed logging
// 		var profile map[string]interface{}
// 		if err := idToken.Claims(&profile); err != nil {
// 			log.Printf("Claims extraction error: %v", err)
// 			return c.Status(fiber.StatusInternalServerError).SendString("Failed to get user info")
// 		}

// 		// Debug log the profile contents
// 		log.Printf("Received profile data: %+v", profile)

// 		// Safely extract profile information
// 		auth0ID, ok := profile["sub"].(string)
// 		if !ok {
// 			log.Printf("Auth0 ID (sub) missing or invalid in profile: %+v", profile)
// 			return c.Status(fiber.StatusInternalServerError).SendString("Invalid profile: missing sub claim")
// 		}

// 		email, ok := profile["email"].(string)
// 		if !ok {
// 			// Try to get email from alternative fields
// 			if emailInterface, exists := profile["name"].(string); exists {
// 				email = emailInterface
// 			} else {
// 				log.Printf("Email missing or invalid in profile: %+v", profile)
// 				// Use auth0ID as email if no email is available
// 				email = auth0ID
// 			}
// 		}

// 		// Database operations
// 		var userID uuid.UUID
// 		var userEmail string
// 		var userAuth0ID string
// 		var userCreatedAt sql.NullTime

// 		// First try to get existing user
// 		existingUser, err := q.GetUserByAuth0ID(c.Context(), auth0ID)
// 		if err != nil {
// 			if err == sql.ErrNoRows {
// 				// User doesn't exist, create new user
// 				log.Printf("Creating new user with auth0ID: %s and email: %s", auth0ID, email)
// 				createParams := db.CreateUserParams{
// 					Auth0ID: auth0ID,
// 					Email:   email,
// 				}
// 				newUser, err := q.CreateUser(c.Context(), createParams)
// 				if err != nil {
// 					log.Printf("User creation error: %v", err)
// 					return c.Status(fiber.StatusInternalServerError).SendString("Failed to create user")
// 				}
// 				userID = newUser.ID
// 				userEmail = newUser.Email
// 				userAuth0ID = newUser.Auth0ID
// 				userCreatedAt = newUser.CreatedAt
// 				log.Printf("Successfully created new user with ID: %s", userID)
// 			} else {
// 				log.Printf("Database error when checking for existing user: %v", err)
// 				return c.Status(fiber.StatusInternalServerError).SendString("Database error")
// 			}
// 		} else {
// 			log.Printf("Found existing user with ID: %s", existingUser.ID)
// 			userID = existingUser.ID
// 			userEmail = existingUser.Email
// 			userAuth0ID = existingUser.Auth0ID
// 			userCreatedAt = existingUser.CreatedAt
// 		}

// 		// Session handling
// 		sess, err := store.Get(c)
// 		if err != nil {
// 			log.Printf("Session retrieval error: %v", err)
// 			return c.Status(fiber.StatusInternalServerError).SendString("Failed to get session")
// 		}

// 		// Create session profile
// 		sessionProfile := map[string]interface{}{
// 			"auth0":      profile,
// 			"user_id":    userID.String(),
// 			"email":      userEmail,
// 			"auth0_id":   userAuth0ID,
// 			"created_at": userCreatedAt.Time,
// 		}

// 		log.Printf("Setting session profile: %+v", sessionProfile)

// 		sess.Set("profile", sessionProfile)
// 		if err := sess.Save(); err != nil {
// 			log.Printf("Session save error: %v", err)
// 			return c.Status(fiber.StatusInternalServerError).SendString("Failed to save session")
// 		}

// 		log.Println("Successfully completed callback handling")
// 		return c.Redirect("/user")
// 	}
// }

// func handleUser(store *session.Store) fiber.Handler {
// 	return func(c *fiber.Ctx) error {
// 		sess, err := store.Get(c)
// 		if err != nil {
// 			return c.Redirect("/")
// 		}

// 		profile := sess.Get("profile")
// 		return c.Render("user", fiber.Map{
// 			"profile": profile,
// 		})
// 	}
// }

// // func handleLogout(store *session.Store) fiber.Handler {
// // 	return func(c *fiber.Ctx) error {
// // 		sess, err := store.Get(c)
// // 		if err != nil {
// // 			return c.Redirect("/")
// // 		}

// // 		if err := sess.Destroy(); err != nil {
// // 			return c.Status(fiber.StatusInternalServerError).SendString("Failed to destroy session")
// // 		}

// // 		// Construct Auth0 logout URL
// // 		domain := os.Getenv("AUTH0_DOMAIN")
// // 		clientID := os.Getenv("AUTH0_CLIENT_ID")
// // 		returnTo := os.Getenv("AUTH0_CALLBACK_URL")

// // 		logoutURL := fmt.Sprintf(
// // 			"https://%s/v2/logout?client_id=%s&returnTo=%s",
// // 			domain,
// // 			clientID,
// // 			returnTo,
// // 		)

// // 		return c.Redirect(logoutURL)
// // 	}
// // }

// func handleLogout(store *session.Store) fiber.Handler {
// 	return func(c *fiber.Ctx) error {
// 		// Destroy the session
// 		sess, err := store.Get(c)
// 		if err != nil {
// 			return c.Redirect("/")
// 		}

// 		if err := sess.Destroy(); err != nil {
// 			return c.Status(fiber.StatusInternalServerError).SendString("Failed to destroy session")
// 		}

// 		// Parse the Auth0 logout URL
// 		logoutUrl, err := url.Parse("https://" + os.Getenv("AUTH0_DOMAIN") + "/v2/logout")
// 		if err != nil {
// 			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
// 		}

// 		// Determine the scheme (http or https)
// 		scheme := "http"
// 		if c.Protocol() == "https" {
// 			scheme = "https"
// 		}

// 		// Create the returnTo URL
// 		returnTo, err := url.Parse(scheme + "://" + c.Hostname())
// 		if err != nil {
// 			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
// 		}

// 		// Add query parameters
// 		parameters := url.Values{}
// 		parameters.Add("returnTo", returnTo.String())
// 		parameters.Add("client_id", os.Getenv("AUTH0_CLIENT_ID"))
// 		logoutUrl.RawQuery = parameters.Encode()

// 		return c.Redirect(logoutUrl.String(), fiber.StatusTemporaryRedirect)
// 	}
// }

// // Helper function to generate random state
// func generateRandomState() string {
// 	b := make([]byte, 32)
// 	rand.Read(b)
// 	return base64.StdEncoding.EncodeToString(b)
// }

// // ConnectDB establishes a connection to the database
// func ConnectDB() *sql.DB {
// 	dsn := os.Getenv("DB_URL")
// 	if dsn == "" {
// 		log.Fatal("DB_URL environment variable not set")
// 	}

// 	db, err := sql.Open("pgx", dsn)
// 	if err != nil {
// 		log.Fatalf("Failed to connect to the database: %v", err)
// 	}

// 	// Verify connection
// 	if err := db.Ping(); err != nil {
// 		log.Fatalf("Failed to ping the database: %v", err)
// 	}

// 	log.Println("Successfully connected to the database")
// 	return db
// }

// func extractToken(c *fiber.Ctx) (string, error) {
// 	authHeader := c.Get("Authorization")
// 	if authHeader == "" {
// 		return "", fmt.Errorf("Authorization header is missing")
// 	}

// 	// Expect the format "Bearer <token>"
// 	parts := strings.Split(authHeader, " ")
// 	if len(parts) != 2 || parts[0] != "Bearer" {
// 		return "", fmt.Errorf("Authorization header format must be Bearer <token>")
// 	}

// 	return parts[1], nil
// }

//////////////////////////////////////////

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
