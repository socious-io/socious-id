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

func paymentsGroup() {

	It("Should create the cards", func() {
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

	It("Should get all the cards", func() {
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

	It("Should get one card", func() {
		for i, data := range cardsData {
			w := httptest.NewRecorder()
			reqBody, _ := json.Marshal(data)
			req, _ := http.NewRequest("GET", fmt.Sprintf("/payments/cards/%s", cardsData[i]["id"]), bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", authTokens[i])
			router.ServeHTTP(w, req)

			body := decodeBody(w.Body)
			Expect(w.Code).To(Equal(http.StatusCreated))
			Expect(data["id"]).To(Equal(body["id"]))
		}
	})

	It("Should create the wallets", func() {
		for i, data := range walletsData {
			w := httptest.NewRecorder()
			reqBody, _ := json.Marshal(data)
			req, _ := http.NewRequest("POST", "/payments/wallets", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", authTokens[i])
			router.ServeHTTP(w, req)

			body := decodeBody(w.Body)
			Expect(w.Code).To(Equal(http.StatusCreated))
			walletsData[i] = body["wallet"].(gin.H)
		}
	})

	It("Should get all the wallets", func() {
		for i, data := range walletsData {
			w := httptest.NewRecorder()
			reqBody, _ := json.Marshal(data)
			req, _ := http.NewRequest("GET", "/payments/wallets", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", authTokens[i])
			router.ServeHTTP(w, req)

			body := decodeBody(w.Body)
			Expect(w.Code).To(Equal(http.StatusCreated))
			Expect(data["id"]).To(Equal(body["id"]))
		}
	})

}
