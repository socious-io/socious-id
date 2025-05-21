package tests_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func paymentsGroup() {

	//TODO: Should deeply mock the stripe lib
	PIt("Should create the cards", func() {
		for i, data := range cardsData {
			w := httptest.NewRecorder()
			reqBody, _ := json.Marshal(data)
			req, _ := http.NewRequest("POST", "/payments/cards", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", authTokens[i])
			router.ServeHTTP(w, req)

			body := decodeBody(w.Body)
			Expect(w.Code).To(Equal(http.StatusCreated))
			cardsData[i] = body["card"].(gin.H)
		}
	})

	//TODO: Should deeply mock the stripe lib
	PIt("Should get all the cards", func() {
		for i, data := range cardsData {
			w := httptest.NewRecorder()
			reqBody, _ := json.Marshal(data)
			req, _ := http.NewRequest("GET", "/payments/cards", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", authTokens[i])
			router.ServeHTTP(w, req)

			body := decodeBody(w.Body)
			Expect(w.Code).To(Equal(http.StatusCreated))
			Expect(data["id"]).To(Equal(body["id"]))
		}
	})

	//TODO: Should deeply mock the stripe lib
	PIt("Should delete a the card", func() {
		for i, data := range cardsData {
			w := httptest.NewRecorder()
			reqBody, _ := json.Marshal(data)
			req, _ := http.NewRequest("DELETE", "/payments/cards", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", authTokens[i])
			router.ServeHTTP(w, req)

			body := decodeBody(w.Body)
			Expect(w.Code).To(Equal(http.StatusCreated))
			Expect(data["id"]).To(Equal(body["id"]))
		}
	})

	It("Should create the wallets", func() {
		for _, data := range walletsData {
			w := httptest.NewRecorder()
			reqBody, _ := json.Marshal(data)
			req, _ := http.NewRequest("POST", "/payments/wallets", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", authTokens[0])
			router.ServeHTTP(w, req)

			body := decodeBody(w.Body)
			Expect(w.Code).To(Equal(http.StatusCreated))
			data["id"] = body["id"]
		}
	})

	It("Should get all the wallets", func() {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/payments/wallets", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", authTokens[0])
		router.ServeHTTP(w, req)

		body := []gin.H{}
		decoder := json.NewDecoder(w.Body)
		decoder.Decode(&body)

		Expect(w.Code).To(Equal(http.StatusOK))
		for i, data := range walletsData {
			Expect(body[i]["id"], data["id"])
		}
	})

}
