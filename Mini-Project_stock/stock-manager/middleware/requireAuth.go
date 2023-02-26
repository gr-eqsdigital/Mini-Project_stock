package middleware

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"example.com/stock-manager/initializers"
	"example.com/stock-manager/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func extractToken(c *gin.Context) (*jwt.Token, error) {
	// Get the cookie off requisitionOwner
	tokenString, err := c.Cookie("Authorization")
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return nil, err
	}

	// Decode and validate it
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("SECRET")), nil
	})
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return nil, err
	}

	return token, nil
}

// Verify user credentials
func RequireAuth(c *gin.Context) {
	token, err := extractToken(c)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Check exp date
		if float64(time.Now().Unix()) > claims["expirationDate"].(float64) {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// Find the user with token
		var user models.User
		initializers.DB.First(&user, claims["userId"])

		if user.ID == 0 {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// Attach to requision
		c.Set("user", user)

		// Reply requision - Continue
		c.Next()
		// Reply requision - Deny Access
	} else {
		c.AbortWithStatus(http.StatusUnauthorized)
	}
}
