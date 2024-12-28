package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Response struct {
	Status       string `json:"status"`
	OriginAction string `json:"origin_action"`
	Data         struct {
		ID uuid.UUID `json:"id"`
	} `json:"data"`
}

type Result struct {
	BrowserID  uuid.UUID `json:"browser_id"`
	DeviceType string    `json:"device_type"`
	Timestamp  int64     `json:"timestamp"`
	UserAgent  string    `json:"user_agent"`
	UserID     uuid.UUID `json:"user_id"`
	Version    string    `json:"version"`
}

type AuthMessage struct {
	ID           uuid.UUID `json:"id"`
	OriginAction string    `json:"origin_action"`
	Result       Result    `json:"result"`
}

type PingMessage struct {
	Action  string                 `json:"action"`
	Data    map[string]interface{} `json:"data"`
	ID      uuid.UUID              `json:"id"`
	Version string                 `json:"version"`
}

type PongMessage struct {
	ID           uuid.UUID `json:"id"`
	OriginAction string    `json:"origin_action"`
}

func createTask(email, password, proxyURL string, targetURL string, wg *sync.WaitGroup, quit chan struct{}) {

	defer wg.Done()

	proxy, err := url.Parse(proxyURL)
	if err != nil {
		log.Fatal("Error parsing proxy URL: ", err)
	}

	var loginResponse []byte
	// var authToken string
	for {
		loginResponse, _, err = Login(email, password, proxy)
		if err != nil {
			fmt.Printf("[Error] Login unsuccessful %v: %v\n", email, err)
			time.Sleep(3 * time.Second)
			continue
		}
		break
	}

	var conn *websocket.Conn
	var initalMessage []byte
	for {
		conn, err = OpenWebSocketConnection("wss://proxy.wynd.network:4444/", proxy)
		if err != nil {
			fmt.Printf("[Error] Open websocket connection %v\n", err)
			continue
		}
		_, initalMessage, err = conn.ReadMessage()
		if err != nil {
			fmt.Printf("[Error] Receiving auth message %v\n", err)
			conn.Close()
			time.Sleep(time.Second)
			continue
		}
		break
	}

	var authResponse AuthMessage
	if err := json.Unmarshal(initalMessage, &authResponse); err != nil {
		fmt.Printf("[Error] Unmarshaling JSON: %v\n", err)
	}

	var loginResponseData Response
	if err := json.Unmarshal(loginResponse, &loginResponseData); err != nil {
		fmt.Printf("[Error] Unmarshaling JSON: %v\n", err)
	} else {
		authMessage := AuthMessage{
			ID:           authResponse.ID,
			OriginAction: "AUTH",
			Result: Result{
				BrowserID:  uuid.New(),
				DeviceType: "extension",
				Timestamp:  time.Now().Unix(),
				UserAgent:  "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
				UserID:     loginResponseData.Data.ID,
				Version:    "3.1.4",
			},
		}
		err := SendMessage(conn, authMessage)
		if err != nil {
			fmt.Printf("[Error] Sending Auth message: %v\n", err)
		}
	}
	ticker := time.NewTicker(20 * time.Second)
	for {
		select {
		case <-ticker.C:
			pingMessage := PingMessage{
				Action:  "PING",
				Data:    make(map[string]interface{}),
				ID:      uuid.New(),
				Version: "1.0.0",
			}
			err := SendMessage(conn, pingMessage)
			if err != nil {
				fmt.Printf("[Error] Sending ping message: %v\n", err)
				createTask(email, password, proxyURL, targetURL, wg, quit)
				return
			} else {
				_, response, _ := conn.ReadMessage()
				var responseData PongMessage
				if err := json.Unmarshal(response, &responseData); err != nil {
					fmt.Printf("[Error] Unmarshaling JSON: %v\n", err)
				} else {
					message := PongMessage{
						ID:           responseData.ID,
						OriginAction: "PONG",
					}
					err := SendMessage(conn, message)
					if err != nil {
						fmt.Printf("[Error] Sending pong message: %v\n", err)
					}
				}
			}
		case <-quit:
			ticker.Stop()
			return
		}
	}
}
