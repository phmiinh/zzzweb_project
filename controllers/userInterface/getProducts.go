package userInterface

import (
	"Zzz_project/models"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

func GetProducts(c *fiber.Ctx) error {
	var products []models.Product
	category := c.Query("category", "")
	search := c.Query("search", "") // Thêm query search
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize := 10

	db := models.DB.Model(&models.Product{})

	if category != "" && search != "" {
		db = db.Where("category = ? AND name LIKE ?", category, "%"+search+"%")
	} else if category != "" {
		db = db.Where("category = ?", category)
	} else if search != "" {
		db = db.Where("name LIKE ?", "%"+search+"%")
	}

	var total int64
	db.Count(&total)

	offset := (page - 1) * pageSize
	if err := db.Offset(offset).Limit(pageSize).Find(&products).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Lỗi khi lấy danh sách sản phẩm"})
	}

	return c.Render("index/shop_products", fiber.Map{
		"products":    products,
		"category":    category,
		"search":      search, // Truyền tham số tìm kiếm vào template
		"currentPage": page,
		"totalPages":  (total + int64(pageSize) - 1) / int64(pageSize),
	}, "layouts/main")
}

func GetSingleProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	var product models.Product

	if err := models.DB.First(&product, id).Error; err != nil {
		return c.Status(404).SendString("Product not found")
	}

	var relatedProducts []models.Product
	models.DB.Where("id != ?", product.ID).Limit(4).Find(&relatedProducts)

	return c.Render("index/single_product", fiber.Map{
		"Product": product,
	}, "layouts/main")
}
