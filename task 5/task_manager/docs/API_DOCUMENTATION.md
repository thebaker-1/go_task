# Task Manager API Documentation

## Overview

The Task Manager API provides a set of endpoints to manage tasks, including creating, retrieving, updating, and deleting tasks. It enables users to organize and track their tasks efficiently.

## Getting Started Guide

To start using the Task Manager API, you need to:

- Use HTTPS to send requests to the API endpoints.
- Send and receive data in JSON format.
- Handle rate limits and usage constraints as described below.

## API Endpoints

### Get All Tasks

- **Endpoint:** `GET /tasks`
- **Description:** Retrieves a list of all tasks.
- **Response:**
  - Status: 200 OK
  - Body: JSON array of task objects

```json
[
  {
    "id": "string",
    "title": "string",
    "description": "string",
    "due_date": "YYYY-MM-DDTHH:MM:SSZ",
    "status": "Pending|In Progress|Completed"
  }
]
```

### Get Task by ID

- **Endpoint:** `GET /tasks/{id}`
- **Description:** Retrieves a task by its ID.
- **Response:**
  - Status: 200 OK
  - Body: JSON task object
  - Status: 404 Not Found if task does not exist

```json
{
  "id": "string",
  "title": "string",
  "description": "string",
  "due_date": "YYYY-MM-DDTHH:MM:SSZ",
  "status": "Pending|In Progress|Completed"
}
```

- **Error Response:**

```json
{
  "message": "task not found"
}
```

### Create a New Task

- **Endpoint:** `POST /tasks`
- **Description:** Creates a new task.
- **Request Body:**

```json
{
  "title": "string",
  "description": "string",
  "due_date": "DD-MM-YYYY",
  "status": "Pending|In Progress|Completed"
}
```

- **Response:**
  - Status: 201 Created
  - Body: JSON task object with assigned ID
- **Error Responses:**
  - 400 Bad Request for invalid input or invalid status value
  - Example:

```json
{
  "message": "Invalid status value"
}
```

### Update a Task

- **Endpoint:** `PUT /tasks/{id}`
- **Description:** Updates fields of an existing task.
- **Request Body:** (Any subset of the following fields)

```json
{
  "title": "string",
  "description": "string",
  "due_date": "DD-MM-YYYY",
  "status": "Pending|In Progress|Completed"
}

```

- **Response:**
  - Status: 200 OK
  - Body: Updated JSON task object
- **Error Responses:**
  - 400 Bad Request for invalid input, invalid status, or no valid fields to update
  - 404 Not Found if task does not exist
  - Example:
  
```json
{
  "message": "No valid fields to update"
}
```

### Delete a Task

- **Endpoint:** `DELETE /tasks/{id}`
- **Description:** Deletes a task by its ID.
- **Response:**
  - Status: 204 No Content
  - Body:

```json
{
  "message": "deleted task successfully"
}
```

- **Error Response:**
  - 404 Not Found if task does not exist
  
```json
{
  "message": "task not found"
}
```

## Authentication

The Task Manager API currently does not require authentication.

## Rate and Usage Limits

The API does not currently enforce rate limits.

## Error Responses

- 400 Bad Request: The request body is invalid or missing required fields.
- 404 Not Found: The requested resource does not exist.
- 500 Internal Server Error: An unexpected error occurred on the server.
