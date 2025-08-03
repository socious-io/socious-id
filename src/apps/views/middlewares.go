package views

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"socious-id/src/apps/auth"
	"socious-id/src/apps/models"
	"socious-id/src/apps/utils"
	"socious-id/src/config"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/microcosm-cc/bluemonday"
	database "github.com/socious-io/pkg_database"
	"github.com/unrolled/secure"

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
		c.Writer.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate")
		c.Writer.Header().Set("Pragma", "no-cache")
		c.Writer.Header().Set("Expires", "0")

		// Continue to handler
		c.Next()
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

func SecureHeaders(env string) gin.HandlerFunc {

	IsDevelopment := env != "production"
	options := secure.Options{
		FrameDeny:          true, // X-Frame-Options: DENY
		ContentTypeNosniff: true, // X-Content-Type-Options: nosniff
		BrowserXssFilter:   true, // X-XSS-Protection: 1; mode=block (legacy)
		// ReferrerPolicy:        "no-referrer",
		ContentSecurityPolicy: "default-src 'self'; script-src 'self' $NONCE; img-src 'self' https: http:;", // Very important for XSS
		// HSTS:
		SSLRedirect:          true,
		STSSeconds:           31536000,
		STSIncludeSubdomains: true,
		STSPreload:           true,
		IsDevelopment:        IsDevelopment,
	}

	return func(c *gin.Context) {
		s := secure.New(options)
		nonce, err := s.ProcessAndReturnNonce(c.Writer, c.Request)
		if err != nil {
			c.AbortWithStatus(500)
			return
		}
		c.Set("nonce", nonce)

		c.Next()
	}
}

func SecureRequest(p *bluemonday.Policy) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check content type
		isUrlEncodedContent := strings.Contains(c.GetHeader("Content-Type"), "application/x-www-form-urlencoded")
		isMultipartContent := strings.Contains(c.GetHeader("Content-Type"), "multipart/form-data")
		isJsonContent := strings.Contains(c.GetHeader("Content-Type"), "application/json")

		// --- 1. Sanitize Query Parameters ---
		q := c.Request.URL.Query()
		utils.SanitizeURLValues(q, p)
		c.Request.URL.RawQuery = q.Encode()

		// --- 2. Sanitize Form Data (application/x-www-form-urlencoded or multipart) ---
		if isUrlEncodedContent || isMultipartContent {
			if err := c.Request.ParseForm(); err != nil {
				c.AbortWithStatusJSON(400, gin.H{
					"error": fmt.Sprintf("Invalid body payload, err: %v", err),
				})
				return
			}
			utils.SanitizeURLValues(c.Request.PostForm, p)
		} else if isJsonContent {
			var bodyBytes []byte
			if c.Request.Body != nil {
				bodyBytes, _ = io.ReadAll(c.Request.Body)
			}

			if len(bodyBytes) > 0 {
				var data map[string]interface{}
				if err := json.Unmarshal(bodyBytes, &data); err != nil {
					c.AbortWithStatusJSON(400, gin.H{
						"error": fmt.Sprintf("Invalid body payload, err: %v", err),
					})
					return
				}

				utils.SanitizeMap(data, p)
				safeBody, _ := json.Marshal(data)
				c.Request.Body = io.NopCloser(bytes.NewReader(safeBody))
			}
		}
		c.Next()
	}
}
