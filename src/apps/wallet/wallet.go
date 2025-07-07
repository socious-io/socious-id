package wallet

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"socious-id/src/apps/shortener"
	"socious-id/src/apps/utils"
	"socious-id/src/config"
	"strings"
	"time"
)

type H map[string]interface{}

type Connect struct {
	ID      string
	URL     string
	ShortID string
}

func CreateConnection(callback string) (*Connect, error) {
	res, status, err := makeRequest("/cloud-agent/connections", "POST", H{"label": "Socious ID Connect"})
	if err != nil {
		return nil, err
	}
	var body H
	if err := json.Unmarshal(res, &body); err != nil {
		return nil, err
	}

	fmt.Println(status, reflect.TypeOf(status), status >= 400)
	if status >= 400 {
		if body["detail"] == nil {
			return nil, errors.New("cannot create connection on agent")
		}
		return nil, errors.New(body["detail"].(string))
	}

	url := strings.ReplaceAll(
		body["invitation"].(map[string]interface{})["invitationUrl"].(string),
		"https://my.domain.com/path",
		config.Config.Wallet.Connect,
	)
	url += fmt.Sprintf("&callback=%s", callback)
	c := &Connect{
		ID:  body["connectionId"].(string),
		URL: url,
	}

	short, err := shortener.New(c.URL)
	if err != nil {
		return nil, err
	}
	c.ShortID = short.ShortID
	return c, nil
}

func ProofRequest(connectionID string, challenge string) (string, error) {
	time.Sleep(time.Second)
	res, _, err := makeRequest("/cloud-agent/present-proof/presentations", "POST", H{
		"connectionId": connectionID,
		"proofs":       []H{},
		"options": H{
			"challenge": challenge,
			"domain":    "socious.io",
		},
	})
	if err != nil {
		return "", err
	}
	var body H

	if err := json.Unmarshal(res, &body); err != nil {
		return "", err
	}
	return body["presentationId"].(string), nil
}

func SendCredentials(connectionID, issuingDID string, claims H) (H, error) {
	time.Sleep(time.Second)

	res, _, err := makeRequest("/cloud-agent/issue-credentials/credential-offers", "POST", H{
		"connectionId":      connectionID,
		"claims":            claims,
		"issuingDID":        issuingDID,
		"schemaId":          nil,
		"automaticIssuance": true,
	})
	if err != nil {
		return H{}, err
	}
	var body H

	if err := json.Unmarshal(res, &body); err != nil {
		return H{}, err
	}
	return body, nil
}

func ProofVerify(presentID string) (H, error) {
	path := fmt.Sprintf("/cloud-agent/present-proof/presentations/%s", presentID)
	res, err := getRequest(path)
	if err != nil {
		return nil, err
	}
	var body H
	if err := json.Unmarshal(res, &body); err != nil {
		return nil, err
	}

	if body["status"].(string) != "PresentationVerified" {
		return nil, fmt.Errorf("presentation not verified")
	}
	_, payload, err := utils.DecodeJWT(body["data"].([]interface{})[0].(string))
	if err != nil {
		return nil, fmt.Errorf("presentation could not decode data")
	}
	var data H
	if err := json.Unmarshal(payload, &data); err != nil {
		return nil, err
	}

	_, payload, err = utils.DecodeJWT(data["vp"].(map[string]interface{})["verifiableCredential"].([]interface{})[0].(string))
	if err != nil {
		return nil, fmt.Errorf("presentation could not decode vc")
	}
	var vc H
	if err := json.Unmarshal(payload, &vc); err != nil {
		return nil, err
	}
	return vc["vc"].(map[string]interface{})["credentialSubject"].(map[string]interface{}), nil
}

func makeRequest(path string, method string, body H) ([]byte, int, error) {
	client := &http.Client{}
	url := fmt.Sprintf("%s%s", config.Config.Wallet.Agent, path)
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, 0, err
	}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, 0, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("apikey", config.Config.Wallet.AgentApiKey)
	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()
	respBody := new(bytes.Buffer)
	respBody.ReadFrom(resp.Body)
	return respBody.Bytes(), resp.StatusCode, nil
}

func getRequest(path string) ([]byte, error) {
	client := &http.Client{}
	url := fmt.Sprintf("%s%s?t=%d", config.Config.Wallet.Agent, path, time.Now().Unix())
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("apikey", config.Config.Wallet.AgentApiKey)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody := new(bytes.Buffer)
	respBody.ReadFrom(resp.Body)
	return respBody.Bytes(), nil
}
