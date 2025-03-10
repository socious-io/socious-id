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

	organizationsData = []gin.H{
		{
			"shortname":   "TestOrgShortName",
			"name":        "TestOrgName",
			"bio":         "TestBio",
			"description": "ThisIsTestDescription",
			"email":       "test-org@test.com",
			"phone":       "+1111111111",
			"city":        "TestCity",
			"country":     "TestCountry",
			"address":     "TestAddress",
			"website":     "TestWebsite",
			"mission":     "TestMission",
			"culture":     "TestCulture",
		},
		{
			"shortname":   "TestOrgShortName2",
			"name":        "TestOrgName2",
			"bio":         "TestBio2",
			"description": "ThisIsTestDescription2",
			"email":       "test-org2@test.com",
			"phone":       "+22222222222",
			"city":        "TestCity2",
			"country":     "TestCountry2",
			"address":     "TestAddress2",
			"website":     "TestWebsite2",
			"mission":     "TestMission2",
			"culture":     "TestCulture2",
		},
	}
)
