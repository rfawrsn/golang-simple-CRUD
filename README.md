# Gin API Server

A simple REST API server built with Go and Gin framework that demonstrates all HTTP methods (GET, POST, PUT, DELETE).

## Features

- **GET** - Retrieve all users or a specific user by ID
- **POST** - Create a new user
- **PUT** - Update an existing user
- **DELETE** - Delete a user by ID
- CORS support for cross-origin requests
- Health check endpoint
- JSON request/response handling
- Error handling and validation

## Prerequisites

- Go 1.21 or higher
- Git (for cloning)

## Installation

1. Clone or download this project
2. Navigate to the project directory
3. Install dependencies:
   ```bash
   go mod tidy
   ```

## Running the Server

Start the server with:
```bash
go run main.go
```

The server will start on `http://localhost:3000`

## API Endpoints

### Health Check
- **GET** `/health` - Check if server is running

### Users API

#### Get All Users
- **GET** `/api/users`
- Returns all users in the system

#### Get User by ID
- **GET** `/api/users/:id`
- Returns a specific user by their ID
- Example: `GET /api/users/1`

#### Create User
- **POST** `/api/users`
- Creates a new user
- Request body:
  ```json
  {
    "name": "John Doe",
    "email": "john@example.com"
  }
  ```

#### Update User
- **PUT** `/api/users/:id`
- Updates an existing user
- Example: `PUT /api/users/1`
- Request body:
  ```json
  {
    "name": "John Updated",
    "email": "john.updated@example.com"
  }
  ```

#### Delete User
- **DELETE** `/api/users/:id`
- Deletes a user by ID
- Example: `DELETE /api/users/1`

## Testing the API

You can test the API using curl, Postman, or any HTTP client.

### Example curl commands:

```bash
# Health check
curl http://localhost:3000/health

# Get all users
curl http://localhost:3000/api/users

# Get user by ID
curl http://localhost:3000/api/users/1

# Create a new user
curl -X POST http://localhost:3000/api/users \
  -H "Content-Type: application/json" \
  -d '{"name": "Alice Johnson", "email": "alice@example.com"}'

# Update a user
curl -X PUT http://localhost:3000/api/users/1 \
  -H "Content-Type: application/json" \
  -d '{"name": "John Updated", "email": "john.updated@example.com"}'

# Delete a user
curl -X DELETE http://localhost:3000/api/users/1
```

## Response Format

All API responses follow a consistent format:

### Success Response
```json
{
  "status": "success",
  "data": {...},
  "message": "Optional message"
}
```

### Error Response
```json
{
  "status": "error",
  "message": "Error description"
}
```

## Project Structure

```
gin-api/
├── main.go      # Main application file
├── go.mod       # Go module file
└── README.md    # This file
```

## Notes

- This is a simple example using in-memory storage
- In a production environment, you would use a database
- The server includes CORS middleware for cross-origin requests
- Error handling includes validation for request bodies and user IDs
- Email uniqueness is enforced for user creation and updates 