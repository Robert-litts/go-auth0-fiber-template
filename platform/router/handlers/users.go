package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

func HandleUser(store *session.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		sess, err := store.Get(c)
		if err != nil {
			return c.Redirect("/")
		}

		profile := sess.Get("profile")
		return c.Render("user", fiber.Map{
			"profile": profile,
		})
	}
}
