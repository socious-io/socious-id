package tests_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	database "github.com/socious-io/pkg_database"
)

func referralsGroup() {
	BeforeAll(func() {
		database.GetDB().DB.Exec(`UPDATE users SET referred_by = $1 WHERE id = $2`, usersData[0]["id"], usersData[1]["id"])
	})

	It("Should create referral achievement", func() {
		for i, data := range referralAchievementsData {
			data["client_id"] = access.ClientID
			data["client_secret"] = secret
			data["referee_id"] = usersData[1]["id"]

			w := httptest.NewRecorder()
			reqBody, _ := json.Marshal(data)
			req, _ := http.NewRequest("POST", "/referrals/achievements", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, req)
			body := decodeBody(w.Body)

			Expect(w.Code).To(Equal(http.StatusCreated))
			referralAchievementsData[i]["id"] = body["id"]
		}
	})

	It("Should fetch referral achievement list", func() {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/referrals/achievements", nil)
		req.Header.Set("Authorization", authTokens[0])
		router.ServeHTTP(w, req)
		body := decodeBody(w.Body)
		Expect(w.Code).To(Equal(http.StatusOK))
		Expect(body["total_count"].(float64)).To(Equal(1.0))
	})

	It("Should fetch referrals", func() {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/referrals", nil)
		req.Header.Set("Authorization", authTokens[0])
		router.ServeHTTP(w, req)
		body := decodeBody(w.Body)
		Expect(w.Code).To(Equal(http.StatusOK))
		Expect(body["total_count"].(float64)).To(Equal(1.0))
	})

	It("Should fetch referral stats", func() {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/referrals/stats", nil)
		req.Header.Set("Authorization", authTokens[0])
		router.ServeHTTP(w, req)
		body := decodeBody(w.Body)
		Expect(w.Code).To(Equal(http.StatusOK))
		Expect(body["total_unclaimed_reward_amount"].(float64)).To(Equal(0.0))
		Expect(body["total_reward_amount"].(float64)).To(Equal(0.0))
		Expect(body["total_count"].(float64)).To(Equal(1.0))
	})

}
