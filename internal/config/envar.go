package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	DatabaseURL string `mapstructure:"DATABASE_URL"`
	JWTSecret   string `mapstructure:"JWT_SECRET"`
	ServerPort  string `mapstructure:"SERVER_PORT"`
	MaxConns    int32  `mapstructure:"DB_MAX_CONNS"`
	MinConns    int32  `mapstructure:"DB_MIN_CONNS"`
	MaxConnLifetime int32 `mapstructure:"DB_MAX_CONN_LIFETIME"`
	RateLimitRPS float64 `mapstructure:"RATE_LIMIT_RPS"`
	RateLimitBurst int `mapstructure:"RATE_LIMIT_BURST"`
}

var AppConfig *Config

func LoadConfig() error {
	viper.AddConfigPath(".")
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		// Don't fail if the .env file is not found
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	}

	// Set default values
	viper.SetDefault("SERVER_PORT", "8080")
	viper.SetDefault("DB_MAX_CONNS", 10)
	viper.SetDefault("DB_MIN_CONNS", 2)
	viper.SetDefault("DB_MAX_CONN_LIFETIME", 300) // 5 minutes
	viper.SetDefault("RATE_LIMIT_RPS", 1.0)
	viper.SetDefault("RATE_LIMIT_BURST", 5)

	err := viper.Unmarshal(&AppConfig)
	return err
}
