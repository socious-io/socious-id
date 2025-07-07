package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"socious-id/src/apps/utils"
	"socious-id/src/config"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AppleUserInfo struct {
	Email string `json:"email"`
	Name  struct {
		FirstName string `json:"first_name,omitempty"`
		LastName  string `json:"last_name,omitempty"`
	} `json:"name,omitempty"`
}

type AppleAccessToken struct {
	AccessToken  string `json:"access_token"`
	IDToken      string `json:"id_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

type AppleClientSecretParams struct {
	ClientID       string
	PrivateKeyPath string
	KeyID          string
	TeamID         string
}

type AppleLoginForm struct {
	Code string `json:"code" form:"code" validate:"required"`
}

func createClientSecret(p AppleClientSecretParams) (string, error) {
	keyData, err := os.ReadFile(p.PrivateKeyPath)
	if err != nil {
		return "", fmt.Errorf("failed to read private key: %w", err)
	}
	privateKey, err := jwt.ParseECPrivateKeyFromPEM(keyData)
	if err != nil {
		return "", fmt.Errorf("failed to parse private key: %w", err)
	}

	// Create client secret
	now := time.Now()
	claims := jwt.MapClaims{
		"iss": p.TeamID,
		"iat": now.Unix(),
		"exp": now.Add(6 * 30 * 24 * time.Hour).Unix(), // 6 months
		"aud": "https://appleid.apple.com",
		"sub": p.ClientID,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	token.Header["kid"] = p.KeyID
	clientSecret, err := token.SignedString(privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign JWT: %w", err)
	}

	return clientSecret, nil
}

func getAppleToken(code string) (*AppleAccessToken, error) {
	// Read and parse private key
	clientSecret, err := createClientSecret(AppleClientSecretParams{
		ClientID:       config.Config.Oauth.Apple.ID,
		KeyID:          config.Config.Oauth.Apple.KeyID,
		TeamID:         config.Config.Oauth.Apple.TeamID,
		PrivateKeyPath: config.Config.Oauth.Apple.PrivateKeyPath,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create client secret: %w", err)
	}

	req, err := http.NewRequest("POST", "https://appleid.apple.com/auth/token", strings.NewReader(url.Values{
		"client_id":     {config.Config.Oauth.Apple.ID},
		"client_secret": {clientSecret},
		"code":          {code},
		"grant_type":    {"authorization_code"},
	}.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return nil, fmt.Errorf("apple token request failed: %v", err)
	}
	defer resp.Body.Close()

	appleToken := new(AppleAccessToken)
	if err := json.NewDecoder(resp.Body).Decode(&appleToken); err != nil {
		return nil, fmt.Errorf("invalid token response: %w", err)
	}

	return appleToken, nil
}

func getUserInfo(token AppleAccessToken) (*AppleUserInfo, error) {
	claims := jwt.MapClaims{}
	_, _, err := jwt.NewParser().ParseUnverified(token.IDToken, claims)
	if err != nil {
		return nil, fmt.Errorf("failed to decode id_token: %w", err)
	}

	userInfo := new(AppleUserInfo)
	utils.Copy(claims, userInfo)

	return userInfo, nil
}

func AppleLoginWithCode(code, ref string) (*AppleUserInfo, error) {
	token, err := getAppleToken(code)
	if err != nil {
		return nil, err
	}

	return getUserInfo(*token)
}
