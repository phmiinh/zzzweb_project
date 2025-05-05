package adminInterface

import (
	"Zzz_project/models"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

func ProductsList(c *fiber.Ctx) error {
	var products []models.Product
	if err := models.DB.Find(&products).Error; err != nil {
		return c.Status(500).SendString("Database error")
	}

	return c.Render("admin/products", fiber.Map{
		"Products": products,
	}, "layouts/admin")
}

func ProductDetail(c *fiber.Ctx) error {
	id := c.Params("id")

	var product models.Product
	if err := models.DB.First(&product, id).Error; err != nil {
		return c.Status(404).SendString("Product not found")
	}

	return c.Render("admin/product_detail", fiber.Map{
		"Product": product,
	}, "layouts/admin")
}

// UpdateProduct handles POST /admin/products/edit/:id
func UpdateProduct(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid product ID")
	}

	var product models.Product
	if err := models.DB.First(&product, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).SendString("Product not found")
	}

	// Cập nhật dữ liệu từ form
	product.Name = c.FormValue("name")
	price, _ := strconv.ParseFloat(c.FormValue("price"), 64)
	oldPrice, _ := strconv.ParseFloat(c.FormValue("old_price"), 64)
	stock, _ := strconv.Atoi(c.FormValue("stock"))

	product.Price = price
	product.OldPrice = oldPrice
	product.Stock = stock
	product.Category = c.FormValue("category")
	product.ImageURL = c.FormValue("image_url")
	product.Description = c.FormValue("description")

	// Lưu lại database
	if err := models.DB.Save(&product).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to update product")
	}

	return c.Redirect("/admin/products")
}

// DeleteProduct handles POST /admin/products/delete/:id
func DeleteProduct(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid product ID")
	}

	if err := models.DB.Delete(&models.Product{}, id).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to delete product")
	}

	return c.Redirect("/admin/products")
}

func ShowAddProductForm(c *fiber.Ctx) error {
	return c.Render("admin/add_product", fiber.Map{}, "layouts/admin")
}

// AddProduct handles POST /admin/products/add
func AddProduct(c *fiber.Ctx) error {
	// Nhận form fields
	name := c.FormValue("name")
	priceStr := c.FormValue("price")
	oldPriceStr := c.FormValue("old_price")
	stockStr := c.FormValue("stock")
	category := c.FormValue("category")
	description := c.FormValue("description")

	price, _ := strconv.ParseFloat(priceStr, 64)
	oldPrice, _ := strconv.ParseFloat(oldPriceStr, 64)
	stock, _ := strconv.Atoi(stockStr)

	var imageURL string

	// Xử lý upload file ảnh
	file, err := c.FormFile("image")
	if err == nil && file != nil {
		uploadPath := "./public/uploads/products/" + file.Filename
		if err := c.SaveFile(file, uploadPath); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to save image")
		}
		imageURL = "/uploads/" + file.Filename
	}

	// Tạo product mới
	newProduct := models.Product{
		Name:        name,
		Price:       price,
		OldPrice:    oldPrice,
		Stock:       stock,
		Category:    category,
		ImageURL:    imageURL,
		Description: description,
	}

	if err := models.DB.Create(&newProduct).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to create product")
	}

	return c.Redirect("/admin/products")
}
