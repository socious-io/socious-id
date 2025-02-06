package views

import (
	"net/http"
	"net/url"
	"socious-id/src/apps/auth"
	"socious-id/src/apps/models"

	"github.com/gin-gonic/gin"
)

func authGroup(router *gin.Engine) {
	g := router.Group("sso")

	router.LoadHTMLGlob("src/apps/templates/*.html")
	router.Static("/public", "src/apps/templates/public")

	g.GET("/login", func(c *gin.Context) {
		redirect_url := c.Query("redirect_url")

		c.HTML(http.StatusOK, "login.html", gin.H{
			"redirect_url": redirect_url,
		})
	})

	g.POST("/login", func(c *gin.Context) {
		redirect_url := c.Query("redirect_url")

		loginForm := new(auth.LoginForm)
		if err := c.ShouldBind(loginForm); err != nil {
			c.HTML(http.StatusBadRequest, "login.html", gin.H{
				"redirect_url": redirect_url,
				"error":        err.Error(),
			})
			return
		}

		u, err := models.GetUserByEmail(loginForm.Email)
		if err != nil || u == nil {
			c.HTML(http.StatusBadRequest, "login.html", gin.H{
				"redirect_url": redirect_url,
				"error":        "Error: User couldn't be found/is not registered on Socious",
			})
			return
		}
		if u.Password == nil {
			c.HTML(http.StatusBadRequest, "login.html", gin.H{
				"redirect_url": redirect_url,
				"error":        "Error: email/password not match",
			})
			return
		}
		if err := auth.CheckPasswordHash(loginForm.Password, *u.Password); err != nil {
			c.HTML(http.StatusBadRequest, "login.html", gin.H{
				"redirect_url": redirect_url,
				"error":        "Error: email/password not match",
			})
			return
		}
		// TODO: make temp-token
		//Add redirect_url to query params
		parsedURL, err := url.Parse(redirect_url)
		if err != nil {
			return
		}
		queryParams := parsedURL.Query()
		queryParams.Add("token", "@TODO:TOKEN")
		parsedURL.RawQuery = queryParams.Encode()
		c.Redirect(http.StatusSeeOther, parsedURL.String())

		return
	})

	g.GET("/register", func(c *gin.Context) {
		redirect_url := c.Query("redirect_url")

		c.HTML(http.StatusOK, "register.html", gin.H{
			"redirect_url": redirect_url,
		})
	})

}
