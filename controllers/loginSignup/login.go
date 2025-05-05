package loginSignup

import (
	"Zzz_project/models"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"golang.org/x/crypto/bcrypt"
)

var Store *session.Store

func LoginHandler(c *fiber.Ctx) error {
	var info models.Info
	if err := c.BodyParser(&info); err != nil {
		return err
	}

	email := info.Email
	password := info.Password
	var user models.User

	result := models.DB.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return c.Render("login_signup", fiber.Map{
			"ErrorLogin": "Invalid email or password",
		})
	}

	if !checkPasswordHash(password, user.Password) {
		return c.Render("login_signup", fiber.Map{
			"ErrorLogin": "Invalid email or password",
		})
	}

	sess, err := Store.Get(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Session error")
	}

	sess.Set("user_id", user.ID)
	sess.Set("role", user.Role)

	if err := sess.Save(); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to save session")
	}

	if user.Role == "admin" {
		return c.Redirect("/admin/dashboard")
	}
	return c.Redirect("/auth/home")
}
func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GetLogin(c *fiber.Ctx) error {
	return c.Render("login_signup", nil)
}
