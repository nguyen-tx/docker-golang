package config

import "github.com/spf13/viper"

type Config struct {
	App          AppConfig
	Postgres     PostgresConfig
	Mongo        MongoConfig
	JWT          JWTConfig
	MQ           MQConfig
	MQTT         MQTTConfig
	AlertService AlertServiceConfig
}

type AppConfig struct {
	Port     string `mapstructure:"APP_PORT"`
	GRPCPort string `mapstructure:"GRPC_PORT"`
	Env      string `mapstructure:"APP_ENV"`
}

type PostgresConfig struct {
	Host     string `mapstructure:"POSTGRES_HOST"`
	Port     string `mapstructure:"POSTGRES_PORT"`
	User     string `mapstructure:"POSTGRES_USER"`
	Password string `mapstructure:"POSTGRES_PASSWORD"`
	DBName   string `mapstructure:"POSTGRES_DB"`
	SSLMode  string `mapstructure:"POSTGRES_SSL_MODE"`
}

type MongoConfig struct {
	URI    string `mapstructure:"MONGO_URI"`
	DBName string `mapstructure:"MONGO_DB"`
}

type JWTConfig struct {
	Secret      string `mapstructure:"JWT_SECRET"`
	ExpiryHours int    `mapstructure:"JWT_EXPIRY_HOURS"`
	RefreshDays int    `mapstructure:"JWT_REFRESH_DAYS"`
}

type MQConfig struct {
	URI string `mapstructure:"RABBITMQ_URI"`
}

type MQTTConfig struct {
	Broker   string `mapstructure:"MQTT_BROKER"`
	ClientID string `mapstructure:"MQTT_CLIENT_ID"`
	Username string `mapstructure:"MQTT_USERNAME"`
	Password string `mapstructure:"MQTT_PASSWORD"`
}

type AlertServiceConfig struct {
	URL string `mapstructure:"ALERT_SERVICE_URL"` // ws://utm-alert:8081/ws/backend
}

func Load() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	_ = viper.ReadInConfig() // fallback to env vars if .env not found

	return &Config{
		App: AppConfig{
			Port:     viper.GetString("APP_PORT"),
			GRPCPort: viper.GetString("GRPC_PORT"),
			Env:      viper.GetString("APP_ENV"),
		},
		Postgres: PostgresConfig{
			Host:     viper.GetString("POSTGRES_HOST"),
			Port:     viper.GetString("POSTGRES_PORT"),
			User:     viper.GetString("POSTGRES_USER"),
			Password: viper.GetString("POSTGRES_PASSWORD"),
			DBName:   viper.GetString("POSTGRES_DB"),
			SSLMode:  viper.GetString("POSTGRES_SSL_MODE"),
		},
		Mongo: MongoConfig{
			URI:    viper.GetString("MONGO_URI"),
			DBName: viper.GetString("MONGO_DB"),
		},
		JWT: JWTConfig{
			Secret:      viper.GetString("JWT_SECRET"),
			ExpiryHours: viper.GetInt("JWT_EXPIRY_HOURS"),
			RefreshDays: viper.GetInt("JWT_REFRESH_DAYS"),
		},
		MQ: MQConfig{
			URI: viper.GetString("RABBITMQ_URI"),
		},
		MQTT: MQTTConfig{
			Broker:   viper.GetString("MQTT_BROKER"),
			ClientID: viper.GetString("MQTT_CLIENT_ID"),
			Username: viper.GetString("MQTT_USERNAME"),
			Password: viper.GetString("MQTT_PASSWORD"),
		},
		AlertService: AlertServiceConfig{
			URL: viper.GetString("ALERT_SERVICE_URL"),
		},
	}, nil
}
