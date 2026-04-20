# Enrollment System Backend (Golang)

This is the skeleton for the backend of the school enrollment system.

## Tech Stack
- **Go** 1.21+
- **GoFiber** (Web Framework)
- **GORM** (ORM)
- **MySQL** (Database)
- **golang-migrate** (Database Migrations)
- **JWT** (Authentication)
- **bcrypt** (Password Hashing)
- **zap** (Logging)
- **Cloudflare R2** (Storage - Planned)

## Project Structure

```text
.
├── cmd/server           # Main applications for this project
├── internal             # Private application and library code
│   ├── common           # Common structs like Response, Pagination
│   ├── config           # Configuration loader
│   ├── database         # Database connection logic
│   ├── middleware       # HTTP middlewares (CORS, Logger, Recovery)
│   └── modules          # Feature modules (e.g., health, auth, users)
├── migrations           # SQL migration files
├── pkg                  # Public library code (e.g., Logger wrapper)
└── README.md
```

## How to run locally

### 1. Requirements
Ensure you have the following installed:
- Go (1.21+)
- MySQL server (or dockerized MySQL)
- `golang-migrate` CLI tool (optional, for running migrations)

### 2. Setup Database
Create a database in your local MySQL instance:
```sql
CREATE DATABASE enrollment_db;
```

### 3. Environment Variables
Copy the example config to a real config file. Note that your actual `.env` will be ignored by git.
```sh
cp .env.example .env
```
Make sure to update `.env` with your actual local database credentials.

### 4. Install Dependencies
```sh
go mod tidy
```

### 5. Run the server
```sh
go run cmd/server/main.go
```
The server will start on the port defined in your `.env` (default is 8080).

### 6. Verify setup
Visit the health check endpoint to confirm everything is running:
```sh
curl http://localhost:8080/api/v1/health
```
