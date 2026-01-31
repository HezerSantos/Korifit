package controllers

import (
	"Korifit/config"
	"Korifit/helpers"
	"errors"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)



type CreateUserJSON struct {
	Email string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	ConfirmPassword string `json:"confirmPassword" binding:"required,eqfield=Password"`
}




func CreateUser(c *gin.Context) {
	var newUser CreateUserJSON

	err := c.ShouldBind(&newUser)
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

	hashedPassowrd, err := helpers.HashPassword(newUser.Password)

	if err != nil {
		helpers.NetworkError(c, err)
		return
	}

	createdUser := config.User{Email: newUser.Email, Password: hashedPassowrd}

	config.DB.Create(&createdUser)

	c.JSON(200, gin.H{"msg": "User Created Successfully"})
}



type VerifyUserJSON struct {
	Email string `json:"email" binding:"email,required"`
	Password string `json:"password" binding:"required"`
}

func VerifyUser(c *gin.Context) {
	var verifyUserJson VerifyUserJSON
	err := c.ShouldBind(&verifyUserJson)

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

	var user config.User

	result := config.DB.Where("email = ?", verifyUserJson.Email).First(&user)

	//Error if record not found
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		helpers.ErrorHelper(
			c,
			helpers.JsonError{
				Message: "User not found (Email)",
				Status:  404,
				Json: helpers.JsonResponseType{Code: "INVALID_USER", Msg: "User not found"},
			},
		)
		return
	} 

	//Network error
	if result.Error != nil {
		helpers.NetworkError(c, result.Error)
		return
	}

	//Verify the password
	passwordResult, err := helpers.ComparePasswordAndHash(verifyUserJson.Password, user.Password)
	

	if err != nil {
		helpers.NetworkError(c, err)
		return
	}

	if !passwordResult {
		helpers.ErrorHelper(
			c,
			helpers.JsonError{
				Message: "User not found (Passwords)",
				Status:  404,
				Json: helpers.JsonResponseType{Code: "INVALID_USER", Msg: "User not found"},
			},
		)
		return
	}


	jwtToken, err := helpers.GenerateUserJWT(user.ID, user.Email, 1)

	if err != nil {
		helpers.NetworkError(c, err)
		return
	}

	domain := ""

	if os.Getenv("GO_ENV") == "production" {
		domain = ".hallowedvisions.com"
	}

	c.SetCookie(
		"__Secure-secure-auth.access", 
		jwtToken, 
		60 * 1000 * 60, 
		"/", 
		domain, 
		true, 
		true,
	)

	c.JSON(200, gin.H{"msg": "User Logged In"})
}

