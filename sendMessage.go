package main

import (
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
)

func SendMessage(conn *websocket.Conn, message interface{}) error {

	payload, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("error encoding JSON: %v", err)
	}

	err = conn.WriteMessage(websocket.TextMessage, payload)
	if err != nil {
		return fmt.Errorf("error writing message to WebSocket: %v", err)
	}

	return nil
}
