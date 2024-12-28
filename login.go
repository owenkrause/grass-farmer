package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
)

type Credentials struct {
	User     string `json:"user"`
	Password string `json:"password"`
}

func Login(email, password string, proxy *url.URL) ([]byte, string, error) {
	credentials := Credentials{User: email, Password: password}
	credentialsJSON, err := json.Marshal(credentials)
	if err != nil {
		return nil, "", fmt.Errorf("error encoding login credentials: %v", err)
	}

	url := "https://api.getgrass.io/auth/login"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(credentialsJSON))
	if err != nil {
		return nil, "", fmt.Errorf("error creating login request: %v", err)
	}

	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Origin", "https://app.getgrass.io")
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("Referer", "https://app.getgrass.io/")
	req.Header.Set("Sec-Ch-Ua", `"Not_A Brand";v="8", "Chromium";v="120", "Google Chrome";v="120"`)
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", `"macOS"`)
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-site")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")

	transport := &http.Transport{Proxy: http.ProxyURL(proxy)}
	client := &http.Client{Transport: transport}

	resp, err := client.Do(req)
	if err != nil {
		return nil, "", fmt.Errorf("error sending login request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("login request failed with status: %v", resp.Status)
	}
	cookie := resp.Header["Set-Cookie"][0]
	re := regexp.MustCompile(`token=([^;]+);`)
	authToken := re.FindStringSubmatch(cookie)[1]

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", fmt.Errorf("error reading response body: %v", err)
	}

	return responseBody, authToken, nil
}
