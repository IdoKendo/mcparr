package config

import (
	"fmt"
	"log"
	"os"
)

// Config holds the application configuration including API endpoints and keys.
type Config struct {
	sonarrURL               string
	sonarrAPIKey            string
	radarrURL               string
	radarrAPIKey            string
	showsRootPath           string
	moviesRootPath          string
	defaultQualityProfileID int
}

// New creates a new Config with values from environment variables.
func New() *Config {
	sonarrApiKey, sonarrExists := os.LookupEnv("SONARR_API_KEY")
	radarrApiKey, radarrExists := os.LookupEnv("RADARR_API_KEY")

	if !sonarrExists || !radarrExists {
		log.Fatal("Missing SONARR_API_KEY and/or RADARR_API_KEY in env")
	}

	return &Config{
		sonarrAPIKey:            sonarrApiKey,
		radarrAPIKey:            radarrApiKey,
		sonarrURL:               envWithDefault("SONARR_URL", "http://localhost:8989"),
		radarrURL:               envWithDefault("RADARR_URL", "http://localhost:7878"),
		showsRootPath:           envWithDefault("SHOWS_ROOT_PATH", "/media/library/shows"),
		moviesRootPath:          envWithDefault("MOVIES_ROOT_PATH", "/media/library/movies"),
		defaultQualityProfileID: envIntWithDefault("DEFAULT_QUALITY_PROFILE_ID", 6),
	}
}

// SonarrURL returns the Sonarr URL.
func (c *Config) SonarrURL() string {
	return c.sonarrURL
}

// SonarrAPIKey returns the Sonarr API key.
func (c *Config) SonarrAPIKey() string {
	return c.sonarrAPIKey
}

// RadarrURL returns the Radarr URL.
func (c *Config) RadarrURL() string {
	return c.radarrURL
}

// RadarrAPIKey returns the Radarr API key.
func (c *Config) RadarrAPIKey() string {
	return c.radarrAPIKey
}

// ShowsRootPath returns the root path for TV shows.
func (c *Config) ShowsRootPath() string {
	return c.showsRootPath
}

// MoviesRootPath returns the root path for movies.
func (c *Config) MoviesRootPath() string {
	return c.moviesRootPath
}

// DefaultQualityProfileID returns the default quality profile ID.
func (c *Config) DefaultQualityProfileID() int {
	return c.defaultQualityProfileID
}

func envWithDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func envIntWithDefault(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	var intValue int
	_, err := fmt.Sscanf(value, "%d", &intValue)
	if err != nil {
		return defaultValue
	}

	return intValue
}
