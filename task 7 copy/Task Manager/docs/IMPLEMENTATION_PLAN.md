# Task Management API - Implementation Details

## API Endpoints and JSON Input Requirements

### 1. User Registration - `POST /register`

- **Description:** Registers a new user.
- **JSON Input:**

  ```json
  {
    "username": "string",   // required
    "password": "string",   // required
    "email": "string",      // optional
    "role": "string"        // optional, e.g., "admin", "user"
  }
  ```

- **Authentication:** No authentication required.

### 2. User Login - `POST /login`

- **Description:** Authenticates a user and returns a JWT token.
- **JSON Input:**

  ```json
  {
    "username": "string",   // required
    "password": "string"    // required
  }
  ```

- **Authentication:** No authentication required.

### 3. Task Endpoints (Require JWT Authentication)

All task-related endpoints require a valid JWT token in the `Authorization` header as a Bearer token.

#### a. Get All Tasks - `GET /tasks`

- **Description:** Retrieves all tasks.
- **Authentication:** Required.

#### b. Get Task by ID - `GET /tasks/:id`

- **Description:** Retrieves a task by its ID.
- **Authentication:** Required.

#### c. Create Task - `POST /tasks`

- **Description:** Creates a new task.
- **JSON Input:**

  ```json
  {
    "title": "string",          // required
    "description": "string",    // optional
    "due_date": "dd-mm-yyyy",   // required, format: 02-01-2006
    "status": "string"          // required, one of "Pending", "In Progress", "Completed"
  }
  ```

- **Authentication:** Required.

#### d. Update Task - `PUT /tasks/:id`

- **Description:** Updates an existing task.
- **JSON Input:** Same as Create Task, all fields optional.
- **Authentication:** Required.

#### e. Delete Task - `DELETE /tasks/:id`

- **Description:** Deletes a task by its ID.
- **Authentication:** Required.

## Authentication

- JWT tokens are issued upon successful login.
- Tokens must be included in the `Authorization` header as `Bearer <token>` for all protected endpoints.
- The JWT secret is managed via environment variables.

## Environment Configuration

- Sensitive data such as JWT secret and MongoDB connection URI are stored in a `.env` file.
- The application loads these variables at startup.
