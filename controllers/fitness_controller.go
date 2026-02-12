package controllers

import (
	"Korifit/config"
	"Korifit/helpers"
	"github.com/gin-gonic/gin"
)

func GetExercises(c *gin.Context) {
	id, _ := c.Get("userId")
	var exercies []config.Exercise
	result := config.DB.Find(&exercies)

	if result.Error != nil {
		helpers.NetworkError(c, result.Error)
		return
	}

	c.JSON(200, gin.H{
		"exercises": exercies,
		"id": id,
	})
}