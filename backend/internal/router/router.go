package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"sykell-backend/internal/handler"
	"sykell-backend/internal/middleware"
)

func Setup() *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	// CORS middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	// API routes
	api := app.Group("/api")

	// Public routes
	api.Post("/login", handler.Login)
	api.Post("/users", handler.CreateUser) // Allow public user registration

	// Protected routes
	protected := api.Group("/", middleware.JWTMiddleware())

	// User routes
	users := protected.Group("/users")
	users.Get("/", handler.GetUsers)
	users.Get("/:id", handler.GetUser)
	users.Put("/:id", handler.UpdateUser)
	users.Delete("/:id", handler.DeleteUser)

	// Crawler routes
	crawler := protected.Group("/crawler")
	crawler.Post("/urls", handler.AddURL)               // Add single URL
	crawler.Post("/urls/bulk", handler.BulkAddURLs)     // Add multiple URLs
	crawler.Get("/urls", handler.GetCrawlURLs)          // Get all crawl URLs with pagination/filtering
	crawler.Get("/urls/:id", handler.GetCrawlResult)    // Get detailed crawl result
	crawler.Post("/urls/:id/crawl", handler.StartCrawl) // Start crawling a URL
	crawler.Delete("/urls", handler.DeleteCrawlURLs)    // Delete multiple URLs
	crawler.Post("/urls/recrawl", handler.ReCrawlURLs)  // Re-crawl multiple URLs
	crawler.Get("/stats", handler.GetCrawlStats)        // Get crawl statistics

	return app
}
