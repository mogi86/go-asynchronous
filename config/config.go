package config

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
	"sync"
)

var once sync.Once
var c Config

// Config is configuration of environment variables
type Config struct {
	QueueName string `envconfig:"QUEUE_NAME" required:"true"`
}

func GetConfig() *Config {
	once.Do(func() {
		err := envconfig.Process("", &c)
		if err != nil {
			fmt.Printf("failed to populates the specified struct based on environment variables: %+v\n", err)
		}
	})

	return &c
}
