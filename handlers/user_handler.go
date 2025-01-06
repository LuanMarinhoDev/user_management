package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"net/http"
	"strconv"
	"time"

	"teste_shipay/backend-challenge/database"
	"teste_shipay/backend-challenge/models"

	"github.com/gin-gonic/gin"
)

func generateRandomPassword() string {
	bytes := make([]byte, 8)
	_, _ = rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func CreateUser(c *gin.Context) {
	var userInput struct {
		Name     string `json:"name" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password"`
		RoleID   uint   `json:"role_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&userInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if userInput.Password == "" {
		userInput.Password = generateRandomPassword()
	}

	user := models.User{
		Name:      userInput.Name,
		Email:     userInput.Email,
		Password:  userInput.Password,
		RoleID:    userInput.RoleID,
		CreatedAt: time.Now().String(),
		UpdatedAt: time.Now().String(),
	}

	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully", "user": user})
}

func GetUserById(c *gin.Context) {
	idParam := c.Param("id")

	log.Println(idParam)

	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	type UserDetails struct {
		UserName         string `json:"user_name"`
		UserEmail        string `json:"user_email"`
		RoleDescription  string `json:"role_description"`
		ClaimDescription string `json:"claim_description"`
	}

	var userDetails []UserDetails

	err = database.DB.Table("users").
		Select("users.name AS user_name, users.email AS user_email, roles.description AS role_description, claims.description AS claim_description").
		Joins("INNER JOIN roles ON users.role_id = roles.id").
		Joins("INNER JOIN user_claims ON users.id = user_claims.user_id").
		Joins("INNER JOIN claims ON user_claims.claim_id = claims.id").
		Where("users.id = ?", id).
		Scan(&userDetails).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching user details"})
		return
	}

	if len(userDetails) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user_details": userDetails,
	})
}
