package config

import (
	"github.com/google/uuid"
)

type User struct {
	ID       uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Email    string    `gorm:"unique" json:"email"`
	Password string    `json:"password"`

	Workouts            []Workout            `gorm:"foreignKey:UserID" json:"workouts"`
	DailyNutritionLists []DailyNutritionList `gorm:"foreignKey:UserID" json:"dailyNutritionLists"`
}

type Exercise struct {
	ID           uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Name         string    `gorm:"unique" json:"name"`
	MuscleTarget string    `json:"muscleTarget"`

	Workouts []Workout `gorm:"many2many:workout_exercises;" json:"workouts"`
}

type Workout struct {
	ID   uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Name string    `gorm:"unique" json:"name"`

	Exercises []Exercise `gorm:"many2many:workout_exercises;" json:"exercises"`
	UserID    uuid.UUID  `gorm:"type:uuid" json:"userId"`
}

type NutritionItem struct {
	ID       uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Name     string    `json:"name"`
	Calories int       `json:"calories"`
	Protein  int       `json:"protein"`

	DailyNutritionListID uuid.UUID `gorm:"type:uuid" json:"dailyNutritionListId"`
}

type DailyNutritionList struct {
	ID   uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Date string    `json:"date"`

	NutritionItems []NutritionItem `json:"nutritionItems"`
	UserID         uuid.UUID       `gorm:"type:uuid" json:"userId"`
}
