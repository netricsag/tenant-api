package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

var (
	CLIENT_ID     string
	CLIENT_SECRET string
)

func GetGithubAccessToken(code string) string {
	requestBodyMap := map[string]string{"client_id": CLIENT_ID, "client_secret": CLIENT_SECRET, "code": code}
	requestJSON, _ := json.Marshal(requestBodyMap)

	req, reqerr := http.NewRequest("POST", "https://github.com/login/oauth/access_token", bytes.NewBuffer(requestJSON))
	if reqerr != nil {
		ErrorLogger.Printf("Request creation failed: %s", reqerr)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, respErr := http.DefaultClient.Do(req)
	if respErr != nil {
		ErrorLogger.Printf("Request failed: %s", respErr)
	}

	respbody, _ := ioutil.ReadAll(resp.Body)

	type GithubAccessTokenResponse struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		Scope       string `json:"scope"`
	}

	var githubAccessTokenResponse GithubAccessTokenResponse
	json.Unmarshal(respbody, &githubAccessTokenResponse)

	return githubAccessTokenResponse.AccessToken
}

func GetGithubData(accessToken string) string {
	req, reqerr := http.NewRequest("GET", "https://api.github.com/user", nil)
	if reqerr != nil {
		ErrorLogger.Printf("Request creation failed: %s", reqerr)
	}

	authorizationHeaderValue := fmt.Sprintf("token %s", accessToken)
	req.Header.Set("Authorization", authorizationHeaderValue)

	resp, respErr := http.DefaultClient.Do(req)
	if respErr != nil {
		ErrorLogger.Printf("Request failed: %s", respErr)
	}

	respbody, _ := ioutil.ReadAll(resp.Body)

	return string(respbody)
}
