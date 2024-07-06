// package middlewares

// import (
// 	"net/http"
// 	"strings"

// 	"github.com/gin-gonic/gin"
// 	"github.com/golang-jwt/jwt/v4"
// )

// func AuthMiddleware() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		authHeader := c.GetHeader("Authorization")
// 		if authHeader == "" {
// 			c.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization header required"})
// 			c.Abort()
// 			return
// 		}

// 		// Format header Authorization: Bearer <token>
// 		tokenString := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))

// 		claims, err := utils.ValidateJWT(tokenString)
// 		if err != nil {
// 			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid token"})
// 			c.Abort()
// 			return
// 		}

//			c.Set("user_id", claims.UserID)
//			c.Next()
//		}
//	}
package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/raihansyahrin/backend_laundry_app.git/utils"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization header missing"})
			c.Abort()
			return
		}

		tokenString := strings.Split(authHeader, "Bearer ")[1]

		claims, err := utils.ValidateJWT(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid token", "error": err.Error()})
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("role", claims.Role)

		c.Next()
	}
}
