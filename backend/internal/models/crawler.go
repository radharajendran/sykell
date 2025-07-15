package models

import (
	"time"
)

type CrawlURL struct {
	ID                     int        `json:"id" db:"id"`
	URL                    string     `json:"url" db:"url"`
	Status                 string     `json:"status" db:"status"`
	Title                  string     `json:"title" db:"title"`
	HTMLVersion            string     `json:"html_version" db:"html_version"`
	H1Count                int        `json:"h1_count" db:"h1_count"`
	H2Count                int        `json:"h2_count" db:"h2_count"`
	H3Count                int        `json:"h3_count" db:"h3_count"`
	H4Count                int        `json:"h4_count" db:"h4_count"`
	H5Count                int        `json:"h5_count" db:"h5_count"`
	H6Count                int        `json:"h6_count" db:"h6_count"`
	InternalLinksCount     int        `json:"internal_links_count" db:"internal_links_count"`
	ExternalLinksCount     int        `json:"external_links_count" db:"external_links_count"`
	InaccessibleLinksCount int        `json:"inaccessible_links_count" db:"inaccessible_links_count"`
	HasLoginForm           bool       `json:"has_login_form" db:"has_login_form"`
	ErrorMessage           string     `json:"error_message" db:"error_message"`
	LastCrawledAt          *time.Time `json:"last_crawled_at" db:"last_crawled_at"`
	CreatedAt              time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt              time.Time  `json:"updated_at" db:"updated_at"`
}

type BrokenLink struct {
	ID           int       `json:"id" db:"id"`
	CrawlURLID   int       `json:"crawl_url_id" db:"crawl_url_id"`
	URL          string    `json:"url" db:"url"`
	StatusCode   int       `json:"status_code" db:"status_code"`
	ErrorMessage string    `json:"error_message" db:"error_message"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

type CrawlResult struct {
	CrawlURL    CrawlURL     `json:"crawl_url"`
	BrokenLinks []BrokenLink `json:"broken_links"`
}

type CrawlRequest struct {
	URL string `json:"url" validate:"required,url"`
}

type BulkCrawlRequest struct {
	URLs []string `json:"urls" validate:"required,dive,url"`
}

type CrawlStats struct {
	TotalURLs     int `json:"total_urls"`
	QueuedURLs    int `json:"queued_urls"`
	RunningURLs   int `json:"running_urls"`
	CompletedURLs int `json:"completed_urls"`
	ErrorURLs     int `json:"error_urls"`
}

// Status constants
const (
	StatusQueued    = "queued"
	StatusRunning   = "running"
	StatusCompleted = "completed"
	StatusError     = "error"
)
