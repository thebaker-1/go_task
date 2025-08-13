# Task Management API - Clean Architecture Documentation

## Overview

This project implements a Task Management API following Clean Architecture principles to ensure maintainability, testability, and scalability. The codebase is organized into distinct layers with clear separation of concerns.

## Architecture Layers

### 1. Domain

- Contains core business entities: `Task` and `User`.
- Entities are pure Go structs without any serialization or persistence tags.
- Represents the business rules and logic independent of external frameworks.

### 2. Usecases

- Implements application-specific business rules.
- Contains use case interfaces and their implementations for tasks and users.
- Orchestrates interactions between domain entities and repositories.
- Defines Data Transfer Objects (DTOs) for communication with delivery and persistence layers.

### 3. Repositories

- Abstracts data access logic.
- Defines repository interfaces for tasks and users.
- Implements MongoDB-based repositories using the official MongoDB Go driver.
- Database connection and collection initialization are encapsulated in the Infrastructure layer.

### 4. Infrastructure

- Implements external dependencies and services.
- Includes MongoDB client setup (`database.go`), JWT services, password hashing, and authentication middleware.
- Manages environment variables and configuration.

### 5. Delivery

- Handles HTTP requests and responses.
- Implements controllers using the Gin framework.
- Converts between domain models and DTOs for API communication.
- Sets up routing and middleware.

## Configuration

- Sensitive configuration such as JWT secret and MongoDB URI are stored in a `.env` file.
- The application loads environment variables at startup.

## Key Design Decisions

- Domain models are kept free of serialization tags to maintain purity.
- DTOs with JSON and BSON tags are defined in the Delivery layer for API and persistence mapping.
- Dependency injection is used to wire repositories, usecases, and controllers.
- JWT authentication is implemented with middleware to protect task-related endpoints.
- The database setup is encapsulated in the Infrastructure layer for separation of concerns.

## Running the Application

1. Ensure MongoDB is running and accessible at the configured URI.
2. Set environment variables in the `.env` file.
3. Build and run the application using `go run main.go`.
4. The API listens on port 8080.

## API Endpoints

- `POST /register` - Register a new user.
- `POST /login` - Authenticate user and receive JWT token.
- `GET /tasks` - Get all tasks (requires JWT).
- `GET /tasks/:id` - Get task by ID (requires JWT).
- `POST /tasks` - Create a new task (requires JWT).
- `PUT /tasks/:id` - Update a task (requires JWT).
- `DELETE /tasks/:id` - Delete a task (requires JWT).

## Testing

- Testing is recommended for all endpoints, including happy paths, error paths, and edge cases.
- Consider using tools like Postman or Curl for manual testing.
- Automated unit and integration tests can be added for usecases and repositories.

## Future Improvements

- Add comprehensive unit and integration tests.
- Implement role-based access control.
- Add pagination and filtering for task lists.
- Improve error handling and logging.
