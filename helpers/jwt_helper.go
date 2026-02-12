package helpers

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var AUTH_SECRET = []byte(os.Getenv("AUTH_SECRET"))

func GenerateUserJWT(userId uuid.UUID, email string, exp int) (string, error) {
	claims := jwt.MapClaims{
		"sub": userId,
		"email": email,
		"exp": time.Now().Add(time.Duration(exp) * time.Hour).Unix(),
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(AUTH_SECRET)
}

func VerifyUserJWT(cookie string) (map[string]interface{}, error) {
	token, err := jwt.Parse(cookie, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS256 {
    		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return AUTH_SECRET, nil
	})

	if err != nil {
        return nil, fmt.Errorf("Unauthorized Token")
    }

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, fmt.Errorf("Unauthorized Token")
	}
}