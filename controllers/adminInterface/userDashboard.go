package adminInterface

import (
	"Zzz_project/controllers/loginSignup"
	"Zzz_project/models"
	"github.com/gofiber/fiber/v2"
	"time"
)

func UsersList(c *fiber.Ctx) error {
	// Lấy session
	sess, err := loginSignup.Store.Get(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Session error")
	}

	// Lấy role từ session
	role, ok := sess.Get("role").(string)
	if !ok || role != "admin" {
		return c.Status(fiber.StatusUnauthorized).SendString("Không có quyền truy cập")
	}

	// Truy vấn danh sách người dùng (role = "customer")
	var users []models.User
	if err := models.DB.Find(&users).Error; err != nil {
		return c.Status(500).SendString("Lỗi khi truy vấn người dùng")
	}

	return c.Render("admin/users", fiber.Map{
		"Users": users,
	}, "layouts/admin")
}

func UserDetail(c *fiber.Ctx) error {
	id := c.Params("id")
	var user models.User
	if err := models.DB.First(&user, id).Error; err != nil {
		return c.Status(404).SendString("Không tìm thấy người dùng")
	}

	return c.Render("admin/user_detail", fiber.Map{
		"User": user,
	}, "layouts/admin")
}

func ChangeUserRole(c *fiber.Ctx) error {
	id := c.Params("id")
	var user models.User
	if err := models.DB.First(&user, id).Error; err != nil {
		return c.Status(404).SendString("Không tìm thấy người dùng")
	}

	newRole := c.FormValue("role")
	if newRole != "admin" && newRole != "customer" {
		return c.Status(400).SendString("Vai trò không hợp lệ")
	}

	user.Role = newRole
	models.DB.Save(&user)

	return c.Redirect("/admin/users")
}
func DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	var user models.User

	if err := models.DB.First(&user, id).Error; err != nil {
		return c.Status(404).SendString("Người dùng không tồn tại")
	}

	if err := models.DB.Delete(&user).Error; err != nil {
		return c.Status(500).SendString("Lỗi khi xoá người dùng")
	}

	return c.Redirect("/admin/users")
}
func UpdateUser(c *fiber.Ctx) error {
	id := c.Params("id")
	var user models.User
	if err := models.DB.First(&user, id).Error; err != nil {
		return c.Status(404).SendString("Người dùng không tồn tại")
	}

	// Lấy dữ liệu form
	user.Email = c.FormValue("email")
	user.Username = c.FormValue("username")
	user.Address = c.FormValue("address")
	user.Role = c.FormValue("role")

	// Chuyển đổi ngày sinh
	birthdateStr := c.FormValue("birthdate")
	birthdate, err := time.Parse("2006-01-02", birthdateStr)
	if err != nil {
		return c.Status(400).SendString("Ngày sinh không hợp lệ")
	}
	user.Birthdate = birthdate

	if err := models.DB.Save(&user).Error; err != nil {
		return c.Status(500).SendString("Lỗi khi cập nhật người dùng")
	}

	return c.Redirect("/admin/users")
}
