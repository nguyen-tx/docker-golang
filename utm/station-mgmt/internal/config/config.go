package config

import "github.com/spf13/viper"

type Config struct {
	App            AppConfig
	NetworkStation ServiceConfig
	RTKStation     ServiceConfig
	ManageInfo     ServiceConfig
}

type AppConfig struct {
	Port     string `mapstructure:"APP_PORT"`
	GRPCPort string `mapstructure:"GRPC_PORT"`
	Env      string `mapstructure:"APP_ENV"`
}

// ServiceConfig giữ địa chỉ gRPC của service con
type ServiceConfig struct {
	Address string
}

func Load() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	_ = viper.ReadInConfig()

	return &Config{
		App: AppConfig{
			Port:     viper.GetString("APP_PORT"),
			GRPCPort: viper.GetString("GRPC_PORT"),
			Env:      viper.GetString("APP_ENV"),
		},
		NetworkStation: ServiceConfig{
			Address: viper.GetString("NETWORK_STATION_GRPC_ADDR"), // e.g. network-station:9090
		},
		RTKStation: ServiceConfig{
			Address: viper.GetString("RTK_STATION_GRPC_ADDR"), // e.g. rtk-station:9090
		},
		ManageInfo: ServiceConfig{
			Address: viper.GetString("MANAGE_INFO_GRPC_ADDR"), // e.g. manage-info:9090
		},
	}, nil
}
