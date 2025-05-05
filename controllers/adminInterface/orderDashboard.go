package adminInterface

import (
	"Zzz_project/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"log"
	"strconv"
	"strings"
)

func OrdersList(c *fiber.Ctx) error {
	var orders []models.Order

	if err := models.DB.Order("created_at desc").Find(&orders).Error; err != nil {
		return c.Status(500).SendString("Failed to fetch orders")
	}

	for i := range orders {
		var total int64
		models.DB.Model(&models.OrderDetail{}).
			Where("order_id = ?", orders[i].ID).
			Select("SUM(price * quantity)").Scan(&total)
		orders[i].Total = total
	}

	return c.Render("admin/orders", fiber.Map{
		"Orders": orders,
	}, "layouts/admin")
}

func OrderDetail(c *fiber.Ctx) error {
	orderID := c.Params("id")

	var order models.Order
	if err := models.DB.
		Preload("OrderDetails.Product").
		Where("id = ?", orderID).
		First(&order).Error; err != nil {
		log.Println("❌ Error fetching order details:", err)
		return c.Status(500).JSON(fiber.Map{"error": "Unable to fetch order details"})
	}

	var total int64
	for _, detail := range order.OrderDetails {
		total += detail.Price * int64(detail.Quantity)
	}
	order.Total = total

	return c.Render("admin/order_detail", fiber.Map{
		"Order":      order,
		"TotalFloat": float64(order.Total),
	}, "layouts/admin")
}

func UpdateOrder(c *fiber.Ctx) error {
	id := c.Params("id")
	var order models.Order

	if err := models.DB.First(&order, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "Order not found",
		})
	}

	// Cập nhật thông tin đơn hàng
	shippingAddress := c.FormValue("shipping_address")
	if shippingAddress != "" {
		parts := strings.SplitN(shippingAddress, ",", 2)
		order.Address = strings.TrimSpace(parts[0])
		if len(parts) > 1 {
			order.City = strings.TrimSpace(parts[1])
		}
	}
	order.Payment = c.FormValue("payment_method")
	order.Status = c.FormValue("status")

	if err := models.DB.Save(&order).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to update order info",
		})
	}

	// Xử lý cập nhật sản phẩm trong đơn
	formValues := c.Request().PostArgs()
	type ErrorResponse struct {
		HasError bool
		Message  string
	}
	var response ErrorResponse

	formValues.VisitAll(func(key, value []byte) {
		k := string(key)
		v := string(value)

		if strings.HasPrefix(k, "quantities[") {
			idStr := strings.TrimSuffix(strings.TrimPrefix(k, "quantities["), "]")
			orderDetailID, _ := strconv.ParseUint(idStr, 10, 64)
			newQuantity, _ := strconv.Atoi(v)

			var orderDetail models.OrderDetail
			if err := models.DB.Preload("Product").First(&orderDetail, orderDetailID).Error; err != nil {
				response = ErrorResponse{true, "Failed to load order detail"}
				return
			}

			oldQuantity := orderDetail.Quantity
			productID := orderDetail.ProductID
			productName := orderDetail.Product.Name
			productStock := orderDetail.Product.Stock

			delta := newQuantity - oldQuantity

			if newQuantity <= 0 {
				// Nếu quantity = 0 ➔ xóa sản phẩm và cộng lại stock
				models.DB.Model(&models.Product{}).Where("id = ?", productID).
					Update("stock", gorm.Expr("stock + ?", oldQuantity))
				models.DB.Delete(&models.OrderDetail{}, orderDetailID)
			} else {
				if delta > 0 {
					// Nếu tăng số lượng ➔ check stock
					if productStock < delta {
						response = ErrorResponse{true, "Not enough stock for product: " + productName}
						return
					}
					// Giảm stock
					models.DB.Model(&models.Product{}).Where("id = ?", productID).
						Update("stock", gorm.Expr("stock - ?", delta))
				}

				if delta < 0 {
					// Nếu giảm quantity ➔ cộng lại stock
					models.DB.Model(&models.Product{}).Where("id = ?", productID).
						Update("stock", gorm.Expr("stock + ?", -delta))
				}

				// Cập nhật quantity mới
				models.DB.Model(&models.OrderDetail{}).Where("id = ?", orderDetailID).
					Update("quantity", newQuantity)
			}
		}
	})

	// Nếu gặp lỗi trong quá trình cập nhật thì trả lỗi
	if response.HasError {
		return c.Status(400).JSON(fiber.Map{
			"error": response.Message,
		})
	}

	// Check nếu đơn hàng không còn sản phẩm ➔ xóa đơn luôn
	var remaining int64
	models.DB.Model(&models.OrderDetail{}).Where("order_id = ?", order.ID).Count(&remaining)
	if remaining == 0 {
		models.DB.Delete(&order)
	}

	return c.SendStatus(200)
}
func DeleteOrder(c *fiber.Ctx) error {
	id := c.Params("id")

	// Xoá tất cả chi tiết đơn hàng trước
	if err := models.DB.Where("order_id = ?", id).Delete(&models.OrderDetail{}).Error; err != nil {
		return c.Status(500).SendString("Failed to delete order details")
	}

	// Xoá luôn đơn hàng
	if err := models.DB.Delete(&models.Order{}, id).Error; err != nil {
		return c.Status(500).SendString("Failed to delete order")
	}

	return c.SendStatus(200)
}
