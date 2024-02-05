package cmd

import (
	"os"

	"github.com/esmailemami/chess/shared/database/psql"
	"github.com/esmailemami/chess/shared/database/redis"
	"github.com/esmailemami/chess/shared/logging"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use: "game",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	initConfig()
	initDB()
	initCache()
}

func initConfig() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./configs")
	viper.AddConfigPath("./game-app/configs")

	viper.SetConfigType("yaml")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		logging.FatalE("Unable to read config file.", err)
	}
}

func initDB() {
	user := viper.GetString("database.username")
	password := viper.GetString("database.password")
	host := viper.GetString("database.host")
	dbName := viper.GetString("database.name")
	port := viper.GetString("database.port")
	sslmode := viper.GetString("database.sslmode")

	// initialize db
	if err := psql.Initialize(user, password, host, dbName, port, sslmode, psql.DefaultConfig); err != nil {
		logging.FatalE("Unable to initialize database.", err)
	}

}

func initCache() {
	var (
		host     = viper.GetString("redis.host")
		port     = viper.GetString("redis.port")
		db       = viper.GetInt("redis.db")
		password = viper.GetString("redis.password")
	)

	redis.Connect(host, port, db, password)
}
