package tests_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func organizationsGroup() {

	It("Should create organizations", func() {
		for i, data := range organizationsData {
			w := httptest.NewRecorder()
			reqBody, _ := json.Marshal(data)
			req, _ := http.NewRequest("POST", "/organizations", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", authTokens[0])
			router.ServeHTTP(w, req)

			body := decodeBody(w.Body)
			Expect(w.Code).To(Equal(http.StatusCreated))
			organizationsData[i]["id"] = body["id"]
		}
	})

	It("Should get all organizations", func() {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/organizations", nil)
		req.Header.Set("Authorization", authTokens[0])
		router.ServeHTTP(w, req)

		fmt.Println(w.Body)

		body := decodeBody(w.Body)
		Expect(w.Code).To(Equal(http.StatusOK))
		bodyExpect(body, gin.H{
			"page":    "<ANY>",
			"limit":   "<ANY>",
			"results": "<ANY>",
			"total":   "<ANY>",
		})
	})

	It("Should get all of the organizations that i am a member of", func() {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/organizations/my", nil)
		req.Header.Set("Authorization", authTokens[0])
		router.ServeHTTP(w, req)

		fmt.Println(w.Body)

		body := decodeBody(w.Body)
		Expect(w.Code).To(Equal(http.StatusOK))
		bodyExpect(body, gin.H{
			"page":    "<ANY>",
			"limit":   "<ANY>",
			"results": "<ANY>",
			"total":   "<ANY>",
		})
	})

	It("Should get a organization", func() {
		for _, data := range organizationsData {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", fmt.Sprintf("/organizations/%s", data["id"]), nil)
			req.Header.Set("Authorization", authTokens[0])
			router.ServeHTTP(w, req)

			body := decodeBody(w.Body)
			Expect(w.Code).To(Equal(http.StatusOK))
			Expect(body["id"]).To(Equal(data["id"]))
		}
	})

	It("Should update organizations", func() {
		for i, data := range organizationsData {
			w := httptest.NewRecorder()

			//Changing data
			name := fmt.Sprintf("name_%d", i+1)
			data["name"] = name

			reqBody, _ := json.Marshal(data)
			req, _ := http.NewRequest("PUT", fmt.Sprintf("/organizations/%s", data["id"]), bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", authTokens[0])
			router.ServeHTTP(w, req)

			body := decodeBody(w.Body)
			Expect(w.Code).To(Equal(http.StatusAccepted))
			Expect(body["name"]).To(Equal(name))
		}
	})

	It("Should remove a organization", func() {
		data := organizationsData[1]

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", fmt.Sprintf("/organizations/%s", data["id"]), nil)
		req.Header.Set("Authorization", authTokens[0])
		router.ServeHTTP(w, req)

		Expect(w.Code).To(Equal(http.StatusOK))
		organizationsData = organizationsData[:len(organizationsData)-1]
	})

	It("Should add a member to organization", func() {
		for _, data := range organizationsData {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", fmt.Sprintf("/organizations/%s/members/%s", data["id"], usersData[1]["id"]), nil)
			req.Header.Set("Authorization", authTokens[0])
			router.ServeHTTP(w, req)

			body := decodeBody(w.Body)
			Expect(w.Code).To(Equal(http.StatusOK))
			Expect(body["id"]).To(Equal(data["id"]))
		}
	})

	It("Should remove a member from organization", func() {
		for _, data := range organizationsData {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("DELETE", fmt.Sprintf("/organizations/%s/members/%s", data["id"], usersData[1]["id"]), nil)
			req.Header.Set("Authorization", authTokens[0])
			router.ServeHTTP(w, req)

			body := decodeBody(w.Body)
			Expect(w.Code).To(Equal(http.StatusOK))
			Expect(body["id"]).To(Equal(data["id"]))
		}
	})

}
