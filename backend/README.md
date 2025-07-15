# Sykell Backend API

A Go REST API using Fiber framework with MySQL database integration and web crawling capabilities.

## Features

- **User Management**: CRUD operations with JWT authentication
- **Web Crawler**: Comprehensive website analysis tool
  - HTML version detection
  - Page title extraction
  - Heading tag counting (H1-H6)
  - Internal vs external link categorization
  - Broken link detection (4xx/5xx status codes)
  - Login form detection
  - Real-time crawl status tracking
- **Database Integration**: MySQL with proper schema and indexing
- **Security**: Password hashing with bcrypt, JWT tokens
- **Background Processing**: Automatic job queue processing
- **RESTful API**: Clean, consistent endpoints
- **Docker Support**: Full containerization

## Prerequisites

- Go 1.21 or higher
- MySQL 5.7 or higher (or Docker)
- Git

## Quick Start with Docker

1. **Clone and start with Docker Compose**
   ```bash
   git clone <repository-url>
   cd sykell/backend
   docker-compose up -d
   ```

The API will be available at `http://localhost:8080`

## Manual Setup

1. **Install dependencies**
   ```bash
   make deps
   # or
   go mod tidy
   ```

2. **Set up MySQL database**
   ```bash
   make setup-db
   # or manually:
   mysql -u root -p < scripts/setup_db.sql
   ```

3. **Configure environment variables**
   ```bash
   cp .env.example .env
   # Edit .env with your database credentials
   ```

4. **Build and run**
   ```bash
   make build && ./bin/sykell-api
   # or for development:
   make run
   ```

## Available Make Commands

```bash
make build      # Build the application
make run        # Run the application
make dev        # Run with hot reload (requires air)
make test       # Run tests
make clean      # Clean build artifacts
make setup-db   # Setup database
make deps       # Install dependencies
make fmt        # Format code
make lint       # Run linter
make test-api   # Test API endpoints
```

## API Endpoints

### Authentication
- `POST /api/login` - User login
- `POST /api/users` - Create new user (public)

### Users (Protected)
- `GET /api/users` - Get all users
- `GET /api/users/:id` - Get user by ID
- `PUT /api/users/:id` - Update user
- `DELETE /api/users/:id` - Delete user

### Web Crawler (Protected)
- `POST /api/crawler/urls` - Add single URL for crawling
- `POST /api/crawler/urls/bulk` - Add multiple URLs for crawling
- `GET /api/crawler/urls` - Get all crawl URLs (paginated, filterable)
- `GET /api/crawler/urls/:id` - Get detailed crawl result
- `POST /api/crawler/urls/:id/crawl` - Start crawling a specific URL
- `DELETE /api/crawler/urls` - Delete multiple crawl URLs
- `POST /api/crawler/urls/recrawl` - Re-crawl multiple URLs
- `GET /api/crawler/stats` - Get crawl statistics

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| DB_HOST | MySQL host | localhost |
| DB_PORT | MySQL port | 3306 |
| DB_USER | MySQL username | root |
| DB_PASSWORD | MySQL password | password |
| DB_NAME | Database name | sykell_db |
| JWT_SECRET | JWT signing secret | your-secret-key |
| PORT | Server port | 8080 |

## Testing the Web Crawler

Run the crawler test script (requires server to be running):
```bash
make test-api
# or
./scripts/test_crawler_api.sh
```

### Manual Testing Examples

#### Add URL for Crawling
```bash
curl -X POST http://localhost:8080/api/crawler/urls \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{"url":"https://example.com"}'
```

#### Start Crawling
```bash
curl -X POST http://localhost:8080/api/crawler/urls/1/crawl \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

#### Get Crawl Results
```bash
curl -X GET http://localhost:8080/api/crawler/urls/1 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

#### Bulk Add URLs
```bash
curl -X POST http://localhost:8080/api/crawler/urls/bulk \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{"urls":["https://example.com","https://httpbin.org","https://github.com"]}'
```

## Database Schema

### Users Table
```sql
CREATE TABLE users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE,
    password VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```

### Crawl URLs Table
```sql
CREATE TABLE crawl_urls (
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
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```

### Broken Links Table
```sql
CREATE TABLE broken_links (
    id INT AUTO_INCREMENT PRIMARY KEY,
    crawl_url_id INT NOT NULL,
    url VARCHAR(2048) NOT NULL,
    status_code INT,
    error_message TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (crawl_url_id) REFERENCES crawl_urls(id) ON DELETE CASCADE
);
```

## Technology Stack

- **Framework**: Fiber (Express-inspired web framework)
- **Database**: MySQL with native Go driver
- **Authentication**: JWT tokens
- **Logging**: Zap (structured logging)
- **Password Hashing**: bcrypt
- **HTML Parsing**: golang.org/x/net/html
- **HTTP Client**: Native Go net/http with custom timeouts
- **Background Jobs**: Go routines with graceful shutdown
- **Containerization**: Docker & Docker Compose

## Web Crawler Features

### Data Collection per URL
- **HTML Version**: Automatically detects HTML version (defaults to HTML5)
- **Page Title**: Extracts the `<title>` tag content
- **Heading Counts**: Counts H1-H6 tags for SEO analysis
- **Link Analysis**: 
  - Categorizes internal vs external links
  - Checks link accessibility (HEAD requests)
  - Identifies broken links with status codes
- **Login Form Detection**: Identifies forms with password fields
- **Real-time Status**: Tracks crawl progress (queued → running → completed/error)

### Performance Features
- **Concurrent Processing**: Limited concurrent requests to avoid overwhelming targets
- **Timeout Handling**: 30-second timeout for page fetches, 10-second for link checks
- **Graceful Error Handling**: Comprehensive error reporting and recovery
- **Background Processing**: Non-blocking crawl execution
- **Database Persistence**: All results stored in MySQL for analysis

### Scalability
- **Pagination**: API responses support pagination and filtering
- **Bulk Operations**: Add multiple URLs, delete, or re-crawl in batches
- **Search & Filter**: Full-text search across URLs and titles
- **Statistics Dashboard**: Real-time stats on crawl queue and completion rates

## Development

For hot reloading during development:
```bash
go install github.com/cosmtrek/air@latest
make dev
```

## Production Deployment

Build optimized binary:
```bash
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o sykell-api main.go
```

Or use Docker:
```bash
docker build -t sykell-api .
docker run -p 8080:8080 sykell-api
```
