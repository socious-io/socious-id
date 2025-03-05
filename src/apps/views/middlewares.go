package views

import (
	"net/http"
	"socious-id/src/apps/models"
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
			c.JSON(http.StatusBadRequest, gin.H{"error": "No organizations found for this user"})
			return
		}

		c.Set("organization", organization)
		c.Next()
	}
}
