package routes

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

var Users = []User{}
var NextUserID = 1

// getUsers returns all users
func GetUsers(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   Users,
		"count":  len(Users),
	})
}

// getUserByID returns a specific user by ID
func GetUserByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid user ID"})
		return
	}
	for _, user := range Users {
		if user.ID == id {
			c.JSON(http.StatusOK, gin.H{"status": "success", "data": user})
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "User not found"})
}

// createUser creates a new user
func CreateUser(c *gin.Context) {
	var newUser User
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid request body"})
		return
	}
	newUser.ID = NextUserID
	NextUserID++
	for _, user := range Users {
		if user.Email == newUser.Email {
			c.JSON(http.StatusConflict, gin.H{"status": "error", "message": "User with this email already exists"})
			return
		}
	}
	Users = append(Users, newUser)
	c.JSON(http.StatusCreated, gin.H{"status": "success", "data": newUser, "message": "User created successfully"})
}

// updateUser updates an existing user
func UpdateUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid user ID"})
		return
	}
	var updateData User
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid request body"})
		return
	}
	for i, user := range Users {
		if user.ID == id {
			if updateData.Name != "" {
				Users[i].Name = updateData.Name
			}
			if updateData.Email != "" {
				for j, otherUser := range Users {
					if j != i && otherUser.Email == updateData.Email {
						c.JSON(http.StatusConflict, gin.H{"status": "error", "message": "Email already taken by another user"})
						return
					}
				}
				Users[i].Email = updateData.Email
			}
			c.JSON(http.StatusOK, gin.H{"status": "success", "data": Users[i], "message": "User updated successfully"})
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "User not found"})
}

// deleteUser deletes a user by ID
func DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid user ID"})
		return
	}
	for i, user := range Users {
		if user.ID == id {
			Users = append(Users[:i], Users[i+1:]...)
			c.JSON(http.StatusOK, gin.H{"status": "success", "message": "User deleted successfully"})
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "User not found"})
}
