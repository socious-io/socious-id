package auth

import (
	"net/http"
	"socious-id/src/apps/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

		//Safeguarding Identity if it was empty
		var identity *models.Identity
		identityStr := c.GetHeader(http.CanonicalHeaderKey("current-identity"))
		identityUUID, err := uuid.Parse(identityStr)
		if err == nil {
			identity, _ = models.GetIdentity(identityUUID)
		} else {
			identity, _ = models.GetIdentity(user.ID)
		}
		c.Set("identity", identity)

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
