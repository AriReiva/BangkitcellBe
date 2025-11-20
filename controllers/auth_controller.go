package controllers

import (
	"BangkitcellBe/config"
	"BangkitcellBe/model"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var JWT_SECRET = []byte(os.Getenv("JWT_SECRET"))

// POST /register
func RegisterUser(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
		Name     string `json:"name" binding:"required,min=6"`
		Role     string `json:"role"` // optional default operator
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": err.Error()})
		return
	}

	// Cek email duplicate
	var existing model.User
	err := config.DB.Where("email = ?", input.Email).First(&existing).Error
	if err == nil && existing.ID != 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Email already exists",
		})
		return
	}

	if input.Role == "" {
		input.Role = "operator"
	}

	if input.Role != "admin" && input.Role != "operator" {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": "Role must be admin or operator"})
		return
	}

	// Hash password
	hashed, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)

	user := model.User{
		Email:    input.Email,
		Password: string(hashed),
		Name:     input.Name,
		Role:     input.Role,
	}

	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "Error saving user"})
		return
	}

	user.Password = "" // hide password

	c.JSON(http.StatusCreated, gin.H{
		"status":  true,
		"message": "Registrasi berhasil",
		"data":    user,
	})
}

// POST /login
func LoginUser(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": err.Error()})
		return
	}

	var user model.User
	if err := config.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"status": false, "message": "Email atau password salah"})
		return
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)) != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"status": false, "message": "Email atau password salah"})
		return
	}

	// Generate JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    user.ID,
		"email": user.Email,
		"role":  user.Role,
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString(JWT_SECRET)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "Failed to generate token"})
		return
	}

	user.Password = ""

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "Login berhasil",
		"token":   tokenString,
		"user":    user,
	})
}

func LogoutUser(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "Logout berhasil (hapus token dari client)",
	})
}
