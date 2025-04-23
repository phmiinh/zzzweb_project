package controllers

import (
	"Zzz_project/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func CheckoutPage(c *fiber.Ctx) error {
	sess, err := Store.Get(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Session error"})
	}

	userID, ok := sess.Get("user_id").(uint)
	if !ok {
		return c.Redirect("/login")
	}

	var cartItems []models.CartItem
	if err := models.DB.Preload("Product").Where("user_id = ?", userID).Find(&cartItems).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Không thể lấy dữ liệu giỏ hàng"})
	}

	var totalAmount float64
	cartData := []fiber.Map{}

	for _, item := range cartItems {
		total := float64(item.Price) * float64(item.Quantity)
		totalAmount += total
		cartData = append(cartData, fiber.Map{
			"ProductID":   item.ProductID,
			"ProductName": item.Product.Name,
			"ImageURL":    item.Product.ImageURL,
			"Price":       item.Price,
			"Quantity":    item.Quantity,
			"Total":       total,
		})
	}

	return c.Render("index/checkout", fiber.Map{
		"CartItems":   cartData,
		"TotalAmount": totalAmount,
	}, "layouts/main")
}
func ProcessCheckout(c *fiber.Ctx) error {
	sess, err := Store.Get(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Session error"})
	}

	userID, ok := sess.Get("user_id").(uint)
	if !ok {
		return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
	}

	form := new(struct {
		FirstName  string `json:"billing_first_name"`
		LastName   string `json:"billing_last_name"`
		Address    string `json:"billing_address_1"`
		City       string `json:"billing_city"`
		Email      string `json:"billing_email"`
		Phone      string `json:"billing_phone"`
		Payment    string `json:"payment_method"`
		OrderNotes string `json:"order_comments"`
	})
	if err := c.BodyParser(form); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid form data"})
	}

	var cartItems []models.CartItem
	if err := models.DB.Preload("Product").Where("user_id = ?", userID).Find(&cartItems).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Lỗi khi lấy giỏ hàng"})
	}

	if len(cartItems) == 0 {
		return c.JSON(fiber.Map{"success": false, "message": "Giỏ hàng trống!"})
	}

	tx := models.DB.Begin()

	order := models.Order{
		UserID:     userID,
		FirstName:  form.FirstName,
		LastName:   form.LastName,
		Address:    form.Address,
		City:       form.City,
		Email:      form.Email,
		Phone:      form.Phone,
		Payment:    form.Payment,
		OrderNotes: form.OrderNotes,
		Status:     "Pending",
	}
	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{"error": "Không thể tạo đơn hàng"})
	}

	for _, item := range cartItems {
		result := tx.Model(&models.Product{}).
			Where("id = ? AND stock >= ?", item.ProductID, item.Quantity).
			Update("stock", gorm.Expr("stock - ?", item.Quantity))

		if result.Error != nil {
			tx.Rollback()
			return c.Status(500).JSON(fiber.Map{"error": "Lỗi khi cập nhật tồn kho"})
		}
		if result.RowsAffected == 0 {
			tx.Rollback()
			return c.JSON(fiber.Map{"success": false, "message": "Sản phẩm " + item.Product.Name + " không đủ hàng!"})
		}

		orderDetail := models.OrderDetail{
			OrderID:   order.ID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     item.Price,
		}
		if err := tx.Create(&orderDetail).Error; err != nil {
			tx.Rollback()
			return c.Status(500).JSON(fiber.Map{"error": "Không thể tạo đơn hàng"})
		}
	}

	if err := tx.Where("user_id = ?", userID).Delete(&models.CartItem{}).Error; err != nil {
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{"error": "Không thể xóa giỏ hàng"})
	}
	tx.Commit()

	return c.JSON(fiber.Map{"success": true, "message": "Đơn hàng đã được đặt thành công!"})
}
