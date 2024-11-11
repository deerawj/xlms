package handlers

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"time"
	"xlms/shared"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

var (
    accessExpiresIn  = 3600
    refreshExpiresIn = 86400
)

func getRndHash() string {
	// Create a random byte slice
	randomData := make([]byte, 32) // 32 bytes of random data (256 bits)
	io.ReadFull(rand.Reader, randomData)
	hashStr := fmt.Sprintf("%x", sha256.Sum256(randomData))

	return hashStr[:32]
}

func Login(c *gin.Context) {
	var user shared.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var storedUser shared.User

	tx := shared.DB.First(&storedUser, "username = ?", user.Username)
	if tx.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": tx.Error.Error()})
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(user.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	accessToken := "access_" + getRndHash()
	refreshToken := "refresh_" + getRndHash()

	tx = shared.DB.Create(&shared.Session{UserID: storedUser.ID, AccessToken: accessToken, RefreshToken: refreshToken})
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": tx.Error.Error()})
		return
	}

	// TODO: set token to secure in production
	c.JSON(http.StatusOK, gin.H{"message": "Login successful", "access_token": accessToken, "refresh_token": refreshToken, "access_expires_in": accessExpiresIn, "refresh_expires_in": refreshExpiresIn})
}

func Logout(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var session shared.Session
	tx := shared.DB.First(&session, "access_token = ?", token)
	if tx.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	tx = shared.DB.Unscoped().Delete(&session)
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": tx.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
}

func RefreshxToken(c *gin.Context) {
	// get refresh token from request json
	var postForm struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := c.ShouldBindJSON(&postForm); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var session shared.Session
	tx := shared.DB.First(&session, "refresh_token = ?", postForm.RefreshToken)
	if tx.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	if session.CreatedAt.Add(time.Duration(refreshExpiresIn) * time.Second).Before(time.Now()) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token expired"})
		shared.DB.Unscoped().Delete(&session)
		return
	}

	accessToken := "access_" + getRndHash()

	tx = shared.DB.Model(&session).Update("access_token", accessToken)
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": tx.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Token refreshed", "access_token": accessToken, "refresh_token": session.RefreshToken, "access_expires_in": accessExpiresIn, "refresh_expires_in": refreshExpiresIn})
}
