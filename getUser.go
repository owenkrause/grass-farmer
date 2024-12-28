package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type UserData struct {
	Status string `json:"status"`
	Data   struct {
		TotalEarningsForUser  int `json:"totalEarningsForUser"`
		TotalReferralEarnings int `json:"totalReferralEarnings"`
	} `json:"data"`
}

type Device struct {
	ID         string  `json:"id"`
	DeviceIP   string  `json:"device_ip"`
	Country    string  `json:"country_code"`
	Earning    float64 `json:"earning"`
	DeviceType string  `json:"device_type"`
}

func GetUser(token string, proxy *url.URL) (UserData, error) {
	url := "https://api.getgrass.io/users/dash"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return UserData{}, fmt.Errorf("error creating user info request: %v", err)
	}

	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Cookie", "token="+token)
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
		return UserData{}, fmt.Errorf("error sending user info request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return UserData{}, fmt.Errorf("user info request failed with status: %v", resp.Status)
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return UserData{}, fmt.Errorf("error reading response body: %v", err)
	}

	var response UserData
	if err := json.Unmarshal(responseBody, &response); err != nil {
		return UserData{}, fmt.Errorf("error unmarshalling JSON: %v", err)
	}

	return response, nil
}
