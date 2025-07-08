package routes

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var JwtSecret = []byte("mysecretkey")

type Claims struct {
	UserID int    `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// RegisterUser handles user registration
func RegisterUser(c *gin.Context) {
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
	for _, user := range Users {
		if user.Email == req.Email {
			c.JSON(http.StatusConflict, gin.H{"status": "error", "message": "User with this email already exists"})
			return
		}
	}
	user := User{
		ID:       NextUserID,
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		Role:     req.Role,
	}
	NextUserID++
	Users = append(Users, user)
	c.JSON(http.StatusCreated, gin.H{"status": "success", "data": user, "message": "User registered successfully"})
}

// LoginUser handles user login and JWT generation
func LoginUser(c *gin.Context) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid request body"})
		return
	}
	for _, user := range Users {
		if user.Email == req.Email && user.Password == req.Password {
			claims := Claims{
				UserID: user.ID,
				Email:  user.Email,
				Role:   user.Role,
				RegisteredClaims: jwt.RegisteredClaims{
					ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
				},
			}
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
			tokenString, err := token.SignedString(JwtSecret)
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
