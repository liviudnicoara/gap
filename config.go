package gap

import (
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

	defaultConfig := Config{
		BaseWorkers:   100,
		MaxWorkers:    1000,
		WorkerTimeout: time.Second,
	}

	if err := viper.ReadInConfig(); err != nil {
		config = &defaultConfig
	}

	if err := viper.Unmarshal(&config); err != nil {
		config = &defaultConfig
	}

	return
}
