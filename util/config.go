package util

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

// Initilize this variable to access the env values
var EnvConfigs *Config

// Config represents the configuration for the task pool.
type Config struct {
	BaseWorkers            int           `mapstructure:"GOAPPPOOL_BASE_WORKERS"`
	MaxWorkers             int           `mapstructure:"GOAPPPOOL_MAX_WORKER"`
	WorkerTimeoutInSeconds time.Duration `mapstructure:"GOAPPPOOL_WORKER_TIMEOUT_SECONDS"`
}

// We will call this in main.go to load the env variables
func InitEnvConfigs() {
	EnvConfigs = loadEnvVariables()
}

// Call to load the variables from env
func loadEnvVariables() (config *Config) {
	// Tell viper the path/location of your env file. If it is root just add "."
	viper.AddConfigPath(".")

	// Tell viper the name of your file
	viper.SetConfigName("app")

	// Tell viper the type of your file
	viper.SetConfigType("env")

	// Viper reads all the variables from env file and log error if any found
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Error reading env file", err)
	}

	// Viper unmarshals the loaded env varialbes into the struct
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatal(err)
	}
	return
}
