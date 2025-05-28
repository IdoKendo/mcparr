package main

import "os"

type Config struct {
	SonarrUrl    string
	SonarrApiKey string

	RadarrUrl    string
	RadarrApiKey string
}

func NewConfig() Config {
	return Config{
		SonarrUrl:    os.Getenv("SONARR_URL"),
		SonarrApiKey: os.Getenv("SONARR_API_KEY"),
		RadarrUrl:    os.Getenv("RADARR_URL"),
		RadarrApiKey: os.Getenv("RADARR_API_KEY"),
	}
}
