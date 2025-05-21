package views

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
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
			user := c.MustGet("user").(*models.User)
			organizations, _ := models.GetOrganizationsByMember(user.ID)
			c.HTML(http.StatusOK, "confirm.html", gin.H{
				"User":          user,
				"Organizations": organizations,
				"AuthSession":   authSession,
			})
		}

		session := sessions.Default(c)

		if session.Get("next") != nil {
			next := session.Get("next").(string)
			session.Delete("next")
			session.Save()
			c.Redirect(http.StatusSeeOther, next)
			return
		}
		// NOTE: look like page sent without any session so detect it's self authorization
		c.Redirect(http.StatusTemporaryRedirect, config.Config.Platforms.Accounts)

	})

	g.POST("/confirm", auth.LoginRequired(), func(c *gin.Context) {
		authSession := loadAuthSession(c)
		if authSession == nil {
			c.HTML(http.StatusNotAcceptable, "confirm.html", gin.H{
				"error": "not accepted without auth session",
			})
			return
		}

		//Auth session is completed further redirection will go through account center
		session := sessions.Default(c)
		session.Delete("auth_session_id")
		session.Save()

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
		params.Add("identity_id", form.IdentityId)

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
		if u.Status == models.StatusTypeInactive {
			c.HTML(http.StatusBadRequest, "login.html", gin.H{
				"error": "Error: User couldn't be found/is not registered on Socious",
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

	g.GET("/google", func(c *gin.Context) {
		url := fmt.Sprintf(
			"https://accounts.google.com/o/oauth2/v2/auth?client_id=%s&redirect_uri=%s&response_type=code&scope=email profile&access_type=offline&prompt=consent",
			config.Config.Oauth.Google.ID,
			fmt.Sprintf("%s/auth/google/callback", config.Config.Host),
		)

		c.Redirect(http.StatusTemporaryRedirect, url)
	})

	g.GET("/google/callback", func(c *gin.Context) {
		ctx := c.MustGet("ctx").(context.Context)
		authorizationCode := c.Query("code")

		googleUserInfo, err := auth.GoogleLoginWithCode(authorizationCode, fmt.Sprintf("%s/auth/google/callback", config.Config.Host))
		if err != nil || googleUserInfo.Email == "" {
			c.HTML(http.StatusBadRequest, "login.html", gin.H{
				"error": "Error: Google login failed",
			})
			return
		}

		u, err := models.GetUserByEmail(googleUserInfo.Email)
		if err != nil && errors.Is(err, sql.ErrNoRows) {
			u = &models.User{
				Username:  auth.GenerateUsername(googleUserInfo.Email),
				Email:     googleUserInfo.Email,
				FirstName: &googleUserInfo.GivenName,
				LastName:  &googleUserInfo.FamilyName,
			}
			if err = u.Create(ctx); err != nil {
				c.HTML(http.StatusBadRequest, "login.html", gin.H{
					"error": err.Error(),
				})
				return
			}
		}

		//Make user active
		if u.Status == models.StatusTypeInactive {
			err := u.UpdateStatus(ctx, models.StatusTypeActive)
			if err != nil {
				c.HTML(http.StatusBadRequest, "login.html", gin.H{
					"error": err.Error(),
				})
				return
			}
		}

		//Set session
		session := sessions.Default(c)
		session.Set("user_id", u.ID.String())
		session.Save()

		//Redirect
		if u.Password == nil {
			c.Redirect(http.StatusSeeOther, "/auth/password/set")
			return
		}
		c.Redirect(http.StatusSeeOther, "/auth/confirm")
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

		// TEMPORARY organization creation need this to reverify auth session as we check OTP is new no worries
		/* if otp.AuthSession != nil && (otp.AuthSession.ExpireAt.Before(time.Now()) || otp.AuthSession.VerifiedAt != nil) {
			c.HTML(http.StatusBadRequest, "otp.html", gin.H{
				"error": "auth session has been expired",
			})
			return
		} */

		if err := otp.Verify(ctx, false); err != nil {
			c.HTML(http.StatusBadRequest, "otp.html", gin.H{
				"error": err.Error(),
			})
			return
		}

		if err := otp.User.Verify(ctx, models.UserVerificationTypeEmail); err != nil {
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

		ctx := c.MustGet("ctx").(context.Context)
		session := sessions.Default(c)

		form := new(auth.OTPForm)
		if err := c.ShouldBind(form); err != nil {
			c.HTML(http.StatusBadRequest, "register.html", gin.H{
				"error": err.Error(),
			})
			return
		}

		//Creating user (Default in INACTIVE state)
		u := &models.User{
			Username: auth.GenerateUsername(form.Email),
			Email:    form.Email,
		}

		//Set Referred By
		referredBy := session.Get("referred_by")
		if referredBy != nil && referredBy.(string) != "" {
			refererUser, err := models.GetUserByEmail(referredBy.(string))
			if err != nil {
				log.Printf("Couldn't find the referer user : %s\n", err.Error())
			}
			u.ReferredBy = &refererUser.ID
			session.Delete("referred_by")
			session.Save()
		}

		if err := u.Create(ctx); err != nil {
			c.HTML(http.StatusBadRequest, "register.html", gin.H{
				"error": "Email is already in use. Please select different email.",
			})
			return
		}

		//Save OTP
		otp := &models.OTP{
			UserID: u.ID,
			Type:   models.VerificationOTP,
		}

		if authSession != nil {
			otp.AuthSessionID = &authSession.ID
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

		log.Printf("OTP for email %s is `%s` \n", u.Email, otp.Code)

		c.Redirect(http.StatusSeeOther, fmt.Sprintf("/auth/otp?email=%s", form.Email))
	})

	g.GET("/register", auth.CheckLogin(), func(c *gin.Context) {
		c.HTML(http.StatusOK, "register.html", gin.H{})
	})

	g.GET("/register/pre", func(c *gin.Context) {
		c.HTML(http.StatusOK, "pre-register.html", gin.H{})
	})

	g.PUT("/password", auth.LoginRequired(), func(c *gin.Context) {
		ctx := c.MustGet("ctx").(context.Context)
		u, _ := c.MustGet("user").(*models.User)

		form := new(auth.ChangePasswordForm)
		if err := c.ShouldBindJSON(form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := auth.CheckPasswordHash(form.CurrentPassword, *u.Password); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "email/password not match"})
			return
		}

		newPassword, err := auth.HashPassword(form.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		u.Password = &newPassword
		if err := u.UpdatePassword(ctx); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusAccepted, gin.H{"message": "success"})
	})

	g.POST("/password/forget", auth.CheckLogin(), func(c *gin.Context) {

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
		if u.Status == models.StatusTypeInactive {
			c.HTML(http.StatusBadRequest, "forget-password.html", gin.H{
				"error": "Error: User couldn't be found/is not registered on Socious",
			})
			return
		}

		//Save OTP
		otp := &models.OTP{
			UserID: u.ID,
			Type:   models.ForgetPasswordOTP,
		}

		authSession := loadAuthSession(c)
		if authSession != nil {
			otp.AuthSessionID = &authSession.ID
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

		log.Printf("OTP for email %s is `%s` \n", u.Email, otp.Code)

		c.Redirect(http.StatusSeeOther, fmt.Sprintf("/auth/otp?email=%s", u.Email))
	})

	g.GET("/password/forget", auth.CheckLogin(), func(c *gin.Context) {
		c.HTML(http.StatusOK, "forget-password.html", gin.H{})
	})

	g.POST("/password/set", auth.LoginRequired(), func(c *gin.Context) {

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
		c.Redirect(http.StatusSeeOther, "/auth/login")
	})

	g.GET("/logout", auth.LoginRequired(), func(c *gin.Context) {
		session := sessions.Default(c)
		session.Delete("user_id")
		session.Save()
		c.Redirect(http.StatusSeeOther, "/auth/login")
	})

	g.POST("/session", clientSecretRequired(), func(c *gin.Context) {
		ctx := c.MustGet("ctx").(context.Context)
		form := new(AuthSessionForm)
		if err := c.ShouldBind(form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		access := c.MustGet("access").(*models.Access)

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
		orgOnboard := c.Query("org_onboard")
		referredBy := c.Query("referred_by")

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

		fmt.Println("setting auth session id", authSession.ID.String())

		session := sessions.Default(c)
		session.Set("auth_session_id", authSession.ID.String())
		session.Set("org_onboard", orgOnboard == "true")
		if referredBy != "" {
			session.Set("referred_by", referredBy)
		}
		session.Save()

		if authMode == models.AuthModeRegister {
			c.Redirect(http.StatusTemporaryRedirect, "/auth/register/pre")
		} else {
			c.Redirect(http.StatusTemporaryRedirect, "/auth/login")
		}

	})

	g.POST("/session/token", clientSecretRequired(), func(c *gin.Context) {
		ctx := c.MustGet("ctx").(context.Context)
		form := new(GetTokenForm)
		if err := c.ShouldBind(form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
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

		// TEMPORARY organization creation need this to reverify auth session as we check OTP is new no worries
		/* if otp.AuthSession.ExpireAt.Before(time.Now()) || otp.AuthSession.VerifiedAt != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "auth session has been expired",
			})
			return
		} */

		tokens, err := auth.Signin(otp.User.ID.String(), otp.User.Email)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		if err := otp.Verify(ctx, true); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusAccepted, tokens)
	})

	g.POST("/refresh", clientSecretRequired(), func(c *gin.Context) {
		form := new(RefreshTokenForm)
		if err := c.ShouldBind(form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		claims, err := auth.VerifyToken(form.RefreshToken)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		tokens, err := auth.Signin(claims.ID, claims.Email)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, tokens)
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
