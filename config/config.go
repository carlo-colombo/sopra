package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// Config holds the application's configuration.
type Config struct {
	OpenSkyClient struct {
		ID     string `mapstructure:"id"`
		Secret string `mapstructure:"secret"`
	} `mapstructure:"opensky_client"`
	Service struct {
		Latitude  float64 `mapstructure:"latitude"`
		Longitude float64 `mapstructure:"longitude"`
		Radius    float64 `mapstructure:"radius"` // in kilometers
	} `mapstructure:"service"`
}

// LoadConfig loads configuration from file and environment variables.
func LoadConfig() (*Config, error) {
	viper.SetConfigName("config") // name of config file (without extension)
	viper.AddConfigPath(".")      // path to look for the config file in the current directory
	viper.SetConfigFile(".env")   // look for .env file

	viper.AutomaticEnv() // read in environment variables that match

	// Set default values
	viper.SetDefault("service.latitude", 47.3769)  // Zurich
	viper.SetDefault("service.longitude", 8.5417) // Zurich
	viper.SetDefault("service.radius", 100.0)      // 100km

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if not critical
			fmt.Println("No config file found, using environment variables and defaults")
		} else {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Override with environment variables if set
	if id := viper.GetString("OPENREDISKY_CLIENT_ID"); id != "" {
		cfg.OpenSkyClient.ID = id
	}
	if secret := viper.GetString("OPENREDISKY_CLIENT_SECRET"); secret != "" {
		cfg.OpenSkyClient.Secret = secret
	}

	return &cfg, nil
}
