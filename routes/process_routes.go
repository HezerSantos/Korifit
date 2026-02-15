package routes

import (
	"Korifit/controllers"
	// "Korifit/middleware"

	"github.com/gin-gonic/gin"
)

func ProcessRoutes(api *gin.RouterGroup) {
	
	users := api.Group("/users")

	users.POST("/signup", controllers.CreateUser)
	users.POST("/login", controllers.VerifyUser)


	fitness := api.Group("/fitness")
	// fitness.Use(middleware.AuthenticateUser)
	fitness.GET("/exercises", controllers.GetExercises)
	fitness.POST("/exercises", controllers.CreateExercise)
	fitness.GET("/exercises/:id", controllers.GetExerciseByID)


	fitness.GET("/workouts", controllers.GetWorkouts)
	// fitness.GET("/workouts/:id")
	// fitness.GET("/nutrition")




	// recipes := api.Group("/recipes")
	// recipes.GET("/")
}