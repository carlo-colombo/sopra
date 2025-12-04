package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestLoadConfig_Defaults(t *testing.T) {
	viper.Reset()

	os.Unsetenv("PORT")
	os.Unsetenv("OPENSKY_CLIENT_ID")
	os.Unsetenv("OPENSKY_CLIENT_SECRET")
	os.Unsetenv("FLIGHTAWARE_API_KEY")
	os.Unsetenv("DEFAULT_LATITUDE")
	os.Unsetenv("DEFAULT_LONGITUDE")
	os.Unsetenv("DEFAULT_RADIUS")

	tmpdir, err := os.MkdirTemp("", "config-test-defaults")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpdir)
	cfg, err := LoadConfig(tmpdir)

	assert.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.Equal(t, 8080, cfg.Port)
	assert.Equal(t, "", cfg.OpenSkyClient.ID)
	assert.Equal(t, "", cfg.OpenSkyClient.Secret)
	assert.Equal(t, 47.3769, cfg.Service.Latitude)
	assert.Equal(t, 8.5417, cfg.Service.Longitude)
	assert.Equal(t, 100.0, cfg.Service.Radius)
}

func TestLoadConfig_Env(t *testing.T) {
	viper.Reset()
	os.Setenv("PORT", "9090")
	os.Setenv("OPENSKY_CLIENT_ID", "env_id")
	os.Setenv("OPENSKY_CLIENT_SECRET", "env_secret")
	os.Setenv("DEFAULT_LATITUDE", "1.2345")
	os.Setenv("DEFAULT_LONGITUDE", "5.4321")
	os.Setenv("DEFAULT_RADIUS", "50.5")
	defer os.Unsetenv("PORT")
	defer os.Unsetenv("OPENSKY_CLIENT_ID")
	defer os.Unsetenv("OPENSKY_CLIENT_SECRET")
	defer os.Unsetenv("DEFAULT_LATITUDE")
	defer os.Unsetenv("DEFAULT_LONGITUDE")
	defer os.Unsetenv("DEFAULT_RADIUS")

	cfg, err := LoadConfig(".")

	assert.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.Equal(t, 9090, cfg.Port)
	assert.Equal(t, "env_id", cfg.OpenSkyClient.ID)
	assert.Equal(t, "env_secret", cfg.OpenSkyClient.Secret)
	assert.Equal(t, 1.2345, cfg.Service.Latitude)
	assert.Equal(t, 5.4321, cfg.Service.Longitude)
	assert.Equal(t, 50.5, cfg.Service.Radius)
}

func TestLoadConfig_File(t *testing.T) {
	viper.Reset()
	configContent := `
port: 7070
opensky_client:
  id: "file_id"
  secret: "file_secret"
service:
  latitude: 10.0
  longitude: 20.0
  radius: 30.0
`
	tmpdir, err := os.MkdirTemp("", "config-test")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpdir)

	tmpfile := filepath.Join(tmpdir, "config.yml")
	err = os.WriteFile(tmpfile, []byte(configContent), 0644)
	assert.NoError(t, err)

	cfg, err := LoadConfig(tmpdir)

	assert.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.Equal(t, 7070, cfg.Port)
	assert.Equal(t, "file_id", cfg.OpenSkyClient.ID)
	assert.Equal(t, "file_secret", cfg.OpenSkyClient.Secret)
	assert.Equal(t, 10.0, cfg.Service.Latitude)
	assert.Equal(t, 20.0, cfg.Service.Longitude)
	assert.Equal(t, 30.0, cfg.Service.Radius)
}
