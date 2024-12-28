package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"os"
)

type Payload struct {
	Email          string `json:"email"`
	Password       string `json:"password"`
	Role           string `json:"role"`
	Referral       string `json:"referral"`
	Username       string `json:"username"`
	RecaptchaToken string `json:"recaptchaToken"`
}

func randomString(length int) string {
	source := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789!@#$%^&*()-=_+[]{}|;:'\",.<>/?"
	s := make([]byte, length)
	for i := range s {
		s[i] = source[rand.Intn(len(source))]
	}
	return string(s)
}

func createAccount(email, proxyUrl string) error {

	proxy, err := url.Parse(proxyUrl)
	if err != nil {
		return fmt.Errorf("error parsing proxy URL: %v", err)
	}
	url := "https://api.getgrass.io/auth/reguser"

	password := randomString(12)

	data := Payload{
		Email:          email,
		Password:       password,
		Role:           "seller",
		Referral:       "xZCNgNDpkddo41o",
		Username:       randomString(8),
		RecaptchaToken: "a",
	}
	payload, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("error encoding login credentials: %v", err)
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("error creating login request: %v", err)
	}

	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Origin", "chrome-extension://ilehaonighjijnmpnagapkhpcdbhclfg")
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("Sec-Ch-Ua", `"Not_A Brand";v="8", "Chromium";v="120", "Google Chrome";v="120"`)
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", `"MacOS"`)
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "none")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 11_6_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")

	transport := &http.Transport{Proxy: http.ProxyURL(proxy)}
	client := &http.Client{Transport: transport}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending login request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("login request failed with status: %v", resp.Status)
	}

	file, err := os.OpenFile("config/dist.txt", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("error opening accounts.txt: %v", err)
	}
	_, err = fmt.Fprintln(file, email+" "+password+" "+proxyUrl)
	if err != nil {
		return fmt.Errorf("error writing to accounts.txt: %v", err)
	}
	defer file.Close()

	return nil
}
