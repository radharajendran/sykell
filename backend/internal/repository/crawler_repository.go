package repository

import (
	"database/sql"
	"fmt"
	"strings"

	"sykell-backend/internal/models"
	"sykell-backend/pkg/database"
	"sykell-backend/pkg/logger"
)

type CrawlerRepository struct {
	db *sql.DB
}

func NewCrawlerRepository() *CrawlerRepository {
	return &CrawlerRepository{
		db: database.DB,
	}
}

func (r *CrawlerRepository) CreateCrawlURL(url string) (*models.CrawlURL, error) {
	query := `
		INSERT INTO crawl_urls (url, status) 
		VALUES (?, ?)
		ON DUPLICATE KEY UPDATE 
			status = VALUES(status),
			updated_at = CURRENT_TIMESTAMP
	`

	result, err := r.db.Exec(query, url, models.StatusQueued)
	if err != nil {
		logger.Sugar().Errorf("Failed to create crawl URL: %v", err)
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		logger.Sugar().Errorf("Failed to get last insert ID: %v", err)
		return nil, err
	}

	return r.GetCrawlURLByID(int(id))
}

func (r *CrawlerRepository) GetCrawlURLByID(id int) (*models.CrawlURL, error) {
	query := `
		SELECT id, url, status, title, html_version, 
			   h1_count, h2_count, h3_count, h4_count, h5_count, h6_count,
			   internal_links_count, external_links_count, inaccessible_links_count,
			   has_login_form, error_message, last_crawled_at, created_at, updated_at
		FROM crawl_urls WHERE id = ?
	`

	var crawlURL models.CrawlURL
	var title, htmlVersion, errorMessage sql.NullString
	var lastCrawledAt sql.NullTime

	err := r.db.QueryRow(query, id).Scan(
		&crawlURL.ID, &crawlURL.URL, &crawlURL.Status, &title, &htmlVersion,
		&crawlURL.H1Count, &crawlURL.H2Count, &crawlURL.H3Count, &crawlURL.H4Count,
		&crawlURL.H5Count, &crawlURL.H6Count, &crawlURL.InternalLinksCount,
		&crawlURL.ExternalLinksCount, &crawlURL.InaccessibleLinksCount,
		&crawlURL.HasLoginForm, &errorMessage, &lastCrawledAt,
		&crawlURL.CreatedAt, &crawlURL.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		logger.Sugar().Errorf("Failed to get crawl URL by ID: %v", err)
		return nil, err
	}

	if title.Valid {
		crawlURL.Title = title.String
	}
	if htmlVersion.Valid {
		crawlURL.HTMLVersion = htmlVersion.String
	}
	if errorMessage.Valid {
		crawlURL.ErrorMessage = errorMessage.String
	}
	if lastCrawledAt.Valid {
		crawlURL.LastCrawledAt = &lastCrawledAt.Time
	}

	return &crawlURL, nil
}

func (r *CrawlerRepository) GetCrawlURLs(limit, offset int, status, search string) ([]models.CrawlURL, int, error) {
	var whereClause []string
	var args []interface{}

	if status != "" {
		whereClause = append(whereClause, "status = ?")
		args = append(args, status)
	}

	if search != "" {
		whereClause = append(whereClause, "(url LIKE ? OR title LIKE ?)")
		searchPattern := "%" + search + "%"
		args = append(args, searchPattern, searchPattern)
	}

	whereSQL := ""
	if len(whereClause) > 0 {
		whereSQL = "WHERE " + strings.Join(whereClause, " AND ")
	}

	// Count total records
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM crawl_urls %s", whereSQL)
	var total int
	err := r.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		logger.Sugar().Errorf("Failed to count crawl URLs: %v", err)
		return nil, 0, err
	}

	// Get paginated records
	query := fmt.Sprintf(`
		SELECT id, url, status, title, html_version, 
			   h1_count, h2_count, h3_count, h4_count, h5_count, h6_count,
			   internal_links_count, external_links_count, inaccessible_links_count,
			   has_login_form, error_message, last_crawled_at, created_at, updated_at
		FROM crawl_urls %s
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`, whereSQL)

	args = append(args, limit, offset)
	rows, err := r.db.Query(query, args...)
	if err != nil {
		logger.Sugar().Errorf("Failed to get crawl URLs: %v", err)
		return nil, 0, err
	}
	defer rows.Close()

	var crawlURLs []models.CrawlURL
	for rows.Next() {
		var crawlURL models.CrawlURL
		var title, htmlVersion, errorMessage sql.NullString
		var lastCrawledAt sql.NullTime

		err := rows.Scan(
			&crawlURL.ID, &crawlURL.URL, &crawlURL.Status, &title, &htmlVersion,
			&crawlURL.H1Count, &crawlURL.H2Count, &crawlURL.H3Count, &crawlURL.H4Count,
			&crawlURL.H5Count, &crawlURL.H6Count, &crawlURL.InternalLinksCount,
			&crawlURL.ExternalLinksCount, &crawlURL.InaccessibleLinksCount,
			&crawlURL.HasLoginForm, &errorMessage, &lastCrawledAt,
			&crawlURL.CreatedAt, &crawlURL.UpdatedAt,
		)

		if err != nil {
			logger.Sugar().Errorf("Failed to scan crawl URL: %v", err)
			return nil, 0, err
		}

		if title.Valid {
			crawlURL.Title = title.String
		}
		if htmlVersion.Valid {
			crawlURL.HTMLVersion = htmlVersion.String
		}
		if errorMessage.Valid {
			crawlURL.ErrorMessage = errorMessage.String
		}
		if lastCrawledAt.Valid {
			crawlURL.LastCrawledAt = &lastCrawledAt.Time
		}

		crawlURLs = append(crawlURLs, crawlURL)
	}

	return crawlURLs, total, nil
}

