package proxy

import (
	"io"
	"net/http"
	"strings"

	"github.com/esmailemami/chess/api-gateway/api/config"
	"github.com/esmailemami/chess/api-gateway/api/middleware"
	"github.com/gorilla/mux"
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
	path := target + strings.TrimPrefix(r.URL.Path, basePath)

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
