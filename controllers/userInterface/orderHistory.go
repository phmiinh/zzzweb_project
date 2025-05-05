package userInterface

import (
	"Zzz_project/controllers/loginSignup"
	"Zzz_project/models"
	"github.com/gofiber/fiber/v2"
	"log"
)

// OrderHistory handler to show user's order history
func OrderHistory(c *fiber.Ctx) error {
	// Get the session and user_id
	sess, err := loginSignup.Store.Get(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Session error"})
	}

	userID, ok := sess.Get("user_id").(uint)
	if !ok {
		return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
	}

	// Fetch orders for the user
	var orders []models.Order
	if err := models.DB.Preload("OrderDetails").Where("user_id = ?", userID).Find(&orders).Error; err != nil {
		log.Println("❌ Error fetching orders:", err)
		return c.Status(500).JSON(fiber.Map{"error": "Unable to fetch orders"})
	}
	for i := range orders {
		var total int64
		for _, d := range orders[i].OrderDetails {
			total += d.Price * int64(d.Quantity)
		}
		orders[i].Total = total
	}

	return c.Render("index/order_history", fiber.Map{
		"Orders": orders,
	}, "layouts/main")
}

func OrderDetails(c *fiber.Ctx) error {
	orderID := c.Params("id")

	var order models.Order
	if err := models.DB.
		Preload("OrderDetails.Product"). // ✅ Thêm dòng này để preload Product
		Where("id = ?", orderID).
		First(&order).Error; err != nil {

		log.Println("❌ Error fetching order details:", err)
		return c.Status(500).JSON(fiber.Map{"error": "Unable to fetch order details"})
	}

	// Tính tổng tiền
	var total int64
	for _, detail := range order.OrderDetails {
		total += detail.Price * int64(detail.Quantity)
	}

	return c.Render("index/order_detail", fiber.Map{
		"Order": order,
		"Total": total,
	}, "layouts/main")
}
