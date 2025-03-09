package tests_test

import "github.com/gin-gonic/gin"

var (
	intKey            = ""
	authTokens        = []string{}
	authOtpCodes      = []string{}
	ssoOtpCodes       = []string{}
	authRefreshTokens = []string{}
	authConfig        = gin.H{
		"client_id":     "string",
		"client_secret": "string",
		"redirect_url":  "string",
	}

	sessionsData = []gin.H{}
	usersData    = []gin.H{
		{
			"first_name": "TestName",
			"last_name":  "TestLastName",
			"username":   "test",
			"email":      "test@test.com",
			"password":   "test123456",
		},
		{
			"first_name": "TestName2",
			"last_name":  "TestLastName2",
			"username":   "test2",
			"email":      "test2@test.com",
			"password":   "test123456",
		},
	}
)
