package config

import "github.com/spf13/viper"

type Config struct {
	Port string // cổng WS server cho cả backend lẫn FE kết nối vào
}

func Load() *Config {
	viper.AutomaticEnv()
	port := viper.GetString("ALERT_PORT")
	if port == "" {
		port = "8081"
	}
	return &Config{Port: port}
}
