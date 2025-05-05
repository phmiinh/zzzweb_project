package userInterface

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"time"
)

type FormData struct {
	Email     string    `json:"email"`
	Address   string    `json:"address"`
	Birthdate time.Time `json:"birthdate"`
}

var validate = validator.New()

func GetFormData(c *fiber.Ctx) (*FormData, error) {
	var form struct {
		Email     string `form:"email" validate:"required,email"`
		Address   string `form:"address" validate:"required,min=10"`
		Birthdate string `form:"birthdate" validate:"required"`
	}

	if err := c.BodyParser(&form); err != nil {
		return nil, err
	}

	if err := validate.Struct(form); err != nil {
		return nil, err
	}

	birthdate, err := time.Parse("2006-01-02", form.Birthdate)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Invalid birthdate format, use YYYY-MM-DD")
	}

	formData := &FormData{
		Email:     form.Email,
		Address:   form.Address,
		Birthdate: birthdate,
	}

	return formData, nil
}
