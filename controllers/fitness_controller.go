package controllers

import (
	"Korifit/config"

	"github.com/gin-gonic/gin"
)

func GetExercises(c *gin.Context) {
	allExercises := config.DB.Find(&config.Exercise{})

	c.JSON(200, gin.H{
		"exercises": allExercises,
	})
}