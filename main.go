package main

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func initDB() {
    var err error
    db, err = gorm.Open(sqlite.Open("lms.db"), &gorm.Config{})
    if err != nil {
        panic("failed to connect database")
    }

    db.AutoMigrate(&User{}, &Session{})
    _pass, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
    db.Create(&User{Username: "user", Password: string(_pass)})
}

func register(c *gin.Context) {
    var user User
    if err := c.ShouldBindJSON(&user); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
        return
    }

    tx := db.Create(&User{Username: user.Username, Password: string(hashedPassword)})
    if tx.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": tx.Error.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}

func login(c *gin.Context) {
    var user User
    if err := c.ShouldBindJSON(&user); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    var storedUser User
    // row := db.QueryRow("SELECT id, username, password FROM users WHERE username = ?", user.Username)
    // err := row.Scan(&storedUser.ID, &storedUser.Username, &storedUser.Password)
    // if err != nil {
    //     c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
    //     return
    // }
    tx := db.First(&storedUser, "username = ?", user.Username)
    if tx.Error != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": tx.Error.Error()})
        return
    }

    err := bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(user.Password))
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
        return
    }

    token := make([]byte, 32)
    _, err = rand.Read(token)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
        return
    }

    tokenStr := base64.StdEncoding.EncodeToString(token)
    
    // token = []byte("token")

    // _, err = db.Exec("INSERT INTO sessions (user_id, token) VALUES (?, ?)", storedUser.ID, token)
    // if err != nil {
    //     c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create session"})
    //     return
    // }
    tx = db.Create(&Session{UserID: storedUser.ID, Token: tokenStr})
    if tx.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": tx.Error.Error()})
        return
    }


    c.SetCookie("token", tokenStr, 3600, "", "", false, true)
    c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}

func ping(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{"message": "pong 2"})
}

func authMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Implement token-based authentication here
        // For simplicity, this example does not include token generation and validation

        token, err := c.Cookie("token")
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
            c.Abort()
            return
        }

        var session Session
        db.First(&session, "token = ?", token)
        if session.ID == 0 {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
            c.Abort()
            return
        }

        c.Next()
    }
}

func main() {
    initDB()

    r := gin.Default()

    r.GET("/ping", ping)

    r.POST("/signup", register)
    r.POST("/signin", login)

    auth := r.Group("/")
    auth.Use(authMiddleware())
    {
        auth.GET("/protected", func(c *gin.Context) {
            c.JSON(http.StatusOK, gin.H{"message": "This is a protected route"})
        })
    }

    r.Run()
}