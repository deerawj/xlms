package handlers

import (
	"xlms/shared"

	"github.com/gin-gonic/gin"
)

func GetUserInfo(c *gin.Context) {
	token := c.GetHeader("Authorization")
	var user shared.User
	var session shared.Session

	shared.DB.First(&session, "access_token = ?", token)
	shared.DB.First(&user, "id = ?", session.UserID)

	c.JSON(200, gin.H{"username": user.Username, "id": user.ID, "created_at": user.CreatedAt})
}