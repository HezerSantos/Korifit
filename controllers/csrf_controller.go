package controllers

import (
	"Korifit/helpers"
	"crypto/rand"
	"encoding/hex"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	r "math/rand"
	"time"
)

var CSRF_SECRET = []byte(os.Getenv("CSRF_SECRET"))

func GetCsrfToken(c *gin.Context) {
	b := make([]byte, 32)
	_, err := rand.Read(b)

	if err != nil {
		helpers.NetworkError(c, err)
		c.Abort()
		return
	}

	hexStr := hex.EncodeToString(b)
	r.Seed(time.Now().UnixNano())

	n := r.Intn(10)

	claims := jwt.MapClaims{
		"csrfToken": hexStr,
		"key":       n,
		"exp":       time.Now().Add(time.Duration(15) * time.Minute).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signed, err := token.SignedString(CSRF_SECRET)

	if err != nil {
		helpers.NetworkError(c, err)
		return
	}

	domain := ""

	if os.Getenv("GO_ENV") == "production" {
		domain = ".hallowedvisions.com"
	}

	c.SetCookie(
		"__Secure-auth.csrf",
		signed,
		5*1000*60,
		"/",
		domain,
		true,
		false,
	)

	c.JSON(200, gin.H{"msg": "csrf retrieved"})
}
