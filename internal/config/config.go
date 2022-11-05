package config

import "github.com/spf13/viper"

type Config struct {
	DBHost     string `mapstructure:"DB_HOST"`
	DBPort     int    `mapstructure:"DB_PORT"`
	DBUser     string `mapstructure:"DB_USER"`
	DBName     string `mapstructure:"DB_NAME"`
	DBSSLMode  string `mapstructure:"DB_SSLMODE"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	SrvPort    string `mapstructure:"SRV_PORT"`
}

func New(folder string) (*Config, error) {
	cfg := new(Config)
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AddConfigPath(folder)

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	if err := viper.Unmarshal(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
