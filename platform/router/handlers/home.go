package handlers

import "github.com/gofiber/fiber/v2"

// Handlers
func HandleHome(c *fiber.Ctx) error {
	return c.Render("home", fiber.Map{})
}
