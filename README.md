# Project Title

## Description

This project is a full-featured backend API built using the Go standard `net/http` package. It is designed as a CRUD (Create, Read, Update, Delete) application with comprehensive authentication mechanisms. The project is structured to demonstrate the use of Go's standard library to build robust and efficient web applications without relying on external web frameworks. It includes user authentication and authorization, employing modern security practices like JWT signing and verification and secure password storage.

## Getting Started

### Prerequisites

Ensure you have the following installed on your system:

- Go 1.23 or later
- PostgreSQL
- Redis
- Docker

### Clone the Repository

git clone https://github.com/yourusername/your-repo-name.git
cd your-repo-name

### Install Packages

go mod download

### Set Up Environment Variables

Create a .env file in the root directory and add the necessary environment variables:

```bash
DB_HOST=localhost
DB_PORT=5432
DB_USER=your_db_user
DB_PASSWORD=your_db_password
DB_NAME=your_db_name
REDIS_HOST=localhost
REDIS_PORT=6379
JWT_SECRET=your_jwt_secret
```

### Run Database Migrations

Ensure that PostgreSQL is running and then run the database migrations:

go run scripts/migrate.go

### Run Tests

To run the tests, use the following command:

go test ./...

### Run the Server

To start the server, use the following command:

go run main.go

## Tech Stacks Used

- **Go 1.23**: The programming language used for developing the backend API.
- **`database/sql`**: The built-in Go package for database interactions.
- **PostgreSQL**: The relational database management system used for data storage.
- **Redis**: The in-memory data structure store used for caching and session management.
- **Docker**: The platform used for containerizing the application.

## JWT Authentication Mechanism

The project uses JSON Web Tokens (JWT) for authentication. The tokens are signed and verified using the EDDSA algorithm with public and private keys. This ensures that the tokens are secure and cannot be tampered with.

## Password Storage Mechanism

Passwords are hashed using the bcrypt algorithm before being stored in the database. Bcrypt is a secure hashing algorithm designed to be computationally intensive, making it difficult for attackers to brute force the hashed passwords.

## FAQ

### Why use the standard `net/http` package instead of web frameworks like Gin?

Using the standard `net/http` package allows for greater control over the request handling process and reduces dependencies. It also provides a better understanding of how HTTP works under the hood, which can be beneficial for learning and for situations where performance and customization are critical.

### Why not use an ORM?

Using the `database/sql` package directly provides more control over database interactions and can lead to more efficient and optimized queries. It avoids the overhead that ORMs can introduce and allows for the use of raw SQL queries, which can be more performant in certain scenarios.
