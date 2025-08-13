package tests_test

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx/types"
)

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
		"policies":      []string{},
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
			"website":     "https://my.website.org",
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
			"website":     "https://my.website2.org",
			"mission":     "TestMission2",
			"culture":     "TestCulture2",
		},
	}

	impactPointsData = []gin.H{
		{
			"type":                  "DONATION",
			"total_points":          100,
			"social_cause":          "HEALTH",
			"social_cause_category": "HEALTH",
			"meta": gin.H{
				"test": "test",
			},
			"value":      0.0,
			"unique_tag": "unique_1",
		},
		{
			"type":                  "DONATION",
			"total_points":          200,
			"social_cause":          "SOCIAL",
			"social_cause_category": "LIFE",
			"meta": gin.H{
				"test": "test",
			},
			"value":      1.0,
			"unique_tag": "unique_2",
		},
	}

	credentialsData = []gin.H{}

	cardsData = []gin.H{
		{
			"token": "test1",
		},
		{
			"token": "test2",
		},
	}
	walletsData = []gin.H{
		{
			"chain":    "Cardano",
			"chain_id": "chain_id_1",
			"address":  "0xexample",
		},
		{
			"chain":    "Ethereum",
			"chain_id": "chain_id_2",
			"address":  "0xexample1",
		},
	}

	referralAchievementsData = []gin.H{
		{
			"referee_id":       "auto:<referee_id>",
			"achievement_type": "VOTE",
			"meta":             types.JSONText(`{"meta_key": "test"}`),
		},
		{
			"referee_id":       "auto:<referee_id>",
			"achievement_type": "REF_KYC",
			"meta":             types.JSONText(`{"meta_key": "test2"}`),
		},
	}

	shorteningURLs = []string{
		"https://google.com",
		"https://app.socious.io",
	}
)
