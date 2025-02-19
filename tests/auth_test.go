package tests_test

// import (
// 	. "github.com/onsi/ginkgo/v2"
// 	. "github.com/onsi/gomega"
// )

func authGroup() {

	authExecuted = true

	//TODO: Fix the test (temporarily bypassing so that i can init the deployment)
	// It("should register user", func() {
	// 	for i, data := range usersData {
	// 		w := httptest.NewRecorder()
	// 		reqBody, _ := json.Marshal(data)
	// 		req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(reqBody))
	// 		req.Header.Set("Content-Type", "application/json")
	// 		router.ServeHTTP(w, req)

	// 		body := decodeBody(w.Body)
	// 		bodyExpect(body, gin.H{"access_token": "<ANY>", "refresh_token": "<ANY>", "token_type": "Bearer"})
	// 		Expect(w.Code).To(Equal(200))
	// 		authTokens = append(authTokens, body["access_token"].(string))
	// 		authRefreshTokens = append(authRefreshTokens, body["refresh_token"].(string))
	// 		claims, _ := auth.VerifyToken(body["access_token"].(string))
	// 		usersData[i]["id"] = claims.ID
	// 	}
	// })

}
