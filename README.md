# Sykell - Web Crawler Dashboard

A full-stack web crawler application with a React TypeScript frontend and Go backend, designed to analyze websites and provide detailed insights about their structure, links, and SEO elements.

## 🚀 Features

### Web Crawler Capabilities
- **HTML Analysis**: Extracts page titles, HTML version, and heading structure (H1-H6)
- **Link Analysis**: Identifies internal, external, and broken links
- **Login Form Detection**: Automatically detects login forms on pages
- **Real-time Progress**: Live updates during crawling process
- **Bulk Operations**: Crawl multiple URLs simultaneously

### Dashboard Features
- **Modern UI**: Clean, responsive interface built with Tailwind CSS
- **Data Visualization**: Charts and graphs for crawl results
- **Search & Filter**: Advanced filtering and search capabilities
- **Bulk Management**: Start, stop, and delete multiple crawl jobs
- **Detailed Analytics**: In-depth analysis of each crawled website

## 🏗️ Architecture

### Backend (Go)
- **Framework**: Fiber v2.50.0 (Express-like HTTP framework)
- **Database**: MySQL 8.0 with proper indexing
- **Authentication**: JWT-based authentication
- **Background Jobs**: Asynchronous crawling with job processor
- **API**: RESTful API with comprehensive endpoints

### Frontend (React TypeScript)
- **Framework**: React 19.1.0 with TypeScript
- **Build Tool**: Vite 4.5.3 for fast development
- **Styling**: Tailwind CSS v4.1.11 for modern design
- **Routing**: React Router Dom for SPA navigation
- **Charts**: Recharts for data visualization
- **Icons**: Heroicons for consistent iconography

## 📦 Project Structure

```
sykell/
├── backend/                 # Go backend application
│   ├── internal/
│   │   ├── handler/        # HTTP request handlers
│   │   ├── middleware/     # Authentication & CORS middleware
│   │   ├── models/         # Data models and structs
│   │   ├── router/         # Route definitions
│   │   └── service/        # Business logic and crawler service
│   ├── pkg/
│   │   ├── database/       # MySQL database connection and queries
│   │   └── logger/         # Structured logging with Zap
│   ├── scripts/            # Database setup scripts
│   ├── docker-compose.yml  # Docker services configuration
│   ├── Dockerfile          # Backend container image
│   ├── go.mod              # Go module dependencies
│   └── main.go             # Application entry point
├── frontend/               # React TypeScript frontend
│   ├── src/
│   │   ├── components/     # React components
│   │   ├── services/       # API service layer
│   │   ├── types/          # TypeScript type definitions
│   │   ├── utils/          # Utility functions and data transforms
│   │   ├── App.tsx         # Main application component
│   │   └── main.tsx        # Application entry point
│   ├── public/             # Static assets
│   ├── package.json        # Node.js dependencies and scripts
│   ├── tailwind.config.js  # Tailwind CSS configuration
│   ├── tsconfig.json       # TypeScript configuration
│   └── vite.config.ts      # Vite build configuration
└── README.md               # Project documentation
```

## 🛠️ Installation & Setup

### Prerequisites
- **Go 1.21+** - [Download Go](https://golang.org/dl/)
- **Node.js 18+** - [Download Node.js](https://nodejs.org/)
- **Docker & Docker Compose** - [Download Docker](https://www.docker.com/get-started)

### Backend Setup

1. **Navigate to backend directory**
   ```bash
   cd backend
   ```

2. **Install Go dependencies**
   ```bash
   go mod tidy
   ```

3. **Start MySQL database**
   ```bash
   docker-compose up -d mysql
   ```

4. **Run the backend server**
   ```bash
   go run main.go
   # or using the Makefile
   make run
   ```

The backend will be available at `http://localhost:8080`

### Frontend Setup

1. **Navigate to frontend directory**
   ```bash
   cd frontend
   ```

2. **Install Node.js dependencies**
   ```bash
   npm install
   ```

3. **Start the development server**
   ```bash
   npm run dev
   ```

The frontend will be available at `http://localhost:5174`

## 🔧 API Endpoints

### Authentication
- `POST /api/login` - User authentication
- `POST /api/users` - User registration

### Crawler Management
- `POST /api/crawler/urls` - Add single URL for crawling
- `POST /api/crawler/urls/bulk` - Add multiple URLs
- `GET /api/crawler/urls` - Get all crawl URLs (with pagination)
- `GET /api/crawler/urls/:id` - Get detailed crawl result
- `POST /api/crawler/urls/:id/crawl` - Start crawling a URL
- `DELETE /api/crawler/urls` - Delete multiple URLs
- `POST /api/crawler/urls/recrawl` - Re-crawl multiple URLs
- `GET /api/crawler/stats` - Get crawling statistics

## 🔐 Environment Configuration

### Backend (.env)
```env
# Database Configuration
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=password
DB_NAME=sykell_db

# JWT Configuration
JWT_SECRET=your-secret-key

# Server Configuration
PORT=8080
ENV=development
```

## 🐳 Docker Deployment

### Full Stack with Docker Compose
```bash
docker-compose up -d
```

This will start:
- MySQL database on port 3306
- Go backend on port 8080
- (Frontend needs to be built separately for production)

### Individual Services
```bash
# MySQL only
docker-compose up -d mysql

# Backend only (requires manual build)
docker build -t sykell-backend .
docker run -p 8080:8080 sykell-backend
```

## 🚦 Usage

1. **Access the Application**
   - Open `http://localhost:5174` in your browser
   - Create an account or login

2. **Add URLs for Crawling**
   - Click "Add URL" button
   - Enter the website URL to analyze
   - Start crawling process

3. **View Results**
   - Monitor crawling progress in real-time
   - View detailed analytics for each website
   - Export or manage crawl results

4. **Bulk Operations**
   - Select multiple URLs
   - Perform bulk start, stop, or delete operations

## 🧪 Development

### Backend Development
```bash
# Run with hot reload (requires air)
go install github.com/cosmtrek/air@latest
make dev

# Run tests
make test

# Build for production
make build
```

### Frontend Development
```bash
# Development server with hot reload
npm run dev

# Type checking
npm run type-check

# Linting
npm run lint

# Build for production
npm run build
```
