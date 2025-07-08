package main

import (
	"net/http"
	"strconv"

	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("mysecretkey") // Simple secret key for JWT

// User represents a user in our system
type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"` // "admin" or "user"
}

// In-memory storage for users (not hardcoded, can register new users)
var users = []User{}
var nextUserID = 1

// JWT Claims

type Claims struct {
	UserID int    `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func main() {
	// Create a new Gin router with default middleware
	r := gin.Default()

	// Add CORS middleware
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "Server is running",
		})
	})

	// Auth endpoints
	r.POST("/api/register", registerUser)
	r.POST("/api/login", loginUser)

	// API routes group
	api := r.Group("/api")
	api.Use(authMiddleware())
	{
		// GET /api/users - Get all users
		api.GET("/users", getUsers)

		// GET /api/users/:id - Get user by ID
		api.GET("/users/:id", getUserByID)

		// POST /api/users - Create a new user
		api.POST("/users", adminOnly(createUser))

		// PUT /api/users/:id - Update user by ID
		api.PUT("/users/:id", adminOnly(updateUser))

		// DELETE /api/users/:id - Delete user by ID
		api.DELETE("/users/:id", adminOnly(deleteUser))
	}

	// Start the server on port 3000
	r.Run(":3000")
}

// Registration endpoint
func registerUser(c *gin.Context) {
	var req struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid request body"})
		return
	}
	if req.Name == "" || req.Email == "" || req.Password == "" || req.Role == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "All fields are required"})
		return
	}
	if req.Role != "admin" && req.Role != "user" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Role must be 'admin' or 'user'"})
		return
	}
	for _, user := range users {
		if user.Email == req.Email {
			c.JSON(http.StatusConflict, gin.H{"status": "error", "message": "User with this email already exists"})
			return
		}
	}
	user := User{
		ID:       nextUserID,
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password, 
		Role:     req.Role,
	}
	nextUserID++
	users = append(users, user)
	c.JSON(http.StatusCreated, gin.H{"status": "success", "data": user, "message": "User registered successfully"})
}

func loginUser(c *gin.Context) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid request body"})
		return
	}
	for _, user := range users {
		if user.Email == req.Email && user.Password == req.Password {
			// Generate JWT
			claims := Claims{
				UserID: user.ID,
				Email:  user.Email,
				Role:   user.Role,
				RegisteredClaims: jwt.RegisteredClaims{
					ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
				},
			}
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
			tokenString, err := token.SignedString(jwtSecret)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to generate token"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"status": "success", "token": tokenString})
			return
		}
	}
	c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Invalid email or password"})
}

// JWT Auth Middleware
func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Missing or invalid Authorization header"})
			c.Abort()
			return
		}
		tokenString := strings.TrimPrefix(header, "Bearer ")
		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Invalid or expired token"})
			c.Abort()
			return
		}
		claims, ok := token.Claims.(*Claims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Invalid token claims"})
			c.Abort()
			return
		}
		c.Set("userID", claims.UserID)
		c.Set("userRole", claims.Role)
		c.Next()
	}
}

// Admin-only middleware
func adminOnly(handler gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("userRole")
		if !exists || role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"status": "error", "message": "Admin access required"})
			return
		}
		handler(c)
	}
}

// getUsers returns all users
func getUsers(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   users,
		"count":  len(users),
	})
}

// getUserByID returns a specific user by ID
func getUserByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid user ID",
		})
		return
	}

	for _, user := range users {
		if user.ID == id {
			c.JSON(http.StatusOK, gin.H{
				"status": "success",
				"data":   user,
			})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{
		"status":  "error",
		"message": "User not found",
	})
}

// createUser creates a new user
func createUser(c *gin.Context) {
	var newUser User

	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid request body",
		})
		return
	}

	// Generate a new ID (in a real app, this would be handled by the database)
	newUser.ID = len(users) + 1

	// Check if user with same email already exists
	for _, user := range users {
		if user.Email == newUser.Email {
			c.JSON(http.StatusConflict, gin.H{
				"status":  "error",
				"message": "User with this email already exists",
			})
			return
		}
	}

	users = append(users, newUser)

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"data":    newUser,
		"message": "User created successfully",
	})
}

// updateUser updates an existing user
func updateUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid user ID",
		})
		return
	}

	var updateData User
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid request body",
		})
		return
	}

	// Find and update the user
	for i, user := range users {
		if user.ID == id {
			// Update fields if provided
			if updateData.Name != "" {
				users[i].Name = updateData.Name
			}
			if updateData.Email != "" {
				// Check if email is already taken by another user
				for j, otherUser := range users {
					if j != i && otherUser.Email == updateData.Email {
						c.JSON(http.StatusConflict, gin.H{
							"status":  "error",
							"message": "Email already taken by another user",
						})
						return
					}
				}
				users[i].Email = updateData.Email
			}

			c.JSON(http.StatusOK, gin.H{
				"status":  "success",
				"data":    users[i],
				"message": "User updated successfully",
			})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{
		"status":  "error",
		"message": "User not found",
	})
}

// deleteUser deletes a user by ID
func deleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid user ID",
		})
		return
	}

	// Find and remove the user
	for i, user := range users {
		if user.ID == id {
			users = append(users[:i], users[i+1:]...)
			c.JSON(http.StatusOK, gin.H{
				"status":  "success",
				"message": "User deleted successfully",
			})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{
		"status":  "error",
		"message": "User not found",
	})
}
