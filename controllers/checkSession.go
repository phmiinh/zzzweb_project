package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

func CheckSession(c *fiber.Ctx, store *session.Store) (interface{}, error) {
	sess, err := store.Get(c)
	if err != nil {
		return nil, c.Status(fiber.StatusInternalServerError).SendString("Session error")
	}

	userID := sess.Get("user_id")
	if userID == nil {
		return nil, c.Redirect("/auth/login")
	}

	return userID, nil
}
