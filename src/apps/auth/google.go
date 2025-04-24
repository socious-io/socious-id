package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"socious-id/src/config"
)

type GoogleAccessToken struct {
	AccessToken string `json:"access_token"`
}
type GoogleUserInfo struct {
	Email      string `json:"email"`
	FamilyName string `json:"family_name"`
	GivenName  string `json:"given_name"`
}

func getGoogleUserInfo(accessToken string) (*GoogleUserInfo, error) {
	req, err := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v2/userinfo", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	userInfo := new(GoogleUserInfo)
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, err
	}

	return userInfo, nil
}

func getGoogleToken(code, ref string) (*GoogleAccessToken, error) {
	form := map[string]any{
		"code":          code,
		"client_id":     config.Config.Oauth.Google.ID,
		"client_secret": config.Config.Oauth.Google.Secret,
		"grant_type":    "authorization_code",
		"redirect_uri":  ref,
	}

	formBytes, err := json.Marshal(form)
	if err != nil {
		return nil, err
	}
	reader := bytes.NewReader(formBytes)

	req, err := http.NewRequest("POST", "https://oauth2.googleapis.com/token", reader)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed to get token")
	}

	accessToken := new(GoogleAccessToken)
	if err := json.NewDecoder(resp.Body).Decode(&accessToken); err != nil {
		return nil, err
	}

	return accessToken, nil
}

func GoogleLoginWithCode(code, ref string) (*GoogleUserInfo, error) {
	googleToken, err := getGoogleToken(code, ref)
	if err != nil {
		return nil, err
	}

	return getGoogleUserInfo(googleToken.AccessToken)
}
