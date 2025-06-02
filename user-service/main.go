package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Address  string `json:"address"`
}

var db *sql.DB

func initDB() {
	var err error
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)
	
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	// Create users table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			username VARCHAR(50) UNIQUE NOT NULL,
			password VARCHAR(50) NOT NULL,
			address TEXT NOT NULL
		)
	`)
	if err != nil {
		log.Fatal(err)
	}

	// Insert sample users if they don't exist
	_, err = db.Exec(`
		INSERT INTO users (username, password, address)
		VALUES 
			('user1', 'pass1', '123 Main St, City1'),
			('user2', 'pass2', '456 Oak St, City2')
		ON CONFLICT (username) DO NOTHING
	`)
	if err != nil {
		log.Fatal(err)
	}
}

func login(c *gin.Context) {
	var credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.BindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	var user User
	err := db.QueryRow("SELECT id, username, password, address FROM users WHERE username = $1 AND password = $2",
		credentials.Username, credentials.Password).Scan(&user.ID, &user.Username, &user.Password, &user.Address)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"address":  user.Address,
		},
	})
}

func getUserById(c *gin.Context) {
	userID := c.Param("id")
	var user User
	err := db.QueryRow("SELECT id, username, address FROM users WHERE id = $1", userID).
		Scan(&user.ID, &user.Username, &user.Address)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func main() {
	initDB()
	r := gin.Default()

	r.POST("/login", login)
	r.GET("/users/:id", getUserById)

	log.Println("User service starting on port 8081...")
	r.Run(":8081")
} 