package config

import (
	"encoding/json"
	"os"
	"sort"

	"github.com/spf13/viper"
)

type MiddlewareConfig struct {
	Name   string                 `json:"name"`
	Type   string                 `json:"type"`
	Config map[string]interface{} `json:"config"`
}

type RouteConfig struct {
	Path        string   `json:"path"`
	Methods     []string `json:"methods"`
	Middlewares []string `json:"middlewares"`
}

type ProxyConfig struct {
	Path        string        `json:"path"`
	Target      string        `json:"target"`
	Middlewares []string      `json:"middlewares"`
	Routes      []RouteConfig `json:"routes"`
}

type GatewayConfig struct {
	Proxies []ProxyConfig `json:"proxies"`
}

func LoadConfiguration() (*GatewayConfig, error) {
	file, err := os.Open(viper.GetString("app.gateway_file_path"))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	config := &GatewayConfig{}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(config)
	if err != nil {
		return nil, err
	}

	for _, proxy := range config.Proxies {
		sortRoutesByPathLength(proxy.Routes)
	}

	return config, nil
}

func sortRoutesByPathLength(routes []RouteConfig) {
	sort.Slice(routes, func(i, j int) bool {
		return len(routes[i].Path) > len(routes[j].Path)
	})
}
