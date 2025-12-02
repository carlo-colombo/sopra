package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// Config holds the application's configuration.
type Config struct {
	Port int `mapstructure:"PORT"`

	OpenSkyClient struct {
		ID     string `mapstructure:"id"`
		Secret string `mapstructure:"secret"`
	} `mapstructure:"opensky_client"`
	Service struct {
		Latitude  float64 `mapstructure:"DEFAULT_LATITUDE"`
		Longitude float64 `mapstructure:"DEFAULT_LONGITUDE"`
		Radius    float64 `mapstructure:"DEFAULT_RADIUS"` // in kilometers
	} `mapstructure:"service"`
}

// LoadConfig loads configuration from file and environment variables.
func LoadConfig() (*Config, error) {
	viper.SetConfigName("config") // name of config file (without extension)
	viper.AddConfigPath(".")      // path to look for the config file in the current directory


	viper.AutomaticEnv() // read in environment variables that match

	// Explicitly bind environment variables for OpenSkyClient
	viper.BindEnv("opensky_client.id", "OPENSKY_CLIENT_ID")
	viper.BindEnv("opensky_client.secret", "OPENSKY_CLIENT_SECRET")

	// Explicitly bind environment variables for Service defaults
	viper.BindEnv("service.default_latitude", "DEFAULT_LATITUDE")
	viper.BindEnv("service.default_longitude", "DEFAULT_LONGITUDE")
	viper.BindEnv("service.default_radius", "DEFAULT_RADIUS")



	// Set default values for new fields
	viper.SetDefault("PORT", 8080)

	viper.SetDefault("OPENSKY_CLIENT_ID", "")
	viper.SetDefault("OPENSKY_CLIENT_SECRET", "")

	// Set defaults from .env specific keys for radius
	viper.SetDefault("DEFAULT_LATITUDE", 47.3769)
	viper.SetDefault("DEFAULT_LONGITUDE", 8.5417)
	viper.SetDefault("DEFAULT_RADIUS", 100.0)

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

	return &cfg, nil
}

// String provides a string representation of the Config struct (for printing)
func (c *Config) String() string {
	return fmt.Sprintf(`
Configuration:
  Port: %d
  OpenSky Client:
    ID: %s
    Secret: %s
  Service Defaults:
    Latitude: %.4f
    Longitude: %.4f
    Radius: %.2f km
`,
		c.Port,
		c.OpenSkyClient.ID, c.OpenSkyClient.Secret,
		c.Service.Latitude, c.Service.Longitude, c.Service.Radius)
}

