package handler

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"

	"sykell-backend/internal/models"
	"sykell-backend/internal/service"
	"sykell-backend/pkg/logger"
)

var crawlerService = service.NewCrawlerService()

// AddURL adds a new URL for crawling
func AddURL(c *fiber.Ctx) error {
	var req models.CrawlRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.URL == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "URL is required",
		})
	}

	crawlURL, err := crawlerService.AddURL(req.URL)
	if err != nil {
		logger.Sugar().Errorf("Failed to add URL: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"data":    crawlURL,
		"message": "URL added successfully",
	})
}

// StartCrawl starts crawling a URL by ID
func StartCrawl(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid URL ID",
		})
	}

	err = crawlerService.CrawlURL(id)
	if err != nil {
		logger.Sugar().Errorf("Failed to start crawl: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Crawl started successfully",
	})
}

// GetCrawlURLs returns paginated list of crawl URLs
func GetCrawlURLs(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	status := c.Query("status", "")
	search := c.Query("search", "")

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	crawlURLs, total, err := crawlerService.GetCrawlURLs(page, limit, status, search)
	if err != nil {
		logger.Sugar().Errorf("Failed to get crawl URLs: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch crawl URLs",
		})
	}

	return c.JSON(fiber.Map{
		"data": crawlURLs,
		"pagination": fiber.Map{
			"page":  page,
			"limit": limit,
			"total": total,
			"pages": (total + limit - 1) / limit,
		},
	})
}

// GetCrawlResult returns detailed crawl result including broken links
func GetCrawlResult(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid URL ID",
		})
	}

	result, err := crawlerService.GetCrawlResult(id)
	if err != nil {
		logger.Sugar().Errorf("Failed to get crawl result: %v", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Crawl result not found",
		})
	}

	return c.JSON(fiber.Map{
		"data": result,
	})
}

// DeleteCrawlURLs deletes multiple crawl URLs by IDs
func DeleteCrawlURLs(c *fiber.Ctx) error {
	var req struct {
		IDs []int `json:"ids"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if len(req.IDs) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No IDs provided",
		})
	}

	err := crawlerService.DeleteCrawlURLs(req.IDs)
	if err != nil {
		logger.Sugar().Errorf("Failed to delete crawl URLs: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete crawl URLs",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Crawl URLs deleted successfully",
	})
}

// ReCrawlURLs re-crawls multiple URLs by IDs
func ReCrawlURLs(c *fiber.Ctx) error {
	var req struct {
		IDs []int `json:"ids"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if len(req.IDs) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No IDs provided",
		})
	}

	err := crawlerService.ReCrawlURLs(req.IDs)
	if err != nil {
		logger.Sugar().Errorf("Failed to re-crawl URLs: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to re-crawl URLs",
		})
	}

	// Start crawling for each URL
	for _, id := range req.IDs {
		go func(crawlID int) {
			if err := crawlerService.CrawlURL(crawlID); err != nil {
				logger.Sugar().Errorf("Failed to start crawl for ID %d: %v", crawlID, err)
			}
		}(id)
	}

	return c.JSON(fiber.Map{
		"message": "Re-crawl started for selected URLs",
	})
}

// GetCrawlStats returns statistics about crawl URLs
func GetCrawlStats(c *fiber.Ctx) error {
	stats, err := crawlerService.GetStats()
	if err != nil {
		logger.Sugar().Errorf("Failed to get crawl stats: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch statistics",
		})
	}

	return c.JSON(fiber.Map{
		"data": stats,
	})
}

// BulkAddURLs adds multiple URLs for crawling
func BulkAddURLs(c *fiber.Ctx) error {
	var req models.BulkCrawlRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if len(req.URLs) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No URLs provided",
		})
	}

	var results []interface{}
	var errors []string

	for _, url := range req.URLs {
		url = strings.TrimSpace(url)
		if url == "" {
			continue
		}

		crawlURL, err := crawlerService.AddURL(url)
		if err != nil {
			errors = append(errors, fmt.Sprintf("Failed to add %s: %v", url, err))
			continue
		}
		results = append(results, crawlURL)
	}

	response := fiber.Map{
		"data":    results,
		"message": fmt.Sprintf("Added %d URLs successfully", len(results)),
	}

	if len(errors) > 0 {
		response["errors"] = errors
	}

	return c.Status(fiber.StatusCreated).JSON(response)
}
