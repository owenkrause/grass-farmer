package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"

	"github.com/gorilla/websocket"
)

func OpenWebSocketConnection(targetURL string, proxy *url.URL) (*websocket.Conn, error) {
	dialer := websocket.Dialer{
		Proxy: http.ProxyURL(proxy),
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	headers := http.Header{}
	headers.Set("Accept-Encoding", "gzip, deflate, br")
	headers.Set("Accept-Language", "en-US,en;q=0.9")
	headers.Set("Cache-Control", "no-cache")
	headers.Set("Host", "proxy.wynd.network:4650")
	headers.Set("Origin", "https://app.getgrass.io")
	headers.Set("Pragma", "no-cache")
	headers.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")

	conn, _, err := dialer.Dial(targetURL, headers)
	if err != nil {
		return nil, fmt.Errorf("error during WebSocket handshake: %v", err)
	}
	return conn, nil
}
