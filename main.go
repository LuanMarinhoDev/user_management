package main

import (
	"teste_shipay/backend-challenge/database"
	"teste_shipay/backend-challenge/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	database.Connect()

	router := gin.Default()

	router.POST("/users", handlers.CreateUser)
	router.GET("/users/:id", handlers.GetUserById)

	router.Run(":8080")
}
