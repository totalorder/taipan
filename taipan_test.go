package taipan

import (
	"os"
	"testing"
)
import "github.com/stretchr/testify/require"

func TestNoProfile(t *testing.T) {
	// Given
	// Default config
	reset()

	// When
	// Config loaded
	config := Get()

	// Then
	// Default config present
	require.Equal(t, 8080, config.GetInt("port"))
	require.Equal(t, "usr", config.GetString("database.username"))
	require.Equal(t, "pwd", config.GetString("database.password"))
}

func TestStagingProfile(t *testing.T) {
	// Given
	// Staging profile
	reset()
	os.Setenv("TAIPAN_PROFILES", "staging")

	// When
	// Config loaded
	config := Get()

	// Then
	// Staging config present
	require.Equal(t, 8080, config.GetInt("port"))
	require.Equal(t, "usr", config.GetString("database.username"))
	require.Equal(t, "pwd-staging", config.GetString("database.password"))
}

func TestProductionProfile(t *testing.T) {
	// Given
	// Production profile
	reset()
	os.Setenv("TAIPAN_PROFILES", "production")

	// When
	// Config loaded
	config := Get()

	// Then
	// Production config present
	require.Equal(t, 80, config.GetInt("port"))
	require.Equal(t, "usr-readonly", config.GetString("database.username"))
	require.Equal(t, "pwd-production", config.GetString("database.password"))
}

func TestStagingProductionProfile(t *testing.T) {
	// Given
	// Staging, production profiles
	reset()
	os.Setenv("TAIPAN_PROFILES", "staging,production")

	// When
	// Config loaded
	config := Get()

	// Then
	// Production config present
	require.Equal(t, 80, config.GetInt("port"))
	require.Equal(t, "usr-readonly", config.GetString("database.username"))
	require.Equal(t, "pwd-production", config.GetString("database.password"))
}

func TestProductionStagingProfile(t *testing.T) {
	// Given
	// Staging, production profiles
	reset()
	os.Setenv("TAIPAN_PROFILES", "production,staging")

	// When
	// Config loaded
	config := Get()

	// Then
	// Production config with staging password present
	require.Equal(t, 80, config.GetInt("port"))
	require.Equal(t, "usr-readonly", config.GetString("database.username"))
	require.Equal(t, "pwd-staging", config.GetString("database.password"))
}

func TestNonExistentProfile(t *testing.T) {
	// Given
	// Non existent profile
	reset()
	os.Setenv("TAIPAN_PROFILES", "non-existent")

	// When
	// Config loaded

	// Then
	// Panic!
	require.Panics(t, func() {
		Get()
	})
}

func TestNonExistentDefaultConfig(t *testing.T) {
	// Given
	// Config path with no default config
	reset()
	os.Setenv("TAIPAN_CONFIG_PATH", "empty-folder")

	// When
	// Config loaded

	// Then
	// Panic!
	require.Panics(t, func() {
		Get()
	})
}

func TestLocalConfig(t *testing.T) {
	// Given
	// Default config with existing config-local.yml
	reset()

	// When
	// Config loaded
	config := Get()

	// Then
	// Value from local config present
	require.Equal(t, "some-value", config.GetString("local-config"))
}

func reset() {
	globalConfig = nil
	os.Unsetenv("TAIPAN_CONFIG_PATH")
	os.Unsetenv("TAIPAN_PROFILES")
}
