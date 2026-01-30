package main

import (
	"Korifit/routes"

	"github.com/gin-gonic/gin"
)
func main() {
	router := gin.New()
	
	api := router.Group("/api")

	routes.ProcessRoutes(api)

	router.Run(":8080")
}