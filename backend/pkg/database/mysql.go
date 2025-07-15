package database

import (
	"database/sql"
	"fmt"
	"os"

	"sykell-backend/pkg/logger"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

func Init() error {
	config := Config{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "3306"),
		User:     getEnv("DB_USER", "root"),
		Password: getEnv("DB_PASSWORD", "password"),
		DBName:   getEnv("DB_NAME", "sykell_db"),
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.User, config.Password, config.Host, config.Port, config.DBName)

	var err error
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		logger.Sugar().Errorf("Failed to connect to database: %v", err)
		return err
	}

	if err = DB.Ping(); err != nil {
		logger.Sugar().Errorf("Failed to ping database: %v", err)
		return err
	}

	logger.Sugar().Info("Successfully connected to MySQL database")

	// Create tables if they don't exist
	if err := createTables(); err != nil {
		logger.Sugar().Errorf("Failed to create tables: %v", err)
		return err
	}

	return nil
}

func createTables() error {
	// Users table
	usersQuery := `
	CREATE TABLE IF NOT EXISTS users (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		email VARCHAR(255) UNIQUE,
		password VARCHAR(255),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
	);`

	_, err := DB.Exec(usersQuery)
	if err != nil {
		logger.Sugar().Errorf("Failed to create users table: %v", err)
		return err
	}

	// Crawl URLs table
	crawlUrlsQuery := `
	CREATE TABLE IF NOT EXISTS crawl_urls (
		id INT AUTO_INCREMENT PRIMARY KEY,
		url VARCHAR(2048) NOT NULL UNIQUE,
		status ENUM('queued', 'running', 'completed', 'error') DEFAULT 'queued',
		title VARCHAR(512),
		html_version VARCHAR(50),
		h1_count INT DEFAULT 0,
		h2_count INT DEFAULT 0,
		h3_count INT DEFAULT 0,
		h4_count INT DEFAULT 0,
		h5_count INT DEFAULT 0,
		h6_count INT DEFAULT 0,
		internal_links_count INT DEFAULT 0,
		external_links_count INT DEFAULT 0,
		inaccessible_links_count INT DEFAULT 0,
		has_login_form BOOLEAN DEFAULT FALSE,
		error_message TEXT,
		last_crawled_at TIMESTAMP NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		INDEX idx_status (status),
		INDEX idx_url (url),
		INDEX idx_created_at (created_at)
	);`

	_, err = DB.Exec(crawlUrlsQuery)
	if err != nil {
		logger.Sugar().Errorf("Failed to create crawl_urls table: %v", err)
		return err
	}

	// Broken links table
	brokenLinksQuery := `
	CREATE TABLE IF NOT EXISTS broken_links (
		id INT AUTO_INCREMENT PRIMARY KEY,
		crawl_url_id INT NOT NULL,
		url VARCHAR(2048) NOT NULL,
		status_code INT,
		error_message TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (crawl_url_id) REFERENCES crawl_urls(id) ON DELETE CASCADE,
		INDEX idx_crawl_url_id (crawl_url_id),
		INDEX idx_status_code (status_code)
	);`

	_, err = DB.Exec(brokenLinksQuery)
	if err != nil {
		logger.Sugar().Errorf("Failed to create broken_links table: %v", err)
		return err
	}

	logger.Sugar().Info("Database tables created successfully")
	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}
