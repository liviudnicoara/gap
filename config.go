package gap

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

// Initilize this variable to access the env values
var EnvConfigs *Config

// Config represents the configuration for the task pool.
type Config struct {
	BaseWorkers   int           `mapstructure:"GAP_BASE_WORKERS"`
	MaxWorkers    int           `mapstructure:"GAP_MAX_WORKER"`
	WorkerTimeout time.Duration `mapstructure:"GAP_WORKER_TIMEOUT"`
}

// Call to load the variables from env
func loadEnvVariables() (config *Config) {
	viper.SetConfigFile(".env")

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
