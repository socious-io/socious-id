package tests_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"socious-id/src/apps/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func credentialsGroup() { //TODO: write tests for badges
	It("should create KYC credentials for current user", func() {
		for i, data := range usersData {
			w := httptest.NewRecorder()

			reqBody, _ := json.Marshal(gin.H{
				"type": models.CredentialTypeKYC,
			})
			req, _ := http.NewRequest("POST", "/credentials", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", authTokens[i])
			router.ServeHTTP(w, req)

			body := decodeBody(w.Body)

			Expect(w.Code).To(Equal(201))
			Expect(body).To(HaveKey("id"))
			Expect(body).To(HaveKey("status"))
			Expect(body["status"]).To(Equal(string(models.CredentialStatusCreated)))
			Expect(body["user_id"]).To(Equal(data["id"]))

			credentialsData = append(credentialsData, body)
		}
	})

	It("should get KYC credential for current user", func() {
		for i, data := range credentialsData {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/credentials?type=KYC", nil)
			req.Header.Set("Authorization", authTokens[i])
			router.ServeHTTP(w, req)

			body := decodeBody(w.Body)
			Expect(w.Code).To(Equal(200))
			Expect(body).To(HaveKey("id"))
			Expect(body["id"]).To(Equal(data["id"]))
			Expect(body["user_id"]).To(Equal(usersData[i]["id"]))
		}
	})

	//TODO: Should deeply mock the wallet lib
	PIt("should initialize credential connection", func() {
		for _, data := range credentialsData {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", fmt.Sprintf("/credentials/%s/connect", data["id"]), nil)
			router.ServeHTTP(w, req)

			body := decodeBody(w.Body)
			Expect(w.Code).To(Equal(200))
			Expect(body).To(HaveKey("id"))
			Expect(body["id"]).To(Equal(data["id"]))
			Expect(body).To(HaveKey("connection_url"))
			Expect(body["connection_url"]).NotTo(BeNil())
			Expect(body).To(HaveKey("connection_id"))
			Expect(body["connection_id"]).NotTo(BeNil())
			Expect(body).To(HaveKey("connection_at"))
			Expect(body["connection_at"]).NotTo(BeNil())
		}
	})

	//TODO: Should deeply mock the wallet lib
	PIt("should reuse existing connection if recent", func() {
		for _, data := range credentialsData {
			w1 := httptest.NewRecorder()
			req1, _ := http.NewRequest("GET", fmt.Sprintf("/credentials/%s/connect", data["id"]), nil)
			router.ServeHTTP(w1, req1)

			body1 := decodeBody(w1.Body)
			firstConnectionURL := body1["connection_url"]
			firstConnectionID := body1["connection_id"]

			// Connect second time (should reuse connection if within 2 minutes)
			w2 := httptest.NewRecorder()
			req2, _ := http.NewRequest("GET", fmt.Sprintf("/credentials/%s/connect", data["id"]), nil)
			router.ServeHTTP(w2, req2)

			body2 := decodeBody(w2.Body)
			Expect(w2.Code).To(Equal(200))
			Expect(body2["connection_url"]).To(Equal(firstConnectionURL))
			Expect(body2["connection_id"]).To(Equal(firstConnectionID))
		}
	})

	//TODO: Should deeply mock the wallet lib
	PIt("should handle credentials callback", func() {
		for _, data := range credentialsData {
			// Initialize connection first
			w1 := httptest.NewRecorder()
			req1, _ := http.NewRequest("GET", fmt.Sprintf("/credentials/%s/connect", data["id"]), nil)
			router.ServeHTTP(w1, req1)
			Expect(w1.Code).To(Equal(200))

			// Now test callback
			w2 := httptest.NewRecorder()
			req2, _ := http.NewRequest("GET", fmt.Sprintf("/credentials/%s/connect", data["id"]), nil)
			router.ServeHTTP(w2, req2)

			body := decodeBody(w2.Body)
			Expect(w2.Code).To(Equal(200))
			Expect(body).To(HaveKey("message"))
			Expect(body["message"]).To(Equal("success"))

			// Verify that credentials status is updated by checking the credentials
			w3 := httptest.NewRecorder()
			req3, _ := http.NewRequest("GET", "/credentials", nil)
			req3.Header.Set("Authorization", authTokens[0])
			router.ServeHTTP(w3, req3)

			credentialStatus := decodeBody(w3.Body)
			Expect(w3.Code).To(Equal(200))
			// Status will depend on your implementation, but it should be updated
			// This test assumes ProofRequest might change the status to "in_progress"
			Expect(credentialStatus["status"]).NotTo(Equal("pending"))
		}
	})

	//TODO: Should deeply mock the wallet lib
	PIt("should verify user when credential is verified", func() {
		// This test requires mocking the credential service to set a credential as verified
		for _, data := range credentialsData {
			// We need to manually set the credential to verified status in the test database
			// This would normally happen through the credential process
			// Mock this by directly updating the database

			credentialId := uuid.MustParse(data["id"].(string))
			userId := uuid.MustParse(usersData[0]["id"].(string))

			credential, _ := models.GetCredential(credentialId)
			credential.Status = models.CredentialStatusVerified
			credential.VerifiedAt = new(time.Time)
			*credential.VerifiedAt = time.Now()
			// credential.Save()

			// Now get the credential, which should trigger user credential
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/credentials", nil)
			req.Header.Set("Authorization", authTokens[0])
			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(200))

			// Also check if user is now verified
			user, _ := models.GetUser(userId)
			Expect(user.IdentityVerifiedAt).NotTo(BeNil())
		}
	})

	//TODO: Should deeply mock the wallet lib
	PIt("should handle error if user credential fails", func() {
		// This test requires mocking a credential as verified but user credential failing
		for _, data := range credentialsData {
			// Set credential to verified status
			credentialId := uuid.MustParse(data["id"].(string))

			credential, _ := models.GetCredential(credentialId)
			credential.Status = models.CredentialStatusVerified
			credential.VerifiedAt = new(time.Time)
			*credential.VerifiedAt = time.Now()
			// credential.Save()

			// Mock the user Verify method to fail
			// This requires mocking or a test-specific implementation
			// For this test, we'll assume there's a way to make the User.Verify method fail

			// Get the credential, which should attempt user credential but fail
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/credentials", nil)
			req.Header.Set("Authorization", authTokens[0])
			router.ServeHTTP(w, req)

			body := decodeBody(w.Body)
			Expect(w.Code).To(Equal(422))
			Expect(body).To(HaveKey("error"))
			Expect(body["error"]).To(Equal("user is verified but couldn't verify user"))
		}
	})
}
