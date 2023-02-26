package middleware

import (
	"net/http"

	"example.com/stock-manager/initializers"
	"example.com/stock-manager/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func ValidateToken(c *gin.Context) {
	token, err := extractToken(c)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Get token from database
		var t models.Token = models.Token{}
		var conn = initializers.DB.Table("tokens").Where("user_id=?", claims["userId"]).Order("updated_at DESC").First(&t)
		if conn.Error != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Failed to get token from database",
			})

			return
		}

		// Check if token is still valid
		if !t.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Token is no longer valid.",
			})

			return
		}
	} else {
		c.AbortWithStatus(http.StatusUnauthorized)
	}

	var requestOwner models.User
	initializers.DB.Take(&requestOwner, token.Claims.(jwt.MapClaims)["userId"])

}

// old check credentials
// tokenString, err := c.Cookie("Authorization")
// if err != nil {
// 	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
// 		"error": "Invalid credentials",
// 	})

// 	// c.Abort()
// 	return
// }

// token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
// 	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
// 	}

// 	return []byte(os.Getenv("SECRET")), nil
// })
// if err != nil {
// 	c.AbortWithStatus(http.StatusUnauthorized)
// 	return
// }
