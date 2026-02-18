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
		"id":        id,
	})
}

type CreateExerciseJSON struct {
	Name         string `json:"name" binding:"required"`
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
				Status:  400,
				Json:    helpers.JsonResponseType{Code: "INVALID_BODY", Msg: "JSON ERROR 001"},
			},
		)
		return
	}

	exercise := config.Exercise{Name: newExerciseJSON.Name, MuscleTarget: newExerciseJSON.MuscleTarget}
	config.DB.Create(&exercise)

	c.JSON(200, gin.H{
		"msg": "Record successfully created",
		"record": map[string]interface{}{
			"id":           exercise.ID,
			"name":         exercise.Name,
			"muscleTarget": exercise.MuscleTarget,
			"workouts":     exercise.Workouts,
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
				Status:  400,
				Json:    helpers.JsonResponseType{Code: "INVALID_BODY", Msg: "Bad Request"},
			},
		)
		return
	}

	exercise := config.Exercise{ID: id}

	config.DB.Find(&exercise)

	c.JSON(200, gin.H{
		"msg":          "Exercise Successfully Retrieved",
		"id":           exercise.ID,
		"name":         exercise.Name,
		"muscleTarget": exercise.MuscleTarget,
		"workouts":     exercise.Workouts,
	})
}

func GetWorkouts(c *gin.Context) {
	userId, exists := c.Get("userId")

	if !exists {
		helpers.NetworkError(c, nil)
		return
	}

	parsedUuid, err := uuid.Parse(fmt.Sprint(userId))

	if err != nil {
		helpers.NetworkError(c, err)
		return
	}

	userWorkouts := config.Workout{UserID: parsedUuid}

	result := config.DB.Preload("Exercises").Find(&userWorkouts)

	if result.Error != nil {
		helpers.NetworkError(c, err)
		return
	}

	c.JSON(200, gin.H{
		"msg":      "Successfully Retrieved User Workouts",
		"id":       parsedUuid,
		"workouts": userWorkouts,
	})
}

type CreateWorkoutJSON struct {
	Name      string   `json:"name" binding:"required"`
	Exercises []string `json:"exercises" binding:"required"`
}

func CreateWorkout(c *gin.Context) {
	var createWorkoutJson CreateWorkoutJSON

	err := c.ShouldBind(&createWorkoutJson)

	if err != nil {
		fmt.Println(err)
		helpers.ErrorHelper(
			c,
			helpers.JsonError{
				Message: "JSON ERROR 001",
				Status:  400,
				Json:    helpers.JsonResponseType{Code: "INVALID_BODY", Msg: "JSON ERROR 001"},
			},
		)
		return
	}

	userId, exists := c.Get("userId")

	if !exists {
		helpers.NetworkError(c, nil)
		return
	}

	parsedUuid, err := uuid.Parse(fmt.Sprint(userId))

	if err != nil {
		helpers.NetworkError(c, err)
		return
	}

	var targetExercises []config.Exercise

	for _, id := range createWorkoutJson.Exercises {
		parsedUuid, err := uuid.Parse(fmt.Sprint(id))

		if err != nil {
			helpers.NetworkError(c, err)
			return
		}

		selectedExercise := config.Exercise{ID: parsedUuid}
		result := config.DB.Find(&selectedExercise)

		if result.Error != nil {
			helpers.NetworkError(c, err)
			return
		}
		targetExercises = append(targetExercises, selectedExercise)
	}

	newWorkout := config.Workout{Name: createWorkoutJson.Name, Exercises: targetExercises, UserID: parsedUuid}

	result := config.DB.Create(&newWorkout)

	if result.Error != nil {
		helpers.NetworkError(c, result.Error)
		return
	}

	c.JSON(200, gin.H{
		"msg":       "Workout Created",
		"id":        newWorkout.ID,
		"userId":    newWorkout.UserID,
		"name":      newWorkout.Name,
		"exercises": newWorkout.Exercises,
	})
}

func GetWorkoutByID(c *gin.Context) {
	workoutId := c.Param("id")

	parsedWorkoutId, err := uuid.Parse(workoutId)

	if err != nil {
		helpers.ErrorHelper(c,
			helpers.JsonError{
				Message: "GetWorkoutByID: ID ERROR",
				Status:  400,
				Json:    helpers.JsonResponseType{Code: "INVALID_PARAM", Msg: "Bad Request"},
			},
		)
		return
	}

	userId, exists := c.Get("userId")

	if !exists {
		helpers.NetworkError(c, nil)
		return
	}

	parsedUserId, err := uuid.Parse(fmt.Sprint(userId))

	if err != nil {
		helpers.NetworkError(c, err)
		return
	}

	queriedWorkout := config.Workout{ID: parsedWorkoutId}
	result := config.DB.Preload("Exercises").Find(&queriedWorkout)

	if result.Error != nil {
		helpers.NetworkError(c, result.Error)
		return
	}

	if result.RowsAffected == 0 {
		helpers.ErrorHelper(c, helpers.JsonError{
			Message: "Workout not found",
			Status:  404,
			Json: helpers.JsonResponseType{
				Msg:  "INVALID_PARAM",
				Code: "INVALID_PARAM",
			},
		})
		return
	}

	if queriedWorkout.UserID != parsedUserId {
		helpers.ErrorHelper(c, helpers.JsonError{
			Message: "Unauthorized",
			Status:  401,
			Json: helpers.JsonResponseType{
				Msg:  "INVALID_ACCESS",
				Code: "INVALID_ACCESS",
			},
		})
		return
	}
	c.JSON(200, gin.H{
		"msg":     "Workout Retrieved",
		"workout": queriedWorkout,
	})

}

