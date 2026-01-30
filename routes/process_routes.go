package routes
import (
	"github.com/gin-gonic/gin"
)

func ProcessRoutes(api *gin.RouterGroup) {
	
	users := api.Group("/users")

	users.GET("/")
	users.GET("/:id")


	fitness := api.Group("/fitness")
	fitness.GET("/exercises")
	fitness.GET("/exercises/:id")
	fitness.GET("/workouts")
	fitness.GET("/workouts/:id")
	fitness.GET("/nutrition")




	recipes := api.Group("/recipes")
	recipes.GET("/")
}