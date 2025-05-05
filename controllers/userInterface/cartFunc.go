package userInterface

import (
	"Zzz_project/controllers/loginSignup"
	"Zzz_project/models"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

func CartFunc(c *fiber.Ctx) error {
	var DB = models.DB
	// Lấy session
	sess, err := loginSignup.Store.Get(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Session error",
		})
	}

	userID := sess.Get("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error":   "User not logged in",
		})
	}
	uid := userID.(uint)

	// Lấy action và dữ liệu form
	action := c.FormValue("action")
	productID := c.FormValue("product_id")
	quantityStr := c.FormValue("quantity")

	// Kiểm tra productID hợp lệ
	pid, err := strconv.Atoi(productID)
	if err != nil || pid <= 0 {
		return c.JSON(fiber.Map{
			"success": false,
			"error":   "Invalid product ID",
		})
	}

	switch action {
	case "remove":
		// Xóa sản phẩm khỏi giỏ
		if err := DB.Where("user_id = ? AND product_id = ?", uid, pid).
			Delete(&models.CartItem{}).Error; err != nil {
			return c.JSON(fiber.Map{
				"success": false,
				"error":   "Failed to remove item",
			})
		}
		return c.JSON(fiber.Map{"success": true})

	case "update":
		// Chuyển quantity về số
		qty, err := strconv.Atoi(quantityStr)
		if err != nil || qty < 1 {
			qty = 1
		}

		// Lấy thông tin sản phẩm từ DB
		var product models.Product
		if err := DB.First(&product, pid).Error; err != nil {
			return c.JSON(fiber.Map{
				"success": false,
				"error":   "Product not found",
			})
		}

		if qty > product.Stock {
			// Trả về lỗi nếu quá tồn kho
			return c.JSON(fiber.Map{
				"success":     false,
				"maxQuantity": product.Stock,
				"error":       "Quantity exceeds stock",
			})
		}

		// Tìm item trong giỏ
		var cartItem models.CartItem
		if err := DB.Where("user_id = ? AND product_id = ?", uid, pid).
			First(&cartItem).Error; err != nil {
			return c.JSON(fiber.Map{
				"success": false,
				"error":   "Cart item not found",
			})
		}

		// Cập nhật số lượng và giá
		cartItem.Quantity = qty
		cartItem.Price = int64(float64(qty) * product.Price)

		if err := DB.Save(&cartItem).Error; err != nil {
			return c.JSON(fiber.Map{
				"success": false,
				"error":   "Failed to update cart",
			})
		}

		return c.JSON(fiber.Map{
			"success": true,
		})

	default:
		return c.JSON(fiber.Map{
			"success": false,
			"error":   "Invalid action",
		})
	}
}