func GetNutritionList(c *gin.Context) {
	userId, exists := c.Get("userId")

	if !exists {
		helpers.ErrorHelper(c, helpers.JsonError{
			Message: "Unauthorized",
			Status:  401,
			Json: helpers.JsonResponseType{
				Msg:  "INVALID_ACCESS",
				Code: "INVALID_ACCESS",
			},
		})
		return
	}

	parsedUserId, err := uuid.Parse(fmt.Sprint(userId))

	if err != nil {
		helpers.NetworkError(c, err)
		return
	}

	var userNutritionList config.DailyNutritionList

	result := config.DB.Where("id = ?", parsedUserId).Preload("NutritionItems").Find(&userNutritionList)

	if result.Error != nil {
		helpers.NetworkError(c, err)
		return
	}

	if result.RowsAffected == 0 {
		userNutritionList := []interface{}{}
		c.JSON(200, gin.H{
			"msg":                "Retrieved Daily Nutrition List",
			"dailyNutritionList": userNutritionList,
		})
		return
	}

	c.JSON(200, gin.H{
		"msg":                "Retrieved Daily Nutrition List",
		"dailyNutritionList": userNutritionList,
	})
}

func GetNutritionListByID(c *gin.Context) {
	listId := c.Param("id")

	parsedListId, err := uuid.Parse(listId)

	if err != nil {
		helpers.ErrorHelper(c,
			helpers.JsonError{
				Message: "GetNutritionListByID: ID ERROR",
				Status:  400,
				Json:    helpers.JsonResponseType{Code: "INVALID_PARAM", Msg: "Bad Request"},
			},
		)
		return
	}
	userId, exists := c.Get("userId")

	if !exists {
		helpers.ErrorHelper(c, helpers.JsonError{
			Message: "Unauthorized",
			Status:  401,
			Json: helpers.JsonResponseType{
				Msg:  "INVALID_ACCESS",
				Code: "INVALID_ACCESS",
			},
		})
		return
	}

	parsedUserId, err := uuid.Parse(fmt.Sprint(userId))

	if err != nil {
		helpers.ErrorHelper(c, helpers.JsonError{
			Message: "Unauthorized",
			Status:  401,
			Json: helpers.JsonResponseType{
				Msg:  "INVALID_ACCESS",
				Code: "INVALID_ACCESS",
			},
		})
		return
	}

	var nutritionList config.DailyNutritionList

	result := config.DB.Preload("NutritionItems").Where("id = ?", parsedListId).Find(&nutritionList)

	if result.Error != nil {
		helpers.NetworkError(c, result.Error)
		return
	}

	if result.RowsAffected == 0 {
		helpers.ErrorHelper(c, helpers.JsonError{
			Message: "Nutrition List not found",
			Status:  404,
			Json: helpers.JsonResponseType{
				Msg:  "INVALID_PARAM",
				Code: "INVALID_PARAM",
			},
		})
		return
	}

	if nutritionList.UserID != parsedUserId {
		helpers.ErrorHelper(c, helpers.JsonError{
			Message: "Unauthorized",
			Status:  401,
			Json: helpers.JsonResponseType{
				Msg:  "INVALID_ACCESS",
				Code: "INVALID_ACCESS",
			},
		})
		return
	}

	c.JSON(200, gin.H{
		"msg":            "Retrieved Daily Nutrition List",
		"id":             nutritionList.ID,
		"date":           nutritionList.Date,
		"nutritionItems": nutritionList.NutritionItems,
		"userId":         nutritionList.UserID,
	})
}

type TemporaryNutritionItem struct {
	Name     string
	Calories int
	Protein  int
}
type NutritionListJson struct {
	Date          string                   `json:"date" binding:"required"`
	NutritionList []TemporaryNutritionItem `json:"nutritionList" binding:"required"`
}

func CreateNutritionList(c *gin.Context) {
	var nutritionListJson NutritionListJson

	err := c.ShouldBind(&nutritionListJson)

	if err != nil {
		helpers.ErrorHelper(c,
			helpers.JsonError{
				Message: "GetNutritionListByID: ID ERROR",
				Status:  400,
				Json:    helpers.JsonResponseType{Code: "INVALID_PARAM", Msg: "Bad Request"},
			},
		)
		return
	}

	userId, exists := c.Get("userId")

	if !exists {
		helpers.ErrorHelper(c,
			helpers.JsonError{
				Message: "Unauthorized",
				Status:  404,
				Json:    helpers.JsonResponseType{Code: "INVALID_ACCESS", Msg: "Unauthorized"},
			},
		)
		return
	}

	parsedUserId, err := uuid.Parse(fmt.Sprint(userId))

	if err != nil {
		helpers.ErrorHelper(c,
			helpers.JsonError{
				Message: "Unauthorized",
				Status:  404,
				Json:    helpers.JsonResponseType{Code: "INVALID_ACCESS", Msg: "Unauthorized"},
			},
		)
		return
	}

	var nutritionList []config.NutritionItem

	for _, item := range nutritionListJson.NutritionList {
		newNutritionItem := config.NutritionItem{
			Name:     item.Name,
			Calories: item.Calories,
			Protein:  item.Protein,
		}
		nutritionList = append(nutritionList, newNutritionItem)
	}

	newNutritionList := config.DailyNutritionList{
		UserID:         parsedUserId,
		Date:           nutritionListJson.Date,
		NutritionItems: nutritionList,
	}

	result := config.DB.Create(&newNutritionList)

	if result.Error != nil {
		helpers.NetworkError(c, result.Error)
		return
	}

	c.JSON(200, gin.H{
		"msg":                  "List Created",
		"dailyNutritionListId": newNutritionList.ID,
		"date":                 newNutritionList.Date,
		"nutritionItems":       newNutritionList.NutritionItems,
	})
}
