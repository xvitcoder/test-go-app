# Test Go App

This is a simple RESTful API built with Go, Gorilla Mux, and PostgreSQL. The application provides basic user management (CRUD) functionality, storing user data in a PostgreSQL database.

## Features
- RESTful endpoints for `/users` resource
- PostgreSQL-backed persistence using `database/sql` and `lib/pq`
- Database migrations included
- Simple, minimal dependencies

## Getting Started

1. **Clone the Repository**
   ```sh
   git clone git@github.com:xvitcoder/test-go-app.git
   cd test-go-app
   ```

2. **Run the migration**
   Apply the migration in `migrations/001_create_users_table.sql` to your PostgreSQL database.

3. **Configure Environment**
   - Copy/edit the `.env` file:
     ```sh
     POSTGRES_DSN=postgres://username:password@localhost:5432/yourdbname?sslmode=disable
     ```

4. **Run the Application**
   ```sh
   go run main.go
   ```

5. **API Endpoints**
   - `GET    /users`           - List all users
   - `GET    /users/{id}`      - Get a single user
   - `POST   /users`           - Create a user (JSON body)
   - `PUT    /users`           - Update a user (JSON body)
   - `DELETE /users/{id}`      - Delete user by ID

---

**Contributions and issues welcome!**
