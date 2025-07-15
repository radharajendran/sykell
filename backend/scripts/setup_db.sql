-- Create database
CREATE DATABASE IF NOT EXISTS sykell_db;
USE sykell_db;

-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE,
    password VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- Create crawl_urls table
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
);

-- Create broken_links table
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
);

-- Insert sample data (optional)
INSERT INTO users (name, email, password) VALUES 
('John Doe', 'john@example.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi'), -- password: password
('Jane Smith', 'jane@example.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi'); -- password: password
