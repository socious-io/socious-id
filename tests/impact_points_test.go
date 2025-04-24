package tests_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func impactPointsGroup() {
	It("Should create impact history", func() {
		for i, data := range impactPointsData {
			w := httptest.NewRecorder()
			data["client_id"] = access.ClientID
			data["client_secret"] = secret
			data["user_id"] = usersData[0]["id"]
			reqBody, _ := json.Marshal(data)
			req, _ := http.NewRequest("POST", "/impact-points", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, req)
			body := decodeBody(w.Body)
			Expect(w.Code).To(Equal(http.StatusCreated))
			impactPointsData[i]["id"] = body["id"]
			fmt.Println(body)
		}
	})

	It("Should fetch impact points list", func() {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/impact-points", nil)
		req.Header.Set("Authorization", authTokens[0])
		router.ServeHTTP(w, req)
		body := decodeBody(w.Body)
		Expect(w.Code).To(Equal(http.StatusOK))
		Expect(body["total_count"].(float64)).To(Equal(2.0))
	})

	It("Should get badges", func() {
		for i := range authTokens {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/impact-points/badges", nil)
			req.Header.Set("Authorization", authTokens[i])
			router.ServeHTTP(w, req)

			body := decodeBody(w.Body)
			Expect(w.Code).To(Equal(http.StatusOK))
			Expect(body).To(HaveKey("badges"))
		}
	})
}
