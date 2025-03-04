package views

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"socious-id/src/apps/auth"
	"socious-id/src/apps/models"
	"socious-id/src/config"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/socious-io/gomail"
)

func authGroup(router *gin.Engine) {
	g := router.Group("auth")

	g.GET("/confirm", auth.LoginRequired(), func(c *gin.Context) {
		if authSession := loadAuthSession(c); authSession != nil {
			c.HTML(http.StatusOK, "confirm.html", gin.H{
				"User":        c.MustGet("user").(*models.User),
				"AuthSession": authSession,
			})
		}
		// NOTE: look like page sent without any session so detect it's self authorization
		c.Redirect(http.StatusPermanentRedirect, config.Config.Platforms.Accounts)

	})

	g.POST("/confirm", auth.LoginRequired(), func(c *gin.Context) {
		authSession := loadAuthSession(c)
		if authSession == nil {
			c.HTML(http.StatusNotAcceptable, "confirm.html", gin.H{
				"error": "not accepted without auth session",
			})
			return
		}

		form := new(ConfirmForm)
		c.ShouldBind(form)
		params := url.Values{}
		params.Add("session", authSession.ID.String())

		user := c.MustGet("user").(*models.User)
		ctx := c.MustGet("ctx").(context.Context)

		if !form.Confirmed {
			params.Add("status", "canceled")

			c.Redirect(http.StatusFound, fmt.Sprintf("%s?%s", authSession.RedirectURL, params.Encode()))
			return
		}

		otp := &models.OTP{
			UserID:        user.ID,
			AuthSessionID: &authSession.ID,
			Type:          models.SSOOTP,
		}

		if err := otp.Create(ctx); err != nil {
			c.HTML(http.StatusNotAcceptable, "confirm.html", gin.H{
				"error": err.Error(),
			})
			return
		}

		params.Add("status", "success")
		params.Add("code", otp.Code)

		c.Redirect(http.StatusFound, fmt.Sprintf("%s?%s", authSession.RedirectURL, params.Encode()))

	})

	g.GET("/login", auth.CheckLogin(), func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", gin.H{})
	})

	g.POST("/login", auth.CheckLogin(), func(c *gin.Context) {
		form := new(auth.LoginForm)
		if err := c.ShouldBind(form); err != nil {
			c.HTML(http.StatusBadRequest, "login.html", gin.H{
				"error": err.Error(),
			})
			return
		}

		u, err := models.GetUserByEmail(form.Email)
		if err != nil || u == nil {
			c.HTML(http.StatusBadRequest, "login.html", gin.H{
				"error": "Error: User couldn't be found/is not registered on Socious",
			})
			return
		}
		if u.Status == models.UserStatusInactive {
			c.HTML(http.StatusBadRequest, "login.html", gin.H{
				"error": "Error: user is not verified",
			})
			return
		}
		if u.Password == nil {
			c.HTML(http.StatusBadRequest, "login.html", gin.H{
				"error": "Error: email/password not match",
			})
			return
		}
		if err := auth.CheckPasswordHash(form.Password, *u.Password); err != nil {
			c.HTML(http.StatusBadRequest, "login.html", gin.H{
				"error": "Error: email/password not match",
			})
			return
		}

		session := sessions.Default(c)
		session.Set("user_id", u.ID.String())
		session.Save()
		// TODO: make sso otp if it has auth session
		c.Redirect(http.StatusFound, "/auth/confirm")
	})

	g.GET("/otp/confirm", func(c *gin.Context) {
		email := c.Query("email")
		code := c.Query("code")
		ctx := c.MustGet("ctx").(context.Context)

		otp, err := models.GetOTPByEmailAndCode(email, code)
		if err != nil {
			c.HTML(http.StatusBadRequest, "otp.html", gin.H{
				"error": err.Error(),
			})
			return
		}

		if otp.ExpireAt.Before(time.Now()) || otp.VerifiedAt != nil {
			c.HTML(http.StatusBadRequest, "otp.html", gin.H{
				"error": "code has been expired",
			})
			return
		}

		if otp.AuthSession == nil {
			c.HTML(http.StatusBadRequest, "otp.html", gin.H{
				"error": "not valid otp for this session",
			})
			return
		}

		if otp.AuthSession.ExpireAt.Before(time.Now()) || otp.AuthSession.VerifiedAt != nil {
			c.HTML(http.StatusBadRequest, "otp.html", gin.H{
				"error": "auth session has been expired",
			})
			return
		}

		if err := otp.Verify(ctx); err != nil {
			c.HTML(http.StatusBadRequest, "otp.html", gin.H{
				"error": err.Error(),
			})
			return
		}

		if err := otp.User.Verify(ctx, models.VerificationTypeEmail); err != nil {
			c.HTML(http.StatusBadRequest, "otp.html", gin.H{
				"error": err.Error(),
			})
			return
		}

		//Saving into session
		session := sessions.Default(c)
		session.Set("user_id", otp.User.ID.String())
		session.Save()

		if otp.Type == models.VerificationOTP {
			c.Redirect(http.StatusSeeOther, "/users/profile")
		} else if otp.Type == models.ForgetPasswordOTP {
			c.Redirect(http.StatusSeeOther, "/auth/password/set")
		}
	})

	g.GET("/otp", func(c *gin.Context) {
		email := c.Query("email")

		c.HTML(http.StatusOK, "otp.html", gin.H{
			"email": email,
		})
	})

	g.POST("/register", auth.CheckLogin(), func(c *gin.Context) {
		authSession := loadAuthSession(c)
		if authSession == nil {
			c.HTML(http.StatusNotAcceptable, "confirm.html", gin.H{
				"error": "not accepted without auth session",
			})
			return
		}

		ctx := c.MustGet("ctx").(context.Context)

		form := new(auth.OTPForm)
		if err := c.ShouldBind(form); err != nil {
			c.HTML(http.StatusBadRequest, "register.html", gin.H{
				"error": err.Error(),
			})
			return
		}

		//Creating user (Default in INACTIVE state)
		u := &models.User{
			Username: form.Email, //TODO: generate username
			Email:    form.Email,
		}

		if err := u.Create(ctx); err != nil {
			c.HTML(http.StatusBadRequest, "register.html", gin.H{
				"error": err.Error(),
			})
			return
		}

		//Save OTP
		otp := &models.OTP{
			UserID:        u.ID,
			AuthSessionID: &authSession.ID,
			Type:          models.VerificationOTP,
		}

		if err := otp.Create(ctx); err != nil {
			c.HTML(http.StatusNotAcceptable, "register.html", gin.H{
				"error": err.Error(),
			})
			return
		}

		//Email OTP
		items := map[string]string{"code": otp.Code}
		gomail.SendEmail(gomail.EmailConfig{
			Approach:    gomail.EmailApproachTemplate,
			Destination: u.Email,
			Title:       "OTP Code",
			TemplateId:  "otp",
			Args:        items,
		})

		c.Redirect(http.StatusSeeOther, fmt.Sprintf("/auth/otp?email=%s", form.Email))
	})

	g.GET("/register", auth.CheckLogin(), func(c *gin.Context) {
		c.HTML(http.StatusOK, "register.html", gin.H{})
	})

	g.GET("/register/pre", func(c *gin.Context) {
		c.HTML(http.StatusOK, "pre-register.html", gin.H{})
	})

	g.POST("/password/forget", auth.CheckLogin(), func(c *gin.Context) {
		authSession := loadAuthSession(c)
		if authSession == nil {
			c.HTML(http.StatusNotAcceptable, "confirm.html", gin.H{
				"error": "not accepted without auth session",
			})
			return
		}

		ctx := c.MustGet("ctx").(context.Context)

		form := new(auth.OTPForm)
		if err := c.ShouldBind(form); err != nil {
			c.HTML(http.StatusBadRequest, "forget-password.html", gin.H{
				"error": err.Error(),
			})
			return
		}

		//Fetching user
		u, err := models.GetUserByEmail(form.Email)
		if err != nil {
			c.HTML(http.StatusBadRequest, "forget-password.html", gin.H{
				"error": err.Error(),
			})
			return
		}

		//Checking user status
		if u.Status == models.UserStatusInactive {
			c.HTML(http.StatusBadRequest, "forget-password.html", gin.H{
				"error": "Error: user is not verified",
			})
			return
		}

		//Save OTP
		otp := &models.OTP{
			UserID:        u.ID,
			AuthSessionID: &authSession.ID,
			Type:          models.ForgetPasswordOTP,
		}

		if err := otp.Create(ctx); err != nil {
			c.HTML(http.StatusNotAcceptable, "forget-password.html", gin.H{
				"error": err.Error(),
			})
			return
		}

		//Email OTP
		items := map[string]string{"code": otp.Code}
		gomail.SendEmail(gomail.EmailConfig{
			Approach:    gomail.EmailApproachTemplate,
			Destination: u.Email,
			Title:       "Forget Password OTP Code",
			TemplateId:  "otp",
			Args:        items,
		})

		c.Redirect(http.StatusSeeOther, fmt.Sprintf("/auth/otp?email=%s", u.Email))
	})

	g.GET("/password/forget", auth.CheckLogin(), func(c *gin.Context) {
		c.HTML(http.StatusOK, "forget-password.html", gin.H{})
	})

	g.POST("/password/set", auth.LoginRequired(), func(c *gin.Context) {
		authSession := loadAuthSession(c)
		if authSession == nil {
			c.HTML(http.StatusNotAcceptable, "confirm.html", gin.H{
				"error": "not accepted without auth session",
			})
			return
		}

		user := c.MustGet("user").(*models.User)
		ctx := c.MustGet("ctx").(context.Context)

		form := new(auth.SetPasswordForm)
		if err := c.ShouldBind(form); err != nil {
			c.HTML(http.StatusBadRequest, "set-password.html", gin.H{
				"error": err.Error(),
			})
			return
		}

		password, _ := auth.HashPassword(form.Password)
		user.Password = &password
		if err := user.UpdatePassword(ctx); err != nil {
			c.HTML(http.StatusBadRequest, "set-password.html", gin.H{
				"error": err.Error(),
			})
			return
		}

		c.Redirect(http.StatusSeeOther, "/auth/password/set/confirm")
	})

	g.GET("/password/set", auth.LoginRequired(), func(c *gin.Context) {
		c.HTML(http.StatusOK, "set-password.html", gin.H{})
	})

	g.GET("/password/set/confirm", func(c *gin.Context) {
		c.HTML(http.StatusOK, "post-set-password.html", gin.H{})
	})

	g.DELETE("/logout", auth.LoginRequired(), func(c *gin.Context) {
		session := sessions.Default(c)
		session.Delete("user_id")
		session.Save()
		c.Redirect(http.StatusPermanentRedirect, "/auth/login")
	})

	g.POST("/session", func(c *gin.Context) {
		ctx := c.MustGet("ctx").(context.Context)
		form := new(AuthSessionForm)

		if err := c.ShouldBind(form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		access, err := models.GetAccessByClientID(form.ClientID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		if err := auth.CheckPasswordHash(form.ClientSecret, access.ClientSecret); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "client access not valid",
			})
			return
		}

		authSession := &models.AuthSession{
			RedirectURL: form.RedirectURL,
			AccessID:    access.ID,
			ExpireAt:    time.Now().Add(time.Minute * 10),
		}

		if err := authSession.Create(ctx); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"auth_session": authSession,
		})

	})

	g.GET("/session/:id", func(c *gin.Context) {
		id, err := uuid.Parse(c.Param("id"))
		authMode := models.AuthModeType(c.Query("auth_mode"))

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		authSession, err := models.GetAuthSession(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		if authSession.ExpireAt.Before(time.Now()) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "session has been expired",
			})
			return
		}

		session := sessions.Default(c)
		session.Set("auth_session_id", authSession.ID.String())
		session.Save()

		if authMode == models.AuthModeRegister {
			c.Redirect(http.StatusPermanentRedirect, "/auth/register/pre")
		} else {
			c.Redirect(http.StatusPermanentRedirect, "/auth/login")
		}

	})

	g.POST("/session/token", func(c *gin.Context) {
		ctx := c.MustGet("ctx").(context.Context)
		form := new(GetTokenForm)
		if err := c.ShouldBind(form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		access, err := models.GetAccessByClientID(form.ClientID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		if err := auth.CheckPasswordHash(form.ClientSecret, access.ClientSecret); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "client access not valid",
			})
			return
		}

		otp, err := models.GetOTPByCode(form.Code)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "not valid auth code",
			})
			return
		}

		if otp.ExpireAt.Before(time.Now()) || otp.VerifiedAt != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "code has been expired",
			})
			return
		}

		if otp.AuthSession == nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "not valid otp for this session",
			})
		}

		if otp.AuthSession.ExpireAt.Before(time.Now()) || otp.AuthSession.VerifiedAt != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "auth session has been expired",
			})
			return
		}

		tokens, err := auth.Signin(otp.User.ID.String(), otp.User.Email)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		if err := otp.Verify(ctx); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusAccepted, tokens)
	})
}

func loadAuthSession(c *gin.Context) *models.AuthSession {
	session := sessions.Default(c)

	if authSessionID := session.Get("auth_session_id"); authSessionID != nil {
		authSession, err := models.GetAuthSession(uuid.MustParse(authSessionID.(string)))
		if err == nil {
			return authSession
		}
	}

	return nil
}
