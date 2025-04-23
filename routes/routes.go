package routes

import (
	"Zzz_project/controllers"
	"github.com/gofiber/fiber/v2"
	"strings"
)

func SetupRoutes(app *fiber.App) {

	authMiddleware := func(c *fiber.Ctx) error {
		sess, err := controllers.Store.Get(c)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Session error")
		}

		userID := sess.Get("user_id")
		if userID == nil {
			return c.Redirect("/auth/login")
		}

		role, _ := sess.Get("role").(string)
		if role == "admin" && !strings.HasPrefix(c.Path(), "/admin") {
			return c.Redirect("/admin/dashboard")
		}

		return c.Next()
	}

	index := func(c *fiber.Ctx) error {
		return c.Render("index/home", fiber.Map{
			"Title": "Home",
		}, "layouts/main")
	}

	//auth := app.Group("/auth")
	//admin := app.Group("/admin")
	//order := admin.Group("/order")
	//
	//order.Get("", index)
	//order.Post("", index)
	//order.Get("/:id", index)
	//order.Put("/:id", index)
	//order.Delete("/:id", index)

	app.Get("/auth/login", controllers.GetLogin)

	app.Post("/auth/login", controllers.LoginHandler)

	app.Get("/auth/signup", func(c *fiber.Ctx) error {
		return c.Render("shop/login_signup", nil)
	})

	app.Post("/auth/signup", controllers.RegisterHandler)

	app.Get("/auth/home", authMiddleware, index)

	app.Get("auth/logout", controllers.LogoutHandler)

	app.Get("/auth/account", controllers.AccountHandler)

	app.Post("/auth/account/update", controllers.UpdateAccountHandler)

	app.Get("/auth/change-password", controllers.ChangePasswordViewHandler)

	app.Post("/auth/change-password", controllers.ChangePasswordHandler)

	app.Get("/auth/shop_products", authMiddleware, controllers.GetProducts)

	app.Get("/auth/product/:id", authMiddleware, controllers.GetSingleProduct)

	app.Post("/auth/add-to-cart", authMiddleware, controllers.AddToCart)

	app.Get("/auth/cart", authMiddleware, controllers.CartHandler)

	app.Post("/cart/action", authMiddleware, controllers.CartFunc)

	app.Get("/auth/checkout", authMiddleware, controllers.CheckoutPage)

	app.Post("/auth/checkout", authMiddleware, controllers.ProcessCheckout)

	app.Get("/admin/dashboard", authMiddleware, controllers.AdminDashboard)

	app.Get("/admin/users", authMiddleware, controllers.AdminUsersHandler)

	app.Post("/admin/users/update-role", authMiddleware, controllers.UpdateUserRole)

	app.Post("/admin/users/delete", authMiddleware, controllers.DeleteUser)

	app.Get("/admin/orders", authMiddleware, controllers.AdminOrdersHandler)

	app.Get("/admin/orders/:id", authMiddleware, controllers.AdminOrderDetailHandler)

	app.Get("/auth/order-history", authMiddleware, controllers.OrderHistory)

	app.Get("/auth/order/:id", authMiddleware, controllers.OrderDetails)

}
