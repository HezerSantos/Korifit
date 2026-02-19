package middleware

import (
	"Korifit/helpers"
	"github.com/gin-gonic/gin"
	"os"
)

var AUTH_SECRET = []byte(os.Getenv("AUTH_SECRET"))

func AuthenticateUser(c *gin.Context) {
	cookie, err := c.Cookie("__Secure-secure-auth.access")

	if err != nil {
		helpers.ErrorHelper(c,
			helpers.JsonError{
				Message: "Secure Auth Cookie NA",
				Status:  404,
				Json: helpers.JsonResponseType{
					Msg:  "Unauthorized",
					Code: "INVALID_AUTH",
				},
			},
		)
		c.Abort()
		return
	}

	claims, err := helpers.VerifyUserJWT(cookie, AUTH_SECRET)

	if err != nil {
		helpers.ErrorHelper(c,
			helpers.JsonError{
				Message: "Secure Auth Cookie NA",
				Status:  404,
				Json: helpers.JsonResponseType{
					Msg:  "Unauthorized",
					Code: "INVALID_AUTH",
				},
			},
		)
		c.Abort()
		return
	}

	id, _ := claims["sub"]
	c.Set("userId", id)
	c.Next()
}
