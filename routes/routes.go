package routes

import (
	"Zzz_project/controllers/adminInterface"
	"Zzz_project/controllers/loginSignup"
	"Zzz_project/controllers/userInterface"
	"github.com/gofiber/fiber/v2"
	"strings"
)

func SetupRoutes(app *fiber.App) {
	auth := app.Group("/auth")
	admin := app.Group("/admin")

	authMiddleware := func(c *fiber.Ctx) error {
		sess, err := loginSignup.Store.Get(c)
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

	auth.Get("/login", loginSignup.GetLogin)

	auth.Post("/login", loginSignup.LoginHandler)

	auth.Get("/signup", func(c *fiber.Ctx) error {
		return c.Render("shop/login_signup", nil)
	})

	auth.Post("/signup", loginSignup.RegisterHandler)

	auth.Get("/home", authMiddleware, userInterface.HomePage)

	auth.Get("/logout", loginSignup.LogoutHandler)

	auth.Get("/account", userInterface.AccountHandler)

	auth.Post("/account/update", userInterface.UpdateAccountHandler)

	auth.Get("/change-password", userInterface.ChangePasswordViewHandler)

	auth.Post("/change-password", userInterface.ChangePasswordHandler)

	auth.Get("/shop_products", authMiddleware, userInterface.GetProducts)

	auth.Get("/product/:id", authMiddleware, userInterface.GetSingleProduct)

	auth.Post("/add-to-cart", authMiddleware, userInterface.AddToCart)

	auth.Get("/cart", authMiddleware, userInterface.CartHandler)

	auth.Post("/cart/action", authMiddleware, userInterface.CartFunc)

	auth.Get("/checkout", authMiddleware, userInterface.CheckoutPage)

	auth.Post("/checkout", authMiddleware, userInterface.ProcessCheckout)

	auth.Get("/order-history", authMiddleware, userInterface.OrderHistory)

	auth.Get("/order/:id", authMiddleware, userInterface.OrderDetails)

	admin.Get("/dashboard", authMiddleware, adminInterface.AdminDashboard)

	admin.Get("/users", authMiddleware, adminInterface.UsersList)

	admin.Get("/users/:id", authMiddleware, adminInterface.UserDetail)

	admin.Post("/users/:id/change-role", authMiddleware, adminInterface.ChangeUserRole)

	admin.Post("/users/:id/delete", authMiddleware, adminInterface.DeleteUser)

	admin.Post("/users/edit/:id", authMiddleware, adminInterface.UpdateUser)

	admin.Get("/products", authMiddleware, adminInterface.ProductsList)

	admin.Get("/products/:id", authMiddleware, adminInterface.ProductDetail)

	admin.Post("/products/edit/:id", authMiddleware, adminInterface.UpdateProduct)

	admin.Post("/products/delete/:id", authMiddleware, adminInterface.DeleteProduct)

	admin.Get("/product/add", authMiddleware, adminInterface.ShowAddProductForm)

	admin.Post("/product/add", authMiddleware, adminInterface.AddProduct)

	admin.Get("/orders", authMiddleware, adminInterface.OrdersList)

	admin.Get("/orders/:id", authMiddleware, adminInterface.OrderDetail)

	admin.Post("/orders/edit/:id", authMiddleware, adminInterface.UpdateOrder)

	admin.Delete("/orders/delete/:id", authMiddleware, adminInterface.DeleteOrder)

	admin.Get("/banners", authMiddleware, adminInterface.ListBanners)

	admin.Post("/banners/create", authMiddleware, adminInterface.CreateBanner)

	admin.Post("/banners/update/:id", authMiddleware, adminInterface.UpdateBanner)

	admin.Delete("/banners/delete/:id", authMiddleware, adminInterface.DeleteBanner)

}
