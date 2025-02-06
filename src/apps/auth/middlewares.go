package auth

import (
	"net/http"
	"socious-id/src/apps/models"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func LoginRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := c.GetHeader("Authorization")
		splited := strings.Split(tokenStr, " ")
		if len(splited) > 1 {
			tokenStr = splited[1]
		} else {
			tokenStr = splited[0]
		}
		if tokenStr == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			c.Abort()
			return
		}

		claims, err := VerifyToken(tokenStr)

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token signature"})
				c.Abort()
				return
			}
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		u, err := models.GetUser(uuid.MustParse(claims.ID))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}
		c.Set("user", u)

		c.Next()
	}
}
