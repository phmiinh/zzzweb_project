package controllers

import (
	"Zzz_project/models"
	"github.com/gofiber/fiber/v2"
)

// Hiển thị danh sách đơn hàng
func AdminOrdersHandler(c *fiber.Ctx) error {
	var orders []models.Order
	if err := models.DB.Preload("OrderDetails").Find(&orders).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Lỗi truy vấn đơn hàng")
	}

	return c.Render("admin_orders", fiber.Map{
		"Title":  "Orders Management",
		"Orders": orders,
	}, "layouts/admin")
}

// Hiển thị chi tiết đơn hàng
func AdminOrderDetailHandler(c *fiber.Ctx) error {
	orderID := c.Params("id")

	var order models.Order
	if err := models.DB.Preload("OrderDetails").First(&order, orderID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).SendString("Đơn hàng không tồn tại")
	}

	return c.Render("admin/admin_order_detail", fiber.Map{
		"Title":        "Order Details",
		"Order":        order,
		"OrderDetails": order.OrderDetails,
	}, "layouts/admin")
}
