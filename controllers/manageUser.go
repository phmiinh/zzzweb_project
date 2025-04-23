package controllers

import (
	"Zzz_project/models"
	"github.com/gofiber/fiber/v2"
)

func UpdateUserRole(c *fiber.Ctx) error {
	type RoleUpdateRequest struct {
		UserID uint   `json:"user_id"`
		Role   string `json:"role"`
	}

	var req RoleUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Dữ liệu không hợp lệ"})
	}

	var user models.User
	if err := models.DB.First(&user, req.UserID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Người dùng không tồn tại"})
	}

	user.Role = req.Role
	if err := models.DB.Save(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Không thể cập nhật vai trò"})
	}

	return c.JSON(fiber.Map{"message": "Cập nhật vai trò thành công"})
}

// Xóa người dùng
func DeleteUser(c *fiber.Ctx) error {
	type DeleteRequest struct {
		UserID uint `json:"user_id"`
	}

	var req DeleteRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Dữ liệu không hợp lệ"})
	}

	var user models.User
	if err := models.DB.First(&user, req.UserID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Người dùng không tồn tại"})
	}

	if err := models.DB.Delete(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Không thể xóa người dùng"})
	}

	return c.JSON(fiber.Map{"message": "Xóa người dùng thành công"})
}
