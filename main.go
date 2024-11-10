package main

import (
    "database/sql"
    "log"
    "net/http"

    "github.com/gin-gonic/gin"
    "golang.org/x/crypto/bcrypt"
    _ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func initDB() {
    var err error
    db, err = sql.Open("sqlite3", "./lms.db")
    if err != nil {
        log.Fatal(err)
    }

    createTable := `
    CREATE TABLE IF NOT EXISTS users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        username TEXT NOT NULL UNIQUE,
        password TEXT NOT NULL,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP
    );
    `
    _, err = db.Exec(createTable)
    if err != nil {
        log.Fatal(err)
    }
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

    _, err = db.Exec("INSERT INTO users (username, password) VALUES (?, ?)", user.Username, string(hashedPassword))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
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
    row := db.QueryRow("SELECT id, username, password FROM users WHERE username = ?", user.Username)
    err := row.Scan(&storedUser.ID, &storedUser.Username, &storedUser.Password)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
        return
    }

    err = bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(user.Password))
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}

func authMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Implement token-based authentication here
        // For simplicity, this example does not include token generation and validation
        c.Next()
    }
}

func main() {
    initDB()

    r := gin.Default()

    r.POST("/register", register)
    r.POST("/login", login)

    auth := r.Group("/")
    auth.Use(authMiddleware())
    {
        auth.GET("/protected", func(c *gin.Context) {
            c.JSON(http.StatusOK, gin.H{"message": "This is a protected route"})
        })
    }

    r.Run()
}