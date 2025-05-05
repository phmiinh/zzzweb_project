package userInterface

import (
	"Zzz_project/controllers/loginSignup"
	"Zzz_project/models"
	"github.com/gofiber/fiber/v2"
)

func CartHandler(c *fiber.Ctx) error {
	// Lấy session
	sess, err := loginSignup.Store.Get(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Session error"})
	}

	// Lấy userID từ session
	userID, ok := sess.Get("user_id").(uint)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	// Lấy dữ liệu giỏ hàng
	var cartItems []models.CartItem
	if err := models.DB.Preload("Product").Where("user_id = ?", userID).Find(&cartItems).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Lỗi khi lấy giỏ hàng"})
	}

	// Xử lý dữ liệu giỏ hàng
	cartData := []fiber.Map{}
	var totalAmount float64 = 0

	for _, item := range cartItems {
		total := float64(item.Price) * float64(item.Quantity)
		totalAmount += total

		cartData = append(cartData, fiber.Map{
			"ID":          item.ID,
			"ProductID":   item.ProductID,
			"ProductName": item.Product.Name,
			"ImageURL":    item.Product.ImageURL,
			"Price":       float64(item.Price),
			"Quantity":    item.Quantity,
			"Total":       total,
		})
	}

	// Render ra giao diện cart.html và truyền thêm TotalAmount
	return c.Render("index/cart", fiber.Map{
		"CartItems":   cartData,
		"TotalAmount": totalAmount,
	}, "layouts/main")
}
