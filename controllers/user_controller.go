package controllers

import (
	"BangkitcellBe/config"
	"BangkitcellBe/model"
	"net/http"
	"golang.org/x/crypto/bcrypt"

	"github.com/gin-gonic/gin"
)

// GET /users
func GetAllUser(c *gin.Context) {
	var users []model.User
	if err := config.DB.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error fetching users"})
		return
	}
	c.JSON(http.StatusOK, users)
}

// GET /users/:id
func GetUserById(c *gin.Context) {
	id := c.Param("id")

	var user model.User
	if err := config.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}

// POST /users
func CreateUser(c *gin.Context) {
	var input struct {
		Name     string `json:"name" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=8"`
		Role     string `json:"role" binding:"required,oneof=admin operator"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashed, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	user := model.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: string(hashed),
		Role:     input.Role,
	}

	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error creating user", "error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}

// PUT /users/:id
func UpdateUser(c *gin.Context) {
	id := c.Param("id")

	var user model.User
	if err := config.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
		return
	}

	var input struct {
		Name     *string `json:"name"`
		Email    *string `json:"email"`
		Password *string `json:"password"`
		Role     *string `json:"role"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.Name != nil {
		user.Name = *input.Name
	}
	if input.Email != nil {
		user.Email = *input.Email
	}
	if input.Password != nil {
		hashed, _ := bcrypt.GenerateFromPassword([]byte(*input.Password), bcrypt.DefaultCost)
		user.Password = string(hashed)
	}
	if input.Role != nil {
		user.Role = *input.Role
	}

	if err := config.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error updating user"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// DELETE /users/:id
func DeleteUser(c *gin.Context) {
	id := c.Param("id")

	var user model.User
	if err := config.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
		return
	}

	if err := config.DB.Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error deleting user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}