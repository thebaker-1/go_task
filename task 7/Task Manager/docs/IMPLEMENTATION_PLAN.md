# Task Management API - Implementation Details

1. User Registration - POST /register
Description: Registers a new user.

JSON Input:

{
  "username": "string",   // required
  "password": "string",   // required
  "email": "string",      // optional
  "role": "string"        // optional, e.g., "admin", "user"
}

JSON Output:

{
  "message": "User registered successfully"
}

Authentication: No authentication required.

2.User Login - POST /login
Description: Authenticates a user and returns a JWT token.

JSON Input:

{
  "username": "string",   // required
  "password": "string"    // required
}

JSON Output:

{
  "token": "jwt-token-string"
}

Authentication: No authentication required.

3.Task Endpoints (Require JWT Authentication)
All task-related endpoints require a valid JWT token in the Authorization header as a Bearer token.

a. Get All Tasks - GET /tasks
Description: Retrieves all tasks.

Authentication: Required.

JSON Input:

{}

JSON Output:

[
  {
    "id": "string",
    "title": "string",
    "description": "string",
    "due_date": "dd-mm-yyyy",
    "status": "string"
  }
]

b. Get Task by ID - GET /tasks/:id
Description: Retrieves a task by its ID.

Authentication: Required.

JSON Output:

{
  "id": "string",
  "title": "string",
  "description": "string",
  "due_date": "dd-mm-yyyy",
  "status": "string"
}

c. Create Task - POST /tasks
Description: Creates a new task.

JSON Input:

{
  "title": "string",          // required
  "description": "string",    // optional
  "due_date": "dd-mm-yyyy",   // required, format: 02-01-2006
  "status": "string"          // required, one of "Pending", "In Progress", "Completed"
}

JSON Output:

{
  "id": "6887c474a264ed9d9ee52053",
  "title": "Finish report",
  "description": "Complete the quarterly report",
  "due_date": "30-09-2025",
  "status": "Pending"
}

Authentication: Required.

d. Update Task - PUT /tasks/:id
Description: Updates an existing task.

JSON Input: Same as Create Task, all fields optional.

JSON Output:

{
  "id": "string",
  "title": "string",
  "description": "string",
  "due_date": "dd-mm-yyyy",
  "status": "string"
}

Authentication: Required.

e. Delete Task - DELETE /tasks/:id
Description: Deletes a task by its ID.

Authentication: Required.

JSON Output:

{
  "message": "deleted task successfully"
}

Authentication
JWT tokens are issued upon successful login.

Tokens must be included in the Authorization header as Bearer token for all protected endpoints.

The JWT secret is managed via environment variables.

Environment Configuration
Sensitive data such as JWT secret and MongoDB connection URI are stored in a .env file.

The application loads these variables at startup.

The JWT secret is loaded early in the application lifecycle to ensure it is available for JWT services, preventing errors related to empty secrets.

Domain models use MongoDB's primitive.ObjectID for IDs to maintain consistency with the database.

Repository, usecase, and controller layers handle conversion between primitive.ObjectID and string representations for API communication.
