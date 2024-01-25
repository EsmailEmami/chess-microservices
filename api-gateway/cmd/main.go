package main

import (
	"log"
	"net/http"

	"github.com/esmailemami/chess/api-gateway/api/proxy"
	"github.com/esmailemami/chess/api-gateway/internal/service"
	"github.com/esmailemami/chess/shared/consul"
	"github.com/esmailemami/chess/shared/logging"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

func main() {
	// auth grpc connection
	conn, err := service.GetAuthGrpcConnection()
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	service.InitializeServices(conn)

	// app listener
	router := mux.NewRouter()

	err = proxy.ProxyRoutes(router)
	if err != nil {
		log.Fatal(err)
	}

	// register consul
	go consul.Register()

	logging.Info("gateway started")
	log.Fatal(http.ListenAndServe(":"+viper.GetString("app.port"), router))
}

func init() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AddConfigPath("../configs/app")
	viper.AddConfigPath("./api-gateway/configs/app")

	viper.SetConfigType("yaml")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		logging.FatalE("Unable to read config file.", err)
	}
}
