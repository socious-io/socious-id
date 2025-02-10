package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func LoginRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := FetchUserByJWT(c)

		if err != nil {
			u, err := FetchUserBySession(c)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
				return
			}
			user = u
		}

		c.Set("user", user)
		c.Next()
	}
}

func CheckLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, err := FetchUserByJWT(c)

		if err != nil {
			_, err := FetchUserBySession(c)
			if err != nil {
				c.Next()
				return
			}
		}
		c.Redirect(http.StatusTemporaryRedirect, "/auth/confirm")
	}
}
