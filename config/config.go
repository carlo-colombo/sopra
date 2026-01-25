package config

import (
	"fmt"
	"log"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Config holds the application's configuration.
type Config struct {
	Print    bool   `mapstructure:"print"`
	Watch    bool   `mapstructure:"watch"`
	Interval int    `mapstructure:"interval"`
	Port     int    `mapstructure:"port"`
	DBPath   string `mapstructure:"db_path"`

	OpenSkyClient struct {
		ID     string `mapstructure:"id"`
		Secret string `mapstructure:"secret"`
	} `mapstructure:"opensky_client"`
	FlightAware struct {
		APIKey string `mapstructure:"api_key"`
	} `mapstructure:"flightaware"`
	Service struct {
		Latitude  float64 `mapstructure:"latitude"`
		Longitude float64 `mapstructure:"longitude"`
		Radius    float64 `mapstructure:"radius"`
	} `mapstructure:"service"`
}

// LoadConfig loads configuration from file and environment variables.
func LoadConfig(path string) (*Config, error) {
	viper.SetConfigName("config")
	viper.AddConfigPath(path)

	viper.AutomaticEnv()

	// Bind environment variables

	_ = viper.BindPFlag("print", pflag.Lookup("print"))
	_ = viper.BindPFlag("watch", pflag.Lookup("watch"))
	_ = viper.BindPFlag("interval", pflag.Lookup("interval"))
	if err := viper.BindEnv("port", "PORT"); err != nil {
		log.Fatalf("failed to bind 'port' env: %v", err)
	}

	if err := viper.BindEnv("opensky_client.id", "OPENSKY_CLIENT_ID"); err != nil {
		log.Fatalf("failed to bind 'opensky_client.id' env: %v", err)
	}

	if err := viper.BindEnv("opensky_client.secret", "OPENSKY_CLIENT_SECRET"); err != nil {
		log.Fatalf("failed to bind 'opensky_client.secret' env: %v", err)
	}
	if err := viper.BindEnv("flightaware.api_key", "FLIGHTAWARE_API_KEY"); err != nil {
		log.Fatalf("failed to bind 'flightaware.api_key' env: %v", err)
	}
	if err := viper.BindEnv("service.latitude", "DEFAULT_LATITUDE"); err != nil {
		log.Fatalf("failed to bind 'service.latitude' env: %v", err)
	}
	if err := viper.BindEnv("service.longitude", "DEFAULT_LONGITUDE"); err != nil {
		log.Fatalf("failed to bind 'service.longitude' env: %v", err)
	}
	if err := viper.BindEnv("service.radius", "DEFAULT_RADIUS"); err != nil {
		log.Fatalf("failed to bind 'service.radius' env: %v", err)
	}
	if err := viper.BindEnv("watch", "WATCH"); err != nil {
		log.Fatalf("failed to bind 'watch' env: %v", err)
	}
	if err := viper.BindEnv("interval", "WATCH_INTERVAL"); err != nil {
		log.Fatalf("failed to bind 'interval' env: %v", err)
	}
	if err := viper.BindEnv("db_path", "DB_PATH"); err != nil {
		log.Fatalf("failed to bind 'db_path' env: %v", err)
	}

	// Set default values

	viper.SetDefault("port", 8080)
	viper.SetDefault("watch", false)
	viper.SetDefault("interval", 300)
	viper.SetDefault("db_path", "sopra.db")

	viper.SetDefault("opensky_client.id", "")

	viper.SetDefault("opensky_client.secret", "")

	viper.SetDefault("flightaware.api_key", "")

	viper.SetDefault("service.latitude", 47.3769)

	viper.SetDefault("service.longitude", 8.5417)

	viper.SetDefault("service.radius", 100.0)

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
	  Print: %t
	  Watch: %t
	  Interval: %ds
	  Port: %d
	  DB Path: %s

	  OpenSky Client:

	    ID: %s

	    Secret: %s

	  FlightAware Client:

	    API Key: %s

	  Service Defaults:

	    Latitude: %.4f

	    Longitude: %.4f

	    Radius: %.2f km

	`,
		c.Print,
		c.Watch,
		c.Interval,
		c.Port,
		c.DBPath,

		c.OpenSkyClient.ID, c.OpenSkyClient.Secret,

		c.FlightAware.APIKey,

		c.Service.Latitude, c.Service.Longitude, c.Service.Radius)

}
