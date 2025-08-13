# Go Task Manager

A simple Task Manager backend built with Go, MongoDB, and Clean Architecture principles.

## Features
- User registration and authentication
- Task CRUD operations (create, read, update, delete)
- JWT-based authentication
- Clean separation of concerns (Domain, Usecases, Repositories, Delivery)

## Project Structure

```
Task Manager/
├── main.go                  # Application entry point
├── Domain/                  # Domain models and interfaces
├── Usecases/                # Business logic
├── Repositories/            # Data access and adapters
├── Infrastructure/          # DB, JWT, password services
├── Delivery/                # HTTP controllers and routers
├── tests/                   # Unit and integration tests
├── docs/                    # Architecture and implementation docs
```

## Setup
1. **Install Go** (v1.18+ recommended)
2. **Clone the repo**
3. **Set environment variables:**
   - `MONGODB_URI` - MongoDB connection string
   - `DATABASE_NAME` - Database name
   - `TASKS_COLLECTION` - Collection for tasks
   - `USERS_COLLECTION` - Collection for users
4. **Run the app:**
   ```bash
   go run main.go
   ```
5. **Run tests:**
   ```bash
   go test ./tests/...
   ```

## API Endpoints
- `POST /register` - Register a new user
- `POST /login` - Authenticate user
- `GET /tasks` - List tasks
- `POST /tasks` - Create task
- `GET /tasks/{id}` - Get task by ID
- `PUT /tasks/{id}` - Update task
- `DELETE /tasks/{id}` - Delete task

## License
MIT
