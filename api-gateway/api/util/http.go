package util

import (
	"net/http"
	"strings"
)

func IsWebSocketRequest(r *http.Request) bool {
	return strings.ToLower(r.Header.Get("Upgrade")) == "websocket" && strings.ToLower(r.Header.Get("Connection")) == "upgrade"
}
