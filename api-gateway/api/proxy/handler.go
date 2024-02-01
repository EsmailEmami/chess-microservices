package proxy

import (
	"io"
	"net/http"
	"strings"

	"github.com/esmailemami/chess/api-gateway/api/config"
	"github.com/esmailemami/chess/api-gateway/api/middleware"
	"github.com/esmailemami/chess/api-gateway/api/util"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

func ProxyRoutes(router *mux.Router) error {
	config, err := config.LoadConfiguration()

	if err != nil {
		return err
	}

	for _, proxyConfig := range config.Proxies {
		route := router.PathPrefix(proxyConfig.Path).Subrouter()

		// proxy middlewares
		route.Use(loadMiddlewares(proxyConfig.Middlewares...)...)

		for _, routeConfig := range proxyConfig.Routes {
			handler := chainMiddlewares(handleRequest(routeConfig, proxyConfig.Path, proxyConfig.Target), routeConfig.Middlewares...)

			if strings.HasSuffix(routeConfig.Path, "*") {
				route.PathPrefix(strings.TrimSuffix(routeConfig.Path, "*")).HandlerFunc(handler).Methods(routeConfig.Methods...)
			} else {
				route.HandleFunc(routeConfig.Path, handler).Methods(routeConfig.Methods...)
			}
		}
	}

	return nil
}

func handleRequest(routeConfig config.RouteConfig, basePath, target string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		proxyRequest(basePath, target, w, r)
	}
}

func proxyRequest(basePath, target string, w http.ResponseWriter, r *http.Request) {
	path := target + strings.TrimPrefix(r.URL.RequestURI(), basePath)

	if util.IsWebSocketRequest(r) {
		proxyWebSocket(cleanWebSocketAddr(path), w, r)
	} else {
		proxyHTTP(path, w, r)
	}
}

func proxyHTTP(path string, w http.ResponseWriter, r *http.Request) {
	req, err := http.NewRequest(r.Method, path, r.Body)
	if err != nil {
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}

	req.Header = r.Header

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to proxy request", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}
	w.WriteHeader(resp.StatusCode)
	_, _ = io.Copy(w, resp.Body)
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func proxyWebSocket(target string, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	header := http.Header{}
	header.Add("UserId", r.Header.Get("UserId"))

	targetConn, _, err := websocket.DefaultDialer.Dial(target, header)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	defer targetConn.Close()

	go forwardWebSocket(targetConn, conn)
	forwardWebSocket(conn, targetConn)
}

func forwardWebSocket(src, dest *websocket.Conn) {
	for {
		messageType, p, err := src.ReadMessage()
		if err != nil {
			break
		}
		err = dest.WriteMessage(messageType, p)
		if err != nil {
			break
		}
	}
}

func cleanWebSocketAddr(path string) string {
	path = strings.Replace(path, "http://", "ws://", 1)
	path = strings.Replace(path, "https://", "wss://", 1)
	return path
}

func loadMiddlewares(middlewaresName ...string) []mux.MiddlewareFunc {
	var middlewares = make([]mux.MiddlewareFunc, 0, len(middlewaresName))

	for _, middlewareName := range middlewaresName {
		switch middlewareName {
		case "authentication":
			middlewares = append(middlewares, middleware.Authenticate)
		case "execute-duration":
			middlewares = append(middlewares, middleware.ExecuteDuration)
		}
	}

	return middlewares
}

func chainMiddlewares(f http.HandlerFunc, middlewareNames ...string) http.HandlerFunc {
	middlewares := loadMiddlewares(middlewareNames...)

	for _, middleware := range middlewares {
		f = middleware(http.HandlerFunc(f)).ServeHTTP
	}

	return http.HandlerFunc(f)
}
