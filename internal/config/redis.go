package config

import "time"

type TimeoutConfig struct {
	Connect time.Duration `json:"connect"`
	Read    time.Duration `json:"read"`
	Write   time.Duration `json:"write"`
	Idle    time.Duration `json:"idle"`
}

func RedisTimeouts() TimeoutConfig {
	return TimeoutConfig{
		Connect: 5 * time.Second,
		Read:    10 * time.Second,
		Write:   10 * time.Second,
		Idle:    60 * time.Second,
	}
}

func RedisHost() string {
	return cfg.RedisHost
}

func RedisPort() string {
	return cfg.RedisPort
}

func RedisPassword() string {
	return cfg.RedisPassword
}

func RedisDB() int {
	return cfg.RedisDB
}
