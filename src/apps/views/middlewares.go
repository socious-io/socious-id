package views

import (
	"bytes"
	"io"
	"net/http"
	"socious-id/src/apps/auth"
	"socious-id/src/apps/models"
	"socious-id/src/config"
	"strconv"
	"strings"

	"github.com/google/uuid"
	database "github.com/socious-io/pkg_database"

	"github.com/gin-gonic/gin"
)

func paginate() gin.HandlerFunc {
	return func(c *gin.Context) {

		page, err := strconv.Atoi(c.Query("page"))
		if err != nil {
			page = 1
		}

		limit, err := strconv.Atoi(c.Query("limit"))
		if err != nil {
			limit = 10
		}
		if page < 1 {
			page = 1
		}
		if limit > 100 || limit < 1 {
			limit = 10
		}
		filters := make([]database.Filter, 0)
		for key, values := range c.Request.URL.Query() {
			if strings.Contains(key, "filter.") && len(values) > 0 {
				filters = append(filters, database.Filter{
					Key:   strings.Replace(key, "filter.", "", -1),
					Value: values[0],
				})
			}
		}

		c.Set("paginate", database.Paginate{
			Limit:   limit,
			Offet:   (page - 1) * limit,
			Filters: filters,
		})
		c.Set("limit", limit)
		c.Set("page", page)
		c.Next()

	}
}

func isOrgMember() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := uuid.MustParse(c.Param("id")) //Org ID
		user := c.MustGet("user").(*models.User)
		organization, err := models.GetOrganizationByMember(id, user.ID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "No organizations found for this user"})
			return
		}

		c.Set("organization", organization)
		c.Next()
	}
}

func clientSecretRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		bodyBytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "failed to read request body"})
			c.Abort()
			return
		}

		// Restore the body so it can be read by ShouldBindJSON
		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		authData := new(ClientSecretForm)
		if err := c.ShouldBindJSON(&authData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			c.Abort()
			return
		}
		// Checking Client
		access, err := models.GetAccessByClientID(authData.ClientID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			c.Abort()
			return
		}

		if err := auth.CheckPasswordHash(authData.ClientSecret, access.ClientSecret); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "client access not valid",
			})
			c.Abort()
			return
		}
		// Store the client ID in context for later use if needed
		c.Set("access", access)

		// Restore the body so it can be read by Handler
		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		c.Next()
	}
}

func NoCache() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Continue to handler
		c.Next()

		c.Writer.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate")
		c.Writer.Header().Set("Pragma", "no-cache")
		c.Writer.Header().Set("Expires", "0")
	}
}

// Administration
func adminAccessRequired() gin.HandlerFunc {
	return func(c *gin.Context) {

		access_token := c.Query("admin_access_token")
		isAdmin := access_token == config.Config.AdminToken

		if !isAdmin {
			c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
			c.Abort()
			return
		}

		c.Next()
	}
}
