package main

import (
	"net/http"
	"time"
	"xlms/shared"
    "xlms/handlers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "pong"})
}

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		var session shared.Session
		tx := shared.DB.First(&session, "access_token = ?", token)
		if tx.Error != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
		if session.UpdatedAt.Add(time.Duration(3600) * time.Second).Before(time.Now()) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token expired"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func main() {
	shared.InitDB()

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * 60 * 60,
	}))

	r.GET("/ping", ping)

	r.POST("/signup", handlers.Register)
	r.POST("/signin", handlers.Login)
    r.POST("/signout", handlers.Logout)
    r.POST("/refresh", handlers.RefreshxToken)

	auth := r.Group("/")
	auth.Use(authMiddleware())
    {
        auth.GET("/user", handlers.GetUserInfo)
    }
	


	r.Run()
}
