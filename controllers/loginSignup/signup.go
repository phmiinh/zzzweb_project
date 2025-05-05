package loginSignup

import (
	"Zzz_project/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"time"
)

var DB *gorm.DB

func RegisterHandler(c *fiber.Ctx) error {
	var info models.Info

	if err := c.BodyParser(&info); err != nil {
		return c.Render("login_signup", fiber.Map{
			"ErrorSignup": "Failed to parse request body",
		})
	}

	var existingUser models.User
	if err := models.DB.Where("email = ?", info.Email).First(&existingUser).Error; err == nil {
		return c.Render("login_signup", fiber.Map{
			"ErrorSignup": "Email already registered",
		})
	}

	if len(info.Password) < 8 {
		return c.Render("login_signup", fiber.Map{
			"ErrorSignup": "Password must be at least 8 characters",
		})
	}

	if len(info.Address) < 10 {
		return c.Render("login_signup", fiber.Map{
			"ErrorSignup": "Address must be at least 10 characters",
		})
	}

	hashedPassword, err := HashPassword(info.Password)
	if err != nil {
		return c.Render("login_signup", fiber.Map{
			"ErrorSignup": "Error hashing password",
		})
	}

	birthdate, err := time.Parse("02/01/2006", info.Birthdate)
	if err != nil {
		return c.Render("login_signup", fiber.Map{
			"ErrorSignup": "Invalid birthdate format (dd/mm/yyyy required)",
		})
	}

	user := models.User{
		Email:     info.Email,
		Username:  info.Username,
		Password:  hashedPassword,
		Role:      "customer",
		Address:   info.Address,
		Birthdate: birthdate,
	}

	if err := models.DB.Create(&user).Error; err != nil {
		return c.Render("login_signup", fiber.Map{
			"ErrorSignup": "Error creating user",
		})
	}

	return c.Render("login_signup", fiber.Map{
		"SuccessSignup": "Sign up successful. Please log in!",
		"isRegistered":  true,
	})
}
