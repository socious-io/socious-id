package tests_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func verificationsGroup() {
	It("should create verification", func() {
		for i, data := range usersData {
			w := httptest.NewRecorder()

			req, _ := http.NewRequest("POST", "/verifications", nil)
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", authTokens[i])
			router.ServeHTTP(w, req)

			body := decodeBody(w.Body)

			Expect(w.Code).To(Equal(201))
			Expect(body).To(HaveKey("id"))
			Expect(body).To(HaveKey("status"))
			Expect(body["status"]).To(Equal("pending"))
			Expect(body["user_id"]).To(Equal(data["id"]))

			verificationsData = append(verificationsData, body)
		}
	})

	It("should get verification for current user", func() {
		for i, data := range verificationsData {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/verifications", nil)
			req.Header.Set("Authorization", authTokens[i])
			router.ServeHTTP(w, req)

			body := decodeBody(w.Body)
			Expect(w.Code).To(Equal(200))
			Expect(body).To(HaveKey("id"))
			Expect(body["id"]).To(Equal(data["id"]))
			Expect(body["user_id"]).To(Equal(usersData[i]["id"]))
		}
	})

	It("should initialize verification connection", func() {

		for _, data := range verificationsData {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", fmt.Sprintf("/verifications/%s/connect", data["id"]), nil)
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

	// It("should reuse existing connection if recent", func() {
	// 	// First create a verification
	// 	createVerification := func() string {
	// 		w := httptest.NewRecorder()
	// 		data := verificationsData[0]
	// 		data["schema_id"] = schemasData[0]["id"]

	// 		for j, attr := range schemasData[0]["attributes"].([]gin.H) {
	// 			data["attributes"].([]gin.H)[j]["attribute_id"] = attr["id"]
	// 		}

	// 		reqBody, _ := json.Marshal(data)
	// 		req, _ := http.NewRequest("POST", "/verifications", bytes.NewBuffer(reqBody))
	// 		req.Header.Set("Content-Type", "application/json")
	// 		req.Header.Set("Authorization", authTokens[0])
	// 		router.ServeHTTP(w, req)

	// 		body := decodeBody(w.Body)
	// 		return body["id"].(string)
	// 	}

	// 	id := createVerification()

	// 	// Connect first time
	// 	w1 := httptest.NewRecorder()
	// 	req1, _ := http.NewRequest("GET", "/verifications/"+id+"/connect", nil)
	// 	router.ServeHTTP(w1, req1)

	// 	body1 := decodeBody(w1.Body)
	// 	firstConnectionURL := body1["connection_url"]
	// 	firstConnectionID := body1["connection_id"]

	// 	// Connect second time (should reuse connection if within 2 minutes)
	// 	w2 := httptest.NewRecorder()
	// 	req2, _ := http.NewRequest("GET", "/verifications/"+id+"/connect", nil)
	// 	router.ServeHTTP(w2, req2)

	// 	body2 := decodeBody(w2.Body)
	// 	Expect(w2.Code).To(Equal(200))
	// 	Expect(body2["connection_url"]).To(Equal(firstConnectionURL))
	// 	Expect(body2["connection_id"]).To(Equal(firstConnectionID))
	// })

	// It("should handle verification callback", func() {
	// 	// First create a verification
	// 	createVerification := func() string {
	// 		w := httptest.NewRecorder()
	// 		data := verificationsData[0]
	// 		data["schema_id"] = schemasData[0]["id"]

	// 		for j, attr := range schemasData[0]["attributes"].([]gin.H) {
	// 			data["attributes"].([]gin.H)[j]["attribute_id"] = attr["id"]
	// 		}

	// 		reqBody, _ := json.Marshal(data)
	// 		req, _ := http.NewRequest("POST", "/verifications", bytes.NewBuffer(reqBody))
	// 		req.Header.Set("Content-Type", "application/json")
	// 		req.Header.Set("Authorization", authTokens[0])
	// 		router.ServeHTTP(w, req)

	// 		body := decodeBody(w.Body)
	// 		return body["id"].(string)
	// 	}

	// 	id := createVerification()

	// 	// Initialize connection first
	// 	w1 := httptest.NewRecorder()
	// 	req1, _ := http.NewRequest("GET", "/verifications/"+id+"/connect", nil)
	// 	router.ServeHTTP(w1, req1)
	// 	Expect(w1.Code).To(Equal(200))

	// 	// Now test callback
	// 	w2 := httptest.NewRecorder()
	// 	req2, _ := http.NewRequest("GET", "/verifications/"+id+"/callback", nil)
	// 	router.ServeHTTP(w2, req2)

	// 	body := decodeBody(w2.Body)
	// 	Expect(w2.Code).To(Equal(200))
	// 	Expect(body).To(HaveKey("message"))
	// 	Expect(body["message"]).To(Equal("success"))

	// 	// Verify that verification status is updated by checking the verification
	// 	w3 := httptest.NewRecorder()
	// 	req3, _ := http.NewRequest("GET", "/verifications", nil)
	// 	req3.Header.Set("Authorization", authTokens[0])
	// 	router.ServeHTTP(w3, req3)

	// 	verificationStatus := decodeBody(w3.Body)
	// 	Expect(w3.Code).To(Equal(200))
	// 	// Status will depend on your implementation, but it should be updated
	// 	// This test assumes ProofRequest might change the status to "in_progress"
	// 	Expect(verificationStatus["status"]).NotTo(Equal("pending"))
	// })

	// It("should verify user when verification is verified", func() {
	// 	// This test requires mocking the verification service to set a verification as verified
	// 	// First create a verification
	// 	createVerification := func() string {
	// 		w := httptest.NewRecorder()
	// 		data := verificationsData[0]
	// 		data["schema_id"] = schemasData[0]["id"]

	// 		for j, attr := range schemasData[0]["attributes"].([]gin.H) {
	// 			data["attributes"].([]gin.H)[j]["attribute_id"] = attr["id"]
	// 		}

	// 		reqBody, _ := json.Marshal(data)
	// 		req, _ := http.NewRequest("POST", "/verifications", bytes.NewBuffer(reqBody))
	// 		req.Header.Set("Content-Type", "application/json")
	// 		req.Header.Set("Authorization", authTokens[0])
	// 		router.ServeHTTP(w, req)

	// 		body := decodeBody(w.Body)
	// 		return body["id"].(string)
	// 	}

	// 	id := createVerification()

	// 	// We need to manually set the verification to verified status in the test database
	// 	// This would normally happen through the verification process
	// 	// Mock this by directly updating the database
	// 	verification, _ := models.GetVerification(uuid.MustParse(id))
	// 	verification.Status = models.VerificationStatusVerified
	// 	verification.VerifiedAt = new(time.Time)
	// 	*verification.VerifiedAt = time.Now()
	// 	verification.Save()

	// 	// Now get the verification, which should trigger user verification
	// 	w := httptest.NewRecorder()
	// 	req, _ := http.NewRequest("GET", "/verifications", nil)
	// 	req.Header.Set("Authorization", authTokens[0])
	// 	router.ServeHTTP(w, req)

	// 	Expect(w.Code).To(Equal(200))

	// 	// Also check if user is now verified
	// 	user, _ := models.GetUser(uuid.MustParse(users[0]["id"].(string)))
	// 	Expect(user.IsVerified(models.UserVerificationTypeIdenity)).To(BeTrue())
	// })

	// It("should handle error if user verification fails", func() {
	// 	// This test requires mocking a verification as verified but user verification failing
	// 	// First create a verification
	// 	createVerification := func() string {
	// 		w := httptest.NewRecorder()
	// 		data := verificationsData[0]
	// 		data["schema_id"] = schemasData[0]["id"]

	// 		for j, attr := range schemasData[0]["attributes"].([]gin.H) {
	// 			data["attributes"].([]gin.H)[j]["attribute_id"] = attr["id"]
	// 		}

	// 		reqBody, _ := json.Marshal(data)
	// 		req, _ := http.NewRequest("POST", "/verifications", bytes.NewBuffer(reqBody))
	// 		req.Header.Set("Content-Type", "application/json")
	// 		req.Header.Set("Authorization", authTokens[0])
	// 		router.ServeHTTP(w, req)

	// 		body := decodeBody(w.Body)
	// 		return body["id"].(string)
	// 	}

	// 	id := createVerification()

	// 	// Set verification to verified status
	// 	verification, _ := models.GetVerification(uuid.MustParse(id))
	// 	verification.Status = models.VerificationStatusVerified
	// 	verification.VerifiedAt = new(time.Time)
	// 	*verification.VerifiedAt = time.Now()
	// 	verification.Save()

	// 	// Mock the user Verify method to fail
	// 	// This requires mocking or a test-specific implementation
	// 	// For this test, we'll assume there's a way to make the User.Verify method fail

	// 	// Get the verification, which should attempt user verification but fail
	// 	w := httptest.NewRecorder()
	// 	req, _ := http.NewRequest("GET", "/verifications", nil)
	// 	req.Header.Set("Authorization", authTokens[0])
	// 	router.ServeHTTP(w, req)

	// 	body := decodeBody(w.Body)
	// 	Expect(w.Code).To(Equal(422))
	// 	Expect(body).To(HaveKey("error"))
	// 	Expect(body["error"]).To(Equal("user is verified but couldn't verify user"))
	// })
}
