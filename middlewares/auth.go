package middlewares

import (
	"fmt"
	"log"
	"models"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"

	"github.com/gin-gonic/gin"
)

func LoginAuth() gin.HandlerFunc {
	return func(c *gin.Context) {

		tokenStr := c.GetHeader("Authorization")

		if tokenStr == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "No token Provided",
			})
			return
		}

		// Check if it use Bearer

		if s := strings.Split(tokenStr, " "); len(s) == 2 {
			tokenStr = s[1]
		}

		// Problem: token is generatted but not valid

		// Add Claim Latter
		token, _ := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Token"})
				return nil, fmt.Errorf("Invalid Token")
			}

			return []byte("secert"), nil
		})

		fmt.Printf("The token:\n %+v\n", token)

		fmt.Printf("Claims:\n%+v\n", token.Claims.(jwt.MapClaims))

		// Find the User and store in c

		claims := token.Claims.(jwt.MapClaims)

		user, err := models.FindUserByEmail(claims["email"].(string))

		if err != nil {
			log.Print(err)
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "User Not Authorised",
			})
		}

		c.Set("user", user)

		if !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"getToken": tokenStr,
				"error":    "Token is not valid",
			})
			return
		}

		c.Next()
		// t := time.Now()
	}
}

func TestingMeiddelWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"testing": "sucess in the middleware",
		})
		c.Next()
	}
}
