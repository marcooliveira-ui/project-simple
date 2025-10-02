# Car Management API

A RESTful API built with Go (Golang) and Gin framework for managing cars with full CRUD operations, pagination, and Swagger documentation.

## Features

- ✅ RESTful API with proper HTTP verbs (GET, POST, PUT, DELETE)
- ✅ CRUD operations for Car entity
- ✅ Pagination with customizable page size and sorting
- ✅ UUID-based identifiers
- ✅ GORM ORM with PostgreSQL
- ✅ Clean Architecture (Domain, Repository, Service, Handler layers)
- ✅ Swagger/OpenAPI documentation
- ✅ Input validation
- ✅ Error handling middleware
- ✅ CORS support
- ✅ Custom logger middleware
- ✅ Environment-based configuration

## Tech Stack

- **Language**: Go 1.25.1
- **Framework**: Gin (v1.11.0)
- **ORM**: GORM (v1.31.0)
- **Database**: PostgreSQL
- **Documentation**: Swagger/OpenAPI (swaggo)
- **UUID**: google/uuid

## Project Structure

```
project-simple/
├── cmd/
│   └── api/
│       └── main.go              # Application entry point
├── internal/
│   ├── config/                  # Configuration management
│   │   └── config.go
│   ├── domain/
│   │   ├── dto/                 # Data Transfer Objects
│   │   │   └── car_dto.go
│   │   └── entity/              # Domain entities
│   │       └── car.go
│   ├── handler/                 # HTTP handlers (controllers)
│   │   ├── car_handler.go
│   │   └── health_handler.go
│   ├── infrastructure/
│   │   └── database/            # Database setup and migrations
│   │       └── database.go
│   ├── middleware/              # Custom middlewares
│   │   ├── cors.go
│   │   ├── error_handler.go
│   │   └── logger.go
│   ├── repository/              # Data access layer
│   │   └── car_repository.go
│   ├── router/                  # Route definitions
│   │   └── router.go
│   └── service/                 # Business logic layer
│       └── car_service.go
├── pkg/
│   └── response/                # Response utilities
│       ├── response.go
│       └── error_response.go
├── .env.example                 # Environment variables template
├── .gitignore
├── go.mod
├── go.sum
└── README.md
```

## Getting Started

### Prerequisites

- Go 1.25.1 or higher
- PostgreSQL 12 or higher
- Git

### Installation

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd project-simple
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Create database**
   ```bash
   psql -U postgres
   CREATE DATABASE car_db;
   \q
   ```

4. **Configure environment variables**
   ```bash
   cp .env.example .env
   ```

   Edit `.env` with your database credentials:
   ```env
   DB_HOST=localhost
   DB_PORT=5432
   DB_USER=postgres
   DB_PASSWORD=your_password
   DB_NAME=car_db
   DB_SSLMODE=disable

   SERVER_PORT=8080
   SERVER_ENV=development
   ```

5. **Install Swagger CLI (optional, for regenerating docs)**
   ```bash
   go install github.com/swaggo/swag/cmd/swag@latest
   ```

6. **Generate Swagger documentation**
   ```bash
   swag init -g cmd/api/main.go -o docs
   ```

7. **Run the application**
   ```bash
   go run cmd/api/main.go
   ```

   The server will start on `http://localhost:8080`

## API Documentation

### Swagger UI

Access the interactive API documentation at:
```
http://localhost:8080/swagger/index.html
```

### Endpoints

#### Health Check
- `GET /api/v1/health` - Check API health status

#### Cars

- `POST /api/v1/cars` - Create a new car
- `GET /api/v1/cars` - Get all cars (with pagination)
- `GET /api/v1/cars/:id` - Get a specific car by ID
- `PUT /api/v1/cars/:id` - Update a car
- `DELETE /api/v1/cars/:id` - Delete a car

### Examples

#### Create a Car
```bash
curl -X POST http://localhost:8080/api/v1/cars \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Honda Civic",
    "engine_version": "2.0"
  }'
```

#### Get All Cars (with pagination)
```bash
curl "http://localhost:8080/api/v1/cars?page=1&page_size=10&sort_by=created_at&sort_dir=desc"
```

#### Get Car by ID
```bash
curl http://localhost:8080/api/v1/cars/{car-uuid}
```

#### Update a Car
```bash
curl -X PUT http://localhost:8080/api/v1/cars/{car-uuid} \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Honda Civic Sport",
    "engine_version": "2.0"
  }'
```

#### Delete a Car
```bash
curl -X DELETE http://localhost:8080/api/v1/cars/{car-uuid}
```

## Car Entity

### Fields
- `id` (UUID) - Unique identifier
- `name` (string) - Car name (required, 2-100 characters)
- `engine_version` (string) - Engine version (required, must be one of: 1.0, 1.4, 1.5, 1.6, 1.8, 2.0, 2.4, 2.5, 3.0, 3.5, 4.0)
- `created_at` (timestamp) - Creation timestamp
- `updated_at` (timestamp) - Last update timestamp

### Validation Rules
- **Name**: Required, min 2 characters, max 100 characters
- **Engine Version**: Required, must be one of the allowed values

## Pagination

All list endpoints support pagination with the following query parameters:

- `page` (int) - Page number (default: 1, min: 1)
- `page_size` (int) - Items per page (default: 10, min: 1, max: 100)
- `sort_by` (string) - Sort field: `name`, `engine_version`, `created_at` (default: `created_at`)
- `sort_dir` (string) - Sort direction: `asc`, `desc` (default: `desc`)

### Pagination Response Format
```json
{
  "message": "Cars retrieved successfully",
  "data": {
    "data": [
      {
        "id": "uuid",
        "name": "Honda Civic",
        "engine_version": "2.0",
        "created_at": "2024-01-01T10:00:00Z",
        "updated_at": "2024-01-01T10:00:00Z"
      }
    ],
    "pagination": {
      "current_page": 1,
      "page_size": 10,
      "total_pages": 5,
      "total_records": 50
    }
  }
}
```

## Development

### Running Tests
```bash
go test ./...
```

### Building for Production
```bash
go build -o bin/api cmd/api/main.go
```

### Regenerating Swagger Docs
After modifying API documentation comments:
```bash
swag init -g cmd/api/main.go -o docs
```

## Architecture

This project follows **Clean Architecture** principles:

1. **Domain Layer** (`internal/domain`): Core business entities and DTOs
2. **Repository Layer** (`internal/repository`): Data access and persistence
3. **Service Layer** (`internal/service`): Business logic
4. **Handler Layer** (`internal/handler`): HTTP request handlers
5. **Infrastructure Layer** (`internal/infrastructure`): External dependencies (database, etc.)

### Dependency Flow
```
Handler → Service → Repository → Database
   ↓         ↓          ↓
  DTO    Business    Entity
         Logic
```

## Error Handling

The API uses consistent error responses:

```json
{
  "error": "Error Type",
  "message": "Human-readable error message",
  "details": "Additional error details (optional)"
}
```

### HTTP Status Codes
- `200 OK` - Successful GET/PUT
- `201 Created` - Successful POST
- `204 No Content` - Successful DELETE
- `400 Bad Request` - Invalid request format
- `404 Not Found` - Resource not found
- `422 Unprocessable Entity` - Validation failed
- `500 Internal Server Error` - Server error

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License.

## Contact

For support or questions, please contact: support@carapi.com
