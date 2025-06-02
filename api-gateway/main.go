package main

import (
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func forwardRequest(c *gin.Context, targetURL string) {
	// 1. Create a new request to forward to the target service
	req, err := http.NewRequest(c.Request.Method, targetURL, c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}

	// 2. Copy all headers from original request to forwarded request
	for key, values := range c.Request.Header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	// 3. Send the request to target service
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Service unavailable"})
		return
	}
	defer resp.Body.Close()

	// 4. Copy response headers from service back to client
	for key, values := range resp.Header {
		for _, value := range values {
			c.Header(key, value)
		}
	}

	// 5. Read response body from service
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response"})
		return
	}

	// 6. Forward response status code and body back to client
	c.Data(resp.StatusCode, resp.Header.Get("Content-Type"), body)
}

func main() {
	r := gin.Default()

	userServiceURL := "http://user-service:8081"
	itemServiceURL := "http://item-service:8082"
	orderServiceURL := "http://order-service:8083"

	// User Service Routes
	r.POST("/api/login", func(c *gin.Context) {
		forwardRequest(c, userServiceURL+"/login")
	})
	r.GET("/api/users/:id", func(c *gin.Context) {
		forwardRequest(c, userServiceURL+"/users/"+c.Param("id"))
	})

	// Item Service Routes
	r.POST("/api/items", func(c *gin.Context) {
		forwardRequest(c, itemServiceURL+"/items")
	})
	r.GET("/api/items/:id", func(c *gin.Context) {
		forwardRequest(c, itemServiceURL+"/items/"+c.Param("id"))
	})

	// Order Service Routes
	r.POST("/api/orders", func(c *gin.Context) {
		forwardRequest(c, orderServiceURL+"/orders")
	})

	log.Println("API Gateway starting on port 8080...")
	r.Run(":8080")
} 