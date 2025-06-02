package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type Order struct {
	ID          int `json:"id"`
	UserID      int `json:"user_id"`
	ItemID      int `json:"item_id"`
	Quantity    int `json:"quantity"`
}

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Address  string `json:"address"`
}

type Item struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
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

	// Create orders table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS orders (
			id SERIAL PRIMARY KEY,
			user_id INTEGER NOT NULL,
			item_id INTEGER NOT NULL,
			quantity INTEGER NOT NULL
		)
	`)
	if err != nil {
		log.Fatal(err)
	}
}

func getUserDetails(userID int) (*User, error) {
	userServiceURL := os.Getenv("USER_SERVICE_URL")
	resp, err := http.Get(fmt.Sprintf("%s/users/%d", userServiceURL, userID))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get user details, status: %d", resp.StatusCode)
	}

	var user User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

func getItemDetails(itemID int) (*Item, error) {
	itemServiceURL := os.Getenv("ITEM_SERVICE_URL")
	resp, err := http.Get(fmt.Sprintf("%s/items/%d", itemServiceURL, itemID))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get item details, status: %d", resp.StatusCode)
	}

	var item Item
	if err := json.NewDecoder(resp.Body).Decode(&item); err != nil {
		return nil, err
	}
	return &item, nil
}

func createOrder(c *gin.Context) {
	var order Order
	if err := c.BindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Get user details
	user, err := getUserDetails(order.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user details"})
		return
	}

	// Get item details
	item, err := getItemDetails(order.ItemID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get item details"})
		return
	}

	// Calculate total amount (dummy calculation)
	totalAmount := float64(order.Quantity) * 10.0 // Assuming each item costs $10

	// Save order to database
	err = db.QueryRow(
		"INSERT INTO orders (user_id, item_id, quantity) VALUES ($1, $2, $3) RETURNING id",
		order.UserID, order.ItemID, order.Quantity,
	).Scan(&order.ID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
		return
	}

	// Return order details
	c.JSON(http.StatusCreated, gin.H{
		"payment_message": fmt.Sprintf("Proceeding to pay $%.2f total amount of order placed", totalAmount),
		"order_confirmation": "Order Successfully Placed. Thank you!",
		"order_details": gin.H{
			"order_id":      order.ID,
			"user_name":     user.Username,
			"item_name":     item.Name,
			"item_quantity": order.Quantity,
			"address":       user.Address,
		},
	})
}

func main() {
	initDB()
	r := gin.Default()

	r.POST("/orders", createOrder)

	log.Println("Order service starting on port 8083...")
	r.Run(":8083")
} 