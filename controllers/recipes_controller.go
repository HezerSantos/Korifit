package controllers

import "github.com/gin-gonic/gin"

func GetRecipes(c *gin.Context) {
	c.JSON(200, gin.H{"msg": "good"})
}