func (r *CrawlerRepository) UpdateCrawlURL(crawlURL *models.CrawlURL) error {
	query := `
		UPDATE crawl_urls SET 
			status = ?, title = ?, html_version = ?,
			h1_count = ?, h2_count = ?, h3_count = ?, h4_count = ?, h5_count = ?, h6_count = ?,
			internal_links_count = ?, external_links_count = ?, inaccessible_links_count = ?,
			has_login_form = ?, error_message = ?, last_crawled_at = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	_, err := r.db.Exec(query,
		crawlURL.Status, crawlURL.Title, crawlURL.HTMLVersion,
		crawlURL.H1Count, crawlURL.H2Count, crawlURL.H3Count, crawlURL.H4Count,
		crawlURL.H5Count, crawlURL.H6Count, crawlURL.InternalLinksCount,
		crawlURL.ExternalLinksCount, crawlURL.InaccessibleLinksCount,
		crawlURL.HasLoginForm, crawlURL.ErrorMessage, crawlURL.LastCrawledAt,
		crawlURL.ID,
	)

	if err != nil {
		logger.Sugar().Errorf("Failed to update crawl URL: %v", err)
		return err
	}

	return nil
}

func (r *CrawlerRepository) DeleteCrawlURLs(ids []int) error {
	if len(ids) == 0 {
		return nil
	}

	placeholders := strings.Repeat("?,", len(ids)-1) + "?"
	query := fmt.Sprintf("DELETE FROM crawl_urls WHERE id IN (%s)", placeholders)

	args := make([]interface{}, len(ids))
	for i, id := range ids {
		args[i] = id
	}

	_, err := r.db.Exec(query, args...)
	if err != nil {
		logger.Sugar().Errorf("Failed to delete crawl URLs: %v", err)
		return err
	}

	return nil
}

func (r *CrawlerRepository) GetBrokenLinks(crawlURLID int) ([]models.BrokenLink, error) {
	query := `
		SELECT id, crawl_url_id, url, status_code, error_message, created_at
		FROM broken_links 
		WHERE crawl_url_id = ?
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, crawlURLID)
	if err != nil {
		logger.Sugar().Errorf("Failed to get broken links: %v", err)
		return nil, err
	}
	defer rows.Close()

	var brokenLinks []models.BrokenLink
	for rows.Next() {
		var brokenLink models.BrokenLink
		var statusCode sql.NullInt64
		var errorMessage sql.NullString

		err := rows.Scan(
			&brokenLink.ID, &brokenLink.CrawlURLID, &brokenLink.URL,
			&statusCode, &errorMessage, &brokenLink.CreatedAt,
		)

		if err != nil {
			logger.Sugar().Errorf("Failed to scan broken link: %v", err)
			return nil, err
		}

		if statusCode.Valid {
			brokenLink.StatusCode = int(statusCode.Int64)
		}
		if errorMessage.Valid {
			brokenLink.ErrorMessage = errorMessage.String
		}

		brokenLinks = append(brokenLinks, brokenLink)
	}

	return brokenLinks, nil
}

func (r *CrawlerRepository) CreateBrokenLink(brokenLink *models.BrokenLink) error {
	query := `
		INSERT INTO broken_links (crawl_url_id, url, status_code, error_message)
		VALUES (?, ?, ?, ?)
	`

	_, err := r.db.Exec(query,
		brokenLink.CrawlURLID, brokenLink.URL,
		brokenLink.StatusCode, brokenLink.ErrorMessage,
	)

	if err != nil {
		logger.Sugar().Errorf("Failed to create broken link: %v", err)
		return err
	}

	return nil
}

func (r *CrawlerRepository) GetCrawlStats() (*models.CrawlStats, error) {
	query := `
		SELECT 
			COUNT(*) as total,
			SUM(CASE WHEN status = 'queued' THEN 1 ELSE 0 END) as queued,
			SUM(CASE WHEN status = 'running' THEN 1 ELSE 0 END) as running,
			SUM(CASE WHEN status = 'completed' THEN 1 ELSE 0 END) as completed,
			SUM(CASE WHEN status = 'error' THEN 1 ELSE 0 END) as error
		FROM crawl_urls
	`

	var stats models.CrawlStats
	err := r.db.QueryRow(query).Scan(
		&stats.TotalURLs, &stats.QueuedURLs, &stats.RunningURLs,
		&stats.CompletedURLs, &stats.ErrorURLs,
	)

	if err != nil {
		logger.Sugar().Errorf("Failed to get crawl stats: %v", err)
		return nil, err
	}

	return &stats, nil
}
