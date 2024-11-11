package handlers

import (
	"net/http"
	"xlms/shared"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *gin.Context) {
	var user shared.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	tx := shared.DB.Create(&shared.User{Username: user.Username, Password: string(hashedPassword)})
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": tx.Error.Error()})
		return
	}
	shared.DB.Find(&user, "username = ?", user.Username)

	accessToken := getRndHash()
	refreshToken := getRndHash()
	accessExpiresIn := 3600
	refreshExpiresIn := 86400

	tx = shared.DB.Create(&shared.Session{UserID: user.ID, AccessToken: accessToken, RefreshToken: refreshToken})
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": tx.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully", "access_token": accessToken, "refresh_token": refreshToken, "access_expires_in": accessExpiresIn, "refresh_expires_in": refreshExpiresIn})
}