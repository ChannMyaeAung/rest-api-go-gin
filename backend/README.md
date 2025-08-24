# REST API with Go, Gin & SQLite

A production-ready REST API for event management with user authentication, built with Go, Gin framework, and SQLite database. Features JWT authentication, event CRUD operations, and attendee management.

## 🚀 Quick Start

1. **Clone and setup:**

   ```bash
   git clone https://github.com/ChannMyaeAung/rest-api-go-gin.git
   cd rest-api-go-gin
   go mod tidy
   ```

2. **Run database migrations:**

   ```bash
   go run ./cmd/migrate up
   ```

3. **Start the server:**

   ```bash
   go run ./cmd/api
   # Or with live reload: air
   ```

4. **Access Swagger documentation:**
   ```
   http://localhost:8080/swagger/index.html
   ```

## 📖 API Documentation

**Live API Documentation:** [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)

The API provides comprehensive Swagger/OpenAPI documentation that you can use to:

- View all available endpoints
- Test API calls directly in the browser
- See request/response schemas
- Understand authentication requirements

### Key Endpoints:

| Method   | Endpoint                                 | Description       | Auth Required |
| -------- | ---------------------------------------- | ----------------- | ------------- |
| `POST`   | `/api/v1/auth/register`                  | Register new user | No            |
| `POST`   | `/api/v1/auth/login`                     | Login user        | No            |
| `GET`    | `/api/v1/events`                         | List all events   | No            |
| `POST`   | `/api/v1/events`                         | Create event      | Yes           |
| `GET`    | `/api/v1/events/{id}`                    | Get event details | No            |
| `PUT`    | `/api/v1/events/{id}`                    | Update event      | Yes (Owner)   |
| `DELETE` | `/api/v1/events/{id}`                    | Delete event      | Yes (Owner)   |
| `POST`   | `/api/v1/events/{id}/attendees/{userId}` | Add attendee      | Yes (Owner)   |
| `GET`    | `/api/v1/events/{id}/attendees`          | List attendees    | No            |
| `DELETE` | `/api/v1/events/{id}/attendees/{userId}` | Remove attendee   | Yes (Owner)   |

### Authentication

The API uses JWT tokens for authentication. Include the token in requests:

```
Authorization: Bearer <your_jwt_token>
```

## 🧪 Testing the API

### Option 1: Swagger UI (Recommended)

Visit `http://localhost:8080/swagger/index.html` and use the interactive documentation to test all endpoints.

### Option 2: cURL Examples

**Register a user:**

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email": "test@example.com", "password": "password123", "name": "Test User"}'
```

**Login:**

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "test@example.com", "password": "password123"}'
```

**Create an event (requires authentication):**

```bash
curl -X POST http://localhost:8080/api/v1/events \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{"name": "Go Conference", "ownerId": 1, "description": "A conference about Go", "date": "2025-05-20", "location": "San Francisco"}'
```

## 🛠 Development Setup

### Prerequisites

- Go 1.24+
- (Optional) Air for live reload during development: https://github.com/cosmtrek/air or https://github.com/air-verse/air
- golang-migrate CLI for creating/applying migrations: https://github.com/golang-migrate/migrate

### Install tools (examples, Windows PowerShell)

- Install golang-migrate (binary):

  ```powershell
  scoop install migrate  # if using scoop
  # or download from https://github.com/golang-migrate/migrate/releases
  ```

- Install Air (optional):

  ```powershell
  scoop install air
  # or follow project README
  ```

- Install Swaggo (for API documentation):
  ```powershell
  go install github.com/swaggo/swag/cmd/swag@latest
  ```

### Create new migrations

From project root, create migration SQL files in the migrations folder:

```powershell
migrate create -ext sql -dir ./cmd/migrate/migrations -seq create_users_table
migrate create -ext sql -dir ./cmd/migrate/migrations -seq create_events_table
migrate create -ext sql -dir ./cmd/migrate/migrations -seq create_attendees_table
```

This creates `*.up.sql` and `*.down.sql` files under `cmd/migrate/migrations`.

### Run migrations

From the project root (so relative paths match), run the migrate command implemented in this repo:

- Up (apply all up migrations):

```powershell
go run ./cmd/migrate up
```

- Down (apply all down migrations / rollback):

```powershell
go run ./cmd/migrate down
```

### Generate API documentation

```powershell
swag init --dir cmd/api --parseDependency --parseInternal --parseDepth 1
```

Notes:

- The program opens `./data.db` by default and reads migrations from `cmd/migrate/migrations`. Run commands from the repository root or adjust paths accordingly.
- The SQLite driver must be available at build time (commonly `github.com/mattn/go-sqlite3` or an alternative pure-Go driver `modernc.org/sqlite`). If you prefer a pure-Go driver to avoid cgo, replace the driver in the code and update imports.

## 🏗 Project Structure

```
rest-api-go-gin/
├── cmd/
│   ├── api/           # HTTP server and handlers
│   │   ├── main.go    # Application entry point
│   │   ├── auth.go    # Authentication handlers
│   │   ├── events.go  # Event CRUD handlers
│   │   ├── users.go   # User handlers
│   │   ├── routes.go  # Route definitions
│   │   ├── server.go  # Server configuration
│   │   ├── middleware.go # JWT middleware
│   │   └── context.go # Context helpers
│   └── migrate/       # Database migration tool
│       ├── main.go    # Migration runner
│       └── migrations/ # SQL migration files
├── internal/
│   ├── database/      # Database models and operations
│   │   ├── models.go  # Database connection setup
│   │   ├── users.go   # User database operations
│   │   ├── events.go  # Event database operations
│   │   └── attendees.go # Attendee database operations
│   └── env/           # Environment configuration
├── docs/              # Swagger documentation (auto-generated)
└── README.md
```

## 🔧 Module / dependency tips

- To clean unused indirect modules:

```powershell
go mod tidy
```

- To check why a module is present:

```powershell
go mod why -m github.com/mattn/go-sqlite3
```

Keep migration SQLs in version control; do not commit the database file (`data.db`) if it is user/local state. Add `data.db` to `.gitignore` if appropriate.

## 🚀 Deployment

For production deployment, consider:

1. **Environment Variables:**

   - Set `GIN_MODE=release`
   - Use secure `JWT_SECRET`
   - Configure appropriate database URL

2. **Database:**

   - Consider PostgreSQL or MySQL for production
   - Set up proper database connections and pooling

3. **Security:**
   - Enable HTTPS
   - Set up CORS properly
   - Add rate limiting
   - Use secure headers middleware

## 📝 License

This project is for demonstration purposes.
go run ./cmd/migrate up

````

- Down (apply all down migrations / rollback):

```powershell
go run ./cmd/migrate down
````

Notes:

- The program opens `./data.db` by default and reads migrations from `cmd/migrate/migrations`. Run commands from the repository root or adjust paths accordingly.
- The SQLite driver must be available at build time (commonly `github.com/mattn/go-sqlite3` or an alternative pure-Go driver `modernc.org/sqlite`). If you prefer a pure-Go driver to avoid cgo, replace the driver in the code and update imports.

## Module / dependency tips

- To clean unused indirect modules:

```powershell
go mod tidy
```

- To check why a module is present:

```powershell
go mod why -m github.com/modernc.org/sqlite
```
