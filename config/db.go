package config

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func DatabaseConnect(){
	fmt.Println("Connecting to Database...")

	dsn := os.Getenv("DATABASE_URL")

	if dsn == "" {
		panic("Missing Database URL")
	}

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	DB = database
	fmt.Println("Finished Connecting to Database")
}

func DatabaseMigrate() {
    fmt.Println("   MIGRATING DATABASE...")
	DB.AutoMigrate(&User{}, &Exercise{}, &Workout{}, &DailyNutritionList{}, &NutritionItem{})
    fmt.Println("   FINISHED MIGRATING DATABASE...")
}