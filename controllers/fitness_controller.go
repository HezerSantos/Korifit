package controllers

import (
	"Korifit/config"
	"Korifit/helpers"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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


type CreateExerciseJSON struct{
	Name string `json:"name" binding:"required"`
	MuscleTarget string `json:"muscleTarget" binding:"required"`
}

func CreateExercise(c *gin.Context) {
	var newExerciseJSON CreateExerciseJSON

	err := c.ShouldBind(&newExerciseJSON)

	if err != nil {
		helpers.ErrorHelper(
			c, 
			helpers.JsonError{
				Message: "JSON ERROR 001", 
				Status: 400, 
				Json: helpers.JsonResponseType{Code: "INVALID_BODY", Msg: "JSON ERROR 001"},
			},
		)
		return
	}

	exercise := config.Exercise{Name: newExerciseJSON.Name, MuscleTarget: newExerciseJSON.MuscleTarget}
	config.DB.Create(&exercise)

	c.JSON(200, gin.H{
		"msg": "Record successfully created",
		"record": map[string]interface{}{
			"id": exercise.ID,
			"name": exercise.Name,
			"muscleTarget": exercise.MuscleTarget,
			"workouts": exercise.Workouts,
		},
	})
}



func GetExerciseByID(c *gin.Context) {
	paramId := c.Param("id")
	id, err := uuid.Parse(paramId)

	if err != nil {
		helpers.ErrorHelper(c, 
			helpers.JsonError{
				Message: "GetExerciseByID: ID ERROR",
				Status: 400,
				Json: helpers.JsonResponseType{Code: "INVALID_BODY", Msg: "Bad Request"},
			},
		)
		return
	}



	exercise := config.Exercise{ID: id}

	config.DB.Find(&exercise)

	c.JSON(200, gin.H{
		"msg": "Exercise Successfully Retrieved",
		"id": exercise.ID,
		"name": exercise.Name,
		"muscleTarget": exercise.MuscleTarget,
		"workouts": exercise.Workouts,
	})
}


func GetWorkouts(c *gin.Context) {
	userId, exists := c.Get("userId")
	
	if !exists {
		helpers.NetworkError(c, nil)
		return
	};

	parsedUuid, err := uuid.Parse(fmt.Sprint(userId))

	if err != nil {
		helpers.NetworkError(c, err)
		return
	}

	userWorkouts := config.Workout{UserID: parsedUuid}

	result := config.DB.Find(&userWorkouts)

	if result.Error != nil {
		helpers.NetworkError(c, err)
		return
	}

	c.JSON(200, gin.H{
		"msg": "Successfully Retrieved User Workouts",
		"id": parsedUuid,
		"workouts": userWorkouts,
	})
}

