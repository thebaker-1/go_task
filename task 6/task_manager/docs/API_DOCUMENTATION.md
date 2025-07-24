# Task Management API Documentation

## Authentication and Authorization

This API uses JSON Web Tokens (JWT) for authentication and authorization.

### 1. User Registration

**Endpoint:** `POST /register`  
**Description:** Creates a new user account.  
**Request Body (JSON):**

```json
{
  "username": "newuser",
  "password": "strongpassword123",
  "email": "user@example.com"
}
```

**Success Response (201 Created):**

```json
{
  "message": "User registered successfully"
}
```

### 2. User Login

**Endpoint:** `POST /login`  
**Description:** Authenticates a user and returns a JWT token. This token must be included in the Authorization header for all protected routes.  
**Request Body (JSON):**

```json
{
  "username": "existinguser",
  "password": "theirpassword"
}
```

**Success Response (200 OK):**

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2N..."
}
```

### 3. Protected Endpoints

To access protected endpoints, you must include the JWT token obtained from the login endpoint in the Authorization header of your requests, in the format:

``` bash
Authorization: Bearer <YOUR_JWT_TOKEN>
```

### 4. Task Endpoints (Protected - All Authenticated Users)

- `GET /tasks`: Retrieve all tasks.
- `GET /tasks/:id`: Retrieve a specific task by ID.
- `POST /tasks`: Create a new task.

### 5. Admin-Only Endpoints (Protected - Requires 'admin' role)

- `PUT /tasks/:id`: Update an existing task.
- `DELETE /tasks/:id`: Delete a task.

### 6. Role-Based Access

- Users with role `"user"` can access general task endpoints.
- Users with role `"admin"` have additional permissions to update and delete tasks.
