package consul

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/esmailemami/chess/shared/logging"
	"github.com/hashicorp/consul/api"
	"github.com/spf13/viper"
)

func Register() {
	if !viper.GetBool("consul.enable") {
		logging.Info("Consul is not enabled in project")
		return
	}

	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		logging.FatalE("Error creating Consul client", err)
	}

	// Register the service with Consul
	registerService(client)

	// Handle OS signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	// Deregister the service from Consul
	deregisterService(client)
	os.Exit(1)
}

func registerService(client *api.Client) {
	var (
		ttl         time.Duration = time.Duration(viper.GetInt("consul.ttl")) * time.Second
		checkID                   = viper.GetString("consul.check_id")
		serviceID                 = viper.GetString("consul.id")
		clusterName               = viper.GetString("consul.cluster_name")
		tags                      = viper.GetStringSlice("consul.tags")
		address                   = viper.GetString("app.address")
		port                      = viper.GetInt("app.port")
	)

	check := &api.AgentServiceCheck{
		DeregisterCriticalServiceAfter: ttl.String(),
		TLSSkipVerify:                  true,
		TTL:                            ttl.String(),
		CheckID:                        checkID,
	}

	registration := &api.AgentServiceRegistration{
		Check:   check,
		ID:      serviceID,
		Name:    clusterName,
		Tags:    tags,
		Address: address,
		Port:    port,
	}

	err := client.Agent().ServiceRegister(registration)
	if err != nil {
		logging.FatalE("Error registering service with Consul", err)
	}

	// Health check loop
	updateHealthCheck(client)

	logging.Info("Service registered with Consul")
}

func updateHealthCheck(client *api.Client) {
	var (
		ttl     time.Duration = time.Duration(viper.GetInt("consul.check_duration")) * time.Second
		checkID               = viper.GetString("consul.check_id")
		ticker                = time.NewTicker(ttl)
	)
	go func() {

		for {
			client.Agent().UpdateTTL(checkID, "application is running", api.HealthPassing)
			<-ticker.C
		}
	}()
}

func deregisterService(client *api.Client) {
	serviceID := viper.GetString("consul.id")

	err := client.Agent().ServiceDeregister(serviceID)
	if err != nil {
		logging.ErrorE("Error deregistering service from Consul:", err)
	} else {
		logging.Info("Service deregistered from Consul")
	}
}
