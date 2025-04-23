package controllers

import (
	"Zzz_project/models"
	"github.com/gofiber/fiber/v2"
	"log"
)

func AddToCart(c *fiber.Ctx) error {
	type CartRequest struct {
		ProductID uint  `json:"product_id"`
		Price     int64 `json:"price"`
		Quantity  int   `json:"quantity"`
	}

	var req CartRequest
	if err := c.BodyParser(&req); err != nil {
		log.Println("Lỗi parse JSON:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Dữ liệu không hợp lệ!"})
	}

	// Lấy ID người dùng từ session
	sess, err := Store.Get(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Lỗi session!"})
	}

	userID := sess.Get("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Bạn chưa đăng nhập!"})
	}

	// Kiểm tra sản phẩm có tồn tại không
	var product models.Product
	if err := models.DB.First(&product, req.ProductID).Error; err != nil {
		return c.JSON(fiber.Map{"success": false, "message": "Sản phẩm không tồn tại!"})
	}

	// Kiểm tra sản phẩm đã có trong giỏ hàng chưa
	var cartItem models.CartItem
	result := models.DB.Where("user_id = ? AND product_id = ?", userID, req.ProductID).First(&cartItem)

	// Tính tổng số lượng sau khi thêm
	newQuantity := req.Quantity
	if result.RowsAffected > 0 {
		newQuantity += cartItem.Quantity
	}

	// Kiểm tra nếu vượt quá tồn kho
	if newQuantity > product.Stock {
		return c.JSON(fiber.Map{"success": false, "message": "Không đủ hàng trong kho!", "maxQuantity": product.Stock})
	}

	// Cập nhật hoặc thêm mới vào giỏ hàng
	if result.RowsAffected > 0 {
		cartItem.Quantity = newQuantity
		if err := models.DB.Save(&cartItem).Error; err != nil {
			log.Println("Lỗi cập nhật giỏ hàng:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Không thể cập nhật giỏ hàng!"})
		}
	} else {
		cartItem = models.CartItem{
			UserID:    userID.(uint),
			ProductID: req.ProductID,
			Quantity:  req.Quantity,
			Price:     req.Price,
		}
		if err := models.DB.Create(&cartItem).Error; err != nil {
			log.Println("Lỗi thêm vào giỏ hàng:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Không thể thêm vào giỏ hàng!"})
		}
	}

	return c.JSON(fiber.Map{"success": true, "message": "Sản phẩm đã được thêm vào giỏ hàng!"})
}
