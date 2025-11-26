# API Server - Production-Ready Go REST API

Industry-standard Go REST API with JWT authentication, PostgreSQL, and clean architecture.

## 🏗️ Architecture

```
Handler → Service → Repository → Database
```

- **Handler**: HTTP layer (request/response)
- **Service**: Business logic and validation
- **Repository**: Database operations
- **Models**: Data structures

## 📁 Project Structure

```
apiserver/
├── cmd/api/main.go              # Application entry point
├── internal/
│   ├── config/                  # Configuration management
│   ├── handlers/                # HTTP handlers
│   ├── middleware/              # Auth middleware
│   ├── models/                  # Data models
│   ├── repository/              # Database layer
│   ├── service/                 # Business logic
│   └── utils/                   # Utilities (JWT, password, response)
├── migrations/                  # SQL migrations
├── docker-compose.yml           # Docker setup
├── .env                         # Environment variables
└── go.mod                       # Go dependencies
```

## 🚀 Quick Start

### 1. Start PostgreSQL

```bash
docker-compose up -d
```

### 2. Run Migrations

```bash
./migrate.sh up
```

### 3. Install Dependencies

```bash
go mod tidy
```

### 4. Run the Server

```bash
go run cmd/api/main.go
```

Server runs on: `http://localhost:8080`

## 📡 API Endpoints

### Public Endpoints

#### Register User
```bash
POST /api/v1/register
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123",
  "name": "John Doe"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "email": "user@example.com",
      "name": "John Doe",
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  }
}
```

#### Login
```bash
POST /api/v1/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123"
}
```

### Protected Endpoints (Require JWT Token)

Add token to header: `Authorization: Bearer <your-token>`

#### Get Current User
```bash
GET /api/v1/me
Authorization: Bearer <token>
```

#### Update Current User
```bash
PUT /api/v1/me
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "Jane Doe",
  "email": "jane@example.com"
}
```

#### Delete Current User
```bash
DELETE /api/v1/me
Authorization: Bearer <token>
```

#### Get All Users
```bash
GET /api/v1/users
Authorization: Bearer <token>
```

#### Get User by ID
```bash
GET /api/v1/users/{id}
Authorization: Bearer <token>
```

## 🧪 Testing with cURL

### Register
```bash
curl -X POST http://localhost:8080/api/v1/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123",
    "name": "Test User"
  }'
```

### Login
```bash
curl -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }'
```

### Get Current User (with token)
```bash
curl -X GET http://localhost:8080/api/v1/me \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

## 🗃️ Database Migrations

### Create New Migration
```bash
./migrate.sh create add_user_avatar
```

### Apply Migrations
```bash
./migrate.sh up
```

### Rollback Last Migration
```bash
./migrate.sh down
```

### Check Migration Version
```bash
./migrate.sh version
```

## 🔧 Environment Variables

Edit `.env` file:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=apiserver
DB_SSLMODE=disable

JWT_SECRET=your-super-secret-key-change-this-in-production

PORT=8080
```

## 🐳 Docker Commands

```bash
# Start database
docker-compose up -d

# Stop database
docker-compose down

# View logs
docker-compose logs -f postgres

# Reset database (delete all data)
docker-compose down -v

# Connect to database
docker exec -it apiserver-db psql -U postgres -d apiserver
```

## 📦 Dependencies

- **Chi** - HTTP router
- **sqlx** - SQL database toolkit
- **JWT** - JSON Web Tokens
- **bcrypt** - Password hashing
- **godotenv** - Environment variables
- **validator** - Request validation

## 🛡️ Security Features

- ✅ Password hashing with bcrypt
- ✅ JWT authentication
- ✅ SQL injection prevention (parameterized queries)
- ✅ CORS configuration
- ✅ Request validation
- ✅ Secure password handling (never returned in responses)

## 📚 Best Practices Used

1. **Clean Architecture** - Separation of concerns
2. **Repository Pattern** - Abstracted database layer
3. **Dependency Injection** - Testable code
4. **Interface-based Design** - Flexible and maintainable
5. **Error Handling** - Proper error messages
6. **Migrations** - Version-controlled database schema
7. **Environment Configuration** - 12-factor app principles
8. **Standardized Responses** - Consistent API responses

## 🔍 Health Check

```bash
curl http://localhost:8080/health
```

Response:
```json
{
  "status": "ok",
  "message": "Server is running"
}
```

## 📖 Learning Resources

- Go documentation: https://go.dev/doc/
- Chi router: https://go-chi.io/
- sqlx: http://jmoiron.github.io/sqlx/
- JWT: https://jwt.io/

## 🤝 Contributing

This is a learning project demonstrating Go best practices for REST API development.

## 📝 License

MIT License - feel free to use for learning and projects.