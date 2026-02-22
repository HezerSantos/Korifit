package routes

import (
	"Korifit/controllers"
	"Korifit/middleware"

	"github.com/gin-gonic/gin"
)

func ProcessRoutes(api *gin.RouterGroup) {
	api.GET("csrf", controllers.GetCsrfToken)

	users := api.Group("/users")

	users.POST("/signup", controllers.CreateUser)
	users.POST("/login", controllers.VerifyUser)

	fitness := api.Group("/fitness")
	fitness.Use(middleware.AuthenticateUser)
	fitness.GET("/exercises", controllers.GetExercises)
	fitness.POST("/exercises", controllers.CreateExercise)
	fitness.GET("/exercises/:id", controllers.GetExerciseByID)

	fitness.GET("/workouts", controllers.GetWorkouts)
	fitness.POST("/workouts", controllers.CreateWorkout)
	fitness.GET("/workouts/:id", controllers.GetWorkoutByID)
	fitness.GET("/nutrition", controllers.GetNutritionList)
	fitness.GET("/nutrition/:id", controllers.GetNutritionListByID)

	recipes := api.Group("/recipes")
	recipes.Use(middleware.AuthenticateCsrf)
	recipes.GET("/", controllers.GetRecipes)
}
