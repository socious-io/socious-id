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

func usersGroup() {

	It("Should get user", func() {
		for i, data := range usersData {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/users", nil)
			req.Header.Set("Authorization", authTokens[i])
			router.ServeHTTP(w, req)

			body := decodeBody(w.Body)
			Expect(w.Code).To(Equal(http.StatusOK))
			Expect(body["id"]).To(Equal(data["id"]))
		}
	})

	It("Should update user", func() {
		for i, data := range usersData {
			w := httptest.NewRecorder()

			//Changing data
			lastName := fmt.Sprintf("test_last_%d", i+1)
			data["last_name"] = lastName

			reqBody, _ := json.Marshal(data)
			req, _ := http.NewRequest("PUT", "/users", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", authTokens[i])
			router.ServeHTTP(w, req)

			body := decodeBody(w.Body)
			Expect(w.Code).To(Equal(http.StatusOK))
			Expect(body["last_name"].(string)).To(Equal(lastName))
		}
	})

}
