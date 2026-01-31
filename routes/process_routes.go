package routes

import (
	"Korifit/controllers"

	"github.com/gin-gonic/gin"
)

func ProcessRoutes(api *gin.RouterGroup) {
	
	users := api.Group("/users")

	users.POST("/signup", controllers.CreateUser)
	users.POST("/login", controllers.VerifyUser)


	fitness := api.Group("/fitness")
	fitness.GET("/exercises", controllers.GetExercises)
	// fitness.GET("/exercises/:id")
	// fitness.GET("/workouts")
	// fitness.GET("/workouts/:id")
	// fitness.GET("/nutrition")




	// recipes := api.Group("/recipes")
	// recipes.GET("/")
}