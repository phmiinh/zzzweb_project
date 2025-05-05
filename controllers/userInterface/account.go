package userInterface

import (
	"Zzz_project/controllers"
	"Zzz_project/controllers/loginSignup"
	"Zzz_project/models"
	"github.com/gofiber/fiber/v2"
)

func AccountHandler(c *fiber.Ctx) error {
	userID, err := controllers.CheckSession(c, loginSignup.Store)
	if err != nil {
		return err
	}

	var user models.User
	result := models.DB.First(&user, userID)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).SendString("User not found")
	}

	var birthdateStr string
	if !user.Birthdate.IsZero() {
		birthdateStr = user.Birthdate.Format("2006-01-02")
	}

	return c.Render("accountfunc/account", fiber.Map{
		"Username":  user.Username,
		"Email":     user.Email,
		"Address":   user.Address,
		"Birthdate": birthdateStr,
	}, "layouts/main")
}

func UpdateAccountHandler(c *fiber.Ctx) error {
	userID, err := controllers.CheckSession(c, loginSignup.Store)
	if err != nil {
		return err
	}

	var user models.User
	if err := models.DB.First(&user, userID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).SendString("User not found")
	}

	formData, err := GetFormData(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	user.Email = formData.Email
	user.Address = formData.Address
	user.Birthdate = formData.Birthdate

	if err := models.DB.Save(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to update account")
	}

	return c.Redirect("/auth/account")
}

func ChangePasswordViewHandler(c *fiber.Ctx) error {
	return c.Render("accountfunc/change_password", nil, "layouts/main")
}

func ChangePasswordHandler(c *fiber.Ctx) error {
	userID, err := controllers.CheckSession(c, loginSignup.Store)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error":   "Unauthorized",
		})
	}

	//lấy thông tin người dùng từ DB
	var user models.User
	if err := models.DB.First(&user, userID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error":   "User not found",
		})
	}

	//parse dữ liệu JSON từ request
	var req struct {
		CurrentPassword string `json:"current_password"`
		NewPassword     string `json:"new_password"`
		ConfirmPassword string `json:"confirm_password"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid request data",
		})
	}

	//kiểm tra mật khẩu cũ có đúng không
	if !loginSignup.CheckPasswordHash(req.CurrentPassword, user.Password) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Incorrect current password",
		})
	}

	//mật khẩu mới khác mật khẩu cũ
	if req.NewPassword == req.CurrentPassword {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "New password cannot be the same as the current password",
		})
	}

	//xác nhận mật khẩu mới
	if req.NewPassword != req.ConfirmPassword {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "New passwords do not match",
		})
	}

	if len(req.NewPassword) < 8 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "New password must be at least 8 characters",
		})
	}

	//hash mật khẩu mới và cập nhật vào DB
	hashedPassword, err := loginSignup.HashPassword(req.NewPassword)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Error hashing new password",
		})
	}

	user.Password = hashedPassword
	if err := models.DB.Save(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to update password",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Password changed successfully",
	})
}
