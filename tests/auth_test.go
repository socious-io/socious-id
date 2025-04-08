package tests_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"socious-id/src/apps/auth"
	"socious-id/src/apps/models"
	"socious-id/src/apps/utils"

	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func authGroup() {

	authExecuted = true
	// var sessionStore sessions.Store
	var sessionCookies []*http.Cookie
	const cookieName = "socious-id-session"
	offset := len(usersData)

	BeforeAll(func() {

		ctx := context.Background()

		//Set access config
		secret := utils.RandomString(24)
		clientID := utils.RandomString(8)
		clientSecret, _ := auth.HashPassword(secret)

		access := &models.Access{
			Name:         "test",
			Description:  "test description",
			ClientID:     clientID,
			ClientSecret: clientSecret,
		}

		if err := access.Create(ctx); err != nil {
			log.Fatal(err)
		}
		authConfig = gin.H{
			"client_id":     clientID,
			"client_secret": secret,
			"redirect_url":  "http://example.com",
		}

	})

	It("It should create session (login+register)", func() {
		for range len(usersData) * 2 {
			w := httptest.NewRecorder()
			reqBody, _ := json.Marshal(authConfig)
			req, _ := http.NewRequest("POST", "/auth/session", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, req)

			body := decodeBody(w.Body)
			bodyExpect(body["auth_session"].(map[string]interface{}),
				gin.H{
					"id":           "<ANY>",
					"redirect_url": "<ANY>",
					"access_id":    "<ANY>",
					"access":       "<ANY>",
					"expire_at":    "<ANY>",
					"verified_at":  "<ANY>",
					"updated_at":   "<ANY>",
					"created_at":   "<ANY>",
				},
			)
			Expect(w.Code).To(Equal(http.StatusCreated))
			sessionsData = append(sessionsData, body["auth_session"].(map[string]interface{}))
		}
	})

	It("It should fetch session url (register)", func() {
		for i := range usersData {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(
				"GET",
				fmt.Sprintf("/auth/session/%s?auth_mode=register", sessionsData[i]["id"]),
				nil,
			)
			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusPermanentRedirect))
			location, exists := w.Result().Header["Location"]
			Expect(exists).To(Equal(true))
			Expect(location[0]).To(Equal("/auth/register/pre"))

			for _, cookie := range w.Result().Cookies() {
				if cookie.Name == cookieName {
					sessionCookies = append(sessionCookies, cookie)
				}
			}
		}
	})

	It("should send otp for the user (register)", func() {
		for i, data := range usersData {
			w := httptest.NewRecorder()
			reqBody, _ := json.Marshal(gin.H{
				"email": data["email"],
			})
			req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(reqBody))
			req.AddCookie(sessionCookies[i])
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusSeeOther))
			location := w.Result().Header.Get("Location")
			Expect(location).To(Equal(fmt.Sprintf("/auth/otp?email=%s", data["email"])))

			//Save OTP
			otp, _ := models.GetOTPByEmail(data["email"].(string))
			authOtpCodes = append(authOtpCodes, otp.Code)
		}
	})

	It("should confirm otp code for the user registration", func() {
		for i, data := range usersData {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(
				"GET",
				fmt.Sprintf(
					"/auth/otp/confirm?email=%s&code=%s",
					data["email"],
					authOtpCodes[i],
				),
				nil,
			)
			req.AddCookie(sessionCookies[i])
			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusSeeOther))
			location := w.Result().Header.Get("Location")
			Expect(location).To(Equal("/users/profile"))

			for _, cookie := range w.Result().Cookies() {
				if cookie.Name == cookieName {
					sessionCookies[i] = cookie
				}
			}
		}
	})

	It("should fill out the update profile form to complete register", func() {
		for i, data := range usersData {
			reqBody, _ := json.Marshal(data)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("PUT", "/users/profile", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			req.AddCookie(sessionCookies[i])
			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusSeeOther))
			location := w.Result().Header.Get("Location")
			Expect(location).To(Equal("/auth/confirm"))
		}
	})

	It("should confirm the session completion (register)", func() {
		for i, data := range usersData {
			w := httptest.NewRecorder()
			reqBody, _ := json.Marshal(gin.H{
				"confirmed": true,
			})
			req, _ := http.NewRequest("POST", "/auth/confirm", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			req.AddCookie(sessionCookies[i])
			router.ServeHTTP(w, req)

			//Save OTP
			otp, _ := models.GetOTPByEmail(data["email"].(string))
			ssoOtpCodes = append(ssoOtpCodes, otp.Code)

			Expect(w.Code).To(Equal(http.StatusFound))
			location := w.Result().Header.Get("Location")
			Expect(location).To(Equal(fmt.Sprintf(
				"%s?code=%s&identity_id=&session=%s&status=success",
				authConfig["redirect_url"],
				ssoOtpCodes[i],
				sessionsData[i]["id"],
			)))
		}
	})

	It("should fetch the tokens user (register)", func() {
		for i := range usersData {
			w := httptest.NewRecorder()
			reqBody, _ := json.Marshal(gin.H{
				"client_id":     authConfig["client_id"],
				"client_secret": authConfig["client_secret"],
				"code":          ssoOtpCodes[i],
			})
			req, _ := http.NewRequest("POST", "/auth/session/token", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, req)

			body := decodeBody(w.Body)
			Expect(w.Code).To(Equal(http.StatusAccepted))
			bodyExpect(body, gin.H{
				"access_token":  "<ANY>",
				"refresh_token": "<ANY>",
				"token_type":    "Bearer",
			})
		}
	})

	It("It should fetch session url (login)", func() {
		for i := range usersData {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(
				"GET",
				fmt.Sprintf("/auth/session/%s?auth_mode=login", sessionsData[i+offset]["id"]),
				nil,
			)
			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusPermanentRedirect))
			location, exists := w.Result().Header["Location"]
			Expect(exists).To(Equal(true))
			Expect(location[0]).To(Equal("/auth/login"))

			for _, cookie := range w.Result().Cookies() {
				if cookie.Name == cookieName {
					sessionCookies = append(sessionCookies, cookie)
				}
			}
		}
	})

	It("should login user", func() {
		for i, data := range usersData {
			w := httptest.NewRecorder()
			reqBody, _ := json.Marshal(data)
			req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(reqBody))
			req.AddCookie(sessionCookies[i+offset])
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusFound))
			location, exists := w.Result().Header["Location"]
			Expect(exists).To(Equal(true))
			Expect(location[0]).To(Equal("/auth/confirm"))

			for _, cookie := range w.Result().Cookies() {
				if cookie.Name == cookieName {
					sessionCookies[i+offset] = cookie
				}
			}
		}
	})

	It("should confirm the session completion (login)", func() {
		for i, data := range usersData {
			w := httptest.NewRecorder()
			reqBody, _ := json.Marshal(gin.H{
				"confirmed": true,
			})
			req, _ := http.NewRequest("POST", "/auth/confirm", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			req.AddCookie(sessionCookies[i+offset])
			router.ServeHTTP(w, req)

			//Save OTP
			otp, _ := models.GetOTPByEmail(data["email"].(string))
			ssoOtpCodes = append(ssoOtpCodes, otp.Code)

			Expect(w.Code).To(Equal(http.StatusFound))
			location := w.Result().Header.Get("Location")
			Expect(location).To(Equal(fmt.Sprintf(
				"%s?code=%s&identity_id=&session=%s&status=success",
				authConfig["redirect_url"],
				ssoOtpCodes[i+offset],
				sessionsData[i+offset]["id"],
			)))
		}
	})

	It("should fetch the tokens user (login)", func() {
		for i := range usersData {
			w := httptest.NewRecorder()
			reqBody, _ := json.Marshal(gin.H{
				"client_id":     authConfig["client_id"],
				"client_secret": authConfig["client_secret"],
				"code":          ssoOtpCodes[i+offset],
			})
			req, _ := http.NewRequest("POST", "/auth/session/token", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, req)

			body := decodeBody(w.Body)
			Expect(w.Code).To(Equal(http.StatusAccepted))
			bodyExpect(body, gin.H{
				"access_token":  "<ANY>",
				"refresh_token": "<ANY>",
				"token_type":    "Bearer",
			})
			authTokens = append(authTokens, body["access_token"].(string))
			authRefreshTokens = append(authRefreshTokens, body["refresh_token"].(string))
			claims, _ := auth.VerifyToken(body["access_token"].(string))
			usersData[i]["id"] = claims.ID
		}
	})
}
