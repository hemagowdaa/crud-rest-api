package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

const (
	DBUser     = "postgres"
	DBPassword = "hemagowda"
	DBName     = "postgres"
)

var db *sql.DB

type Userz struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}

func main() {
	// Establish a database connection
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", DBUser, DBPassword, DBName)
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}
	defer db.Close()

	// Initialize the Gin router
	router := gin.Default()

	// Define the API endpoints
	router.GET("/users", getUsers)
	router.GET("/users/:id", getUser)
	router.POST("/users", createUser)
	router.PUT("/users/:id", updateUser)
	router.DELETE("/users/:id", deleteUser)

	// Start the server
	log.Fatal(router.Run(":8000"))
}

func getUsers(c *gin.Context) {
	users := []Userz{}
	rows, err := db.Query("SELECT id, email FROM userz")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var user Userz
		if err := rows.Scan(&user.ID, &user.Email); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		users = append(users, user)
	}

	c.JSON(http.StatusOK, users)
}

func getUser(c *gin.Context) {
	id := c.Param("id")

	var user Userz
	row := db.QueryRow("SELECT id, email FROM userz WHERE id=$1", id)
	if err := row.Scan(&user.ID, &user.Email); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func createUser(c *gin.Context) {
	var user Userz
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// Insert the user into the database
	_, err := db.Exec("INSERT INTO userz (ID, email) VALUES ($1, $2)", user.ID, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}

func updateUser(c *gin.Context) {
	id := c.Param("id")

	var user Userz
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// Update the user in the database
	_, err := db.Exec("UPDATE userz SET email=$1 WHERE id=$2", user.Email, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

func deleteUser(c *gin.Context) {
	id := c.Param("id")

	// Delete the user from the database
	_, err := db.Exec("DELETE FROM userz WHERE id=$1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
