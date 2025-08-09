package tools

import (
	"context"
	"testing"
)

type MockConfig struct {
	sonarrURL               string
	sonarrAPIKey            string
	radarrURL               string
	radarrAPIKey            string
	showsRootPath           string
	moviesRootPath          string
	defaultQualityProfileID int
}

func (m *MockConfig) SonarrURL() string {
	return m.sonarrURL
}

func (m *MockConfig) SonarrAPIKey() string {
	return m.sonarrAPIKey
}

func (m *MockConfig) RadarrURL() string {
	return m.radarrURL
}

func (m *MockConfig) RadarrAPIKey() string {
	return m.radarrAPIKey
}

func (m *MockConfig) ShowsRootPath() string {
	return m.showsRootPath
}

func (m *MockConfig) MoviesRootPath() string {
	return m.moviesRootPath
}

func (m *MockConfig) DefaultQualityProfileID() int {
	return m.defaultQualityProfileID
}

func TestGetTools(t *testing.T) {
	cfg := &MockConfig{
		sonarrURL:               "http://sonarr.test",
		sonarrAPIKey:            "sonarr-api-key",
		radarrURL:               "http://radarr.test",
		radarrAPIKey:            "radarr-api-key",
		showsRootPath:           "/test/shows",
		moviesRootPath:          "/test/movies",
		defaultQualityProfileID: 10,
	}

	sonarrClient := &mockSonarrClient{}
	radarrClient := &mockRadarrClient{}

	mediaTools := New(cfg, sonarrClient, radarrClient)

	tools := mediaTools.Tools()

	if len(tools) != 3 {
		t.Errorf("Expected 3 tools, got %d", len(tools))
	}
}

type mockSonarrClient struct{}

func (m *mockSonarrClient) LookupSeries(ctx context.Context, name string) ([]Series, error) {
	return []Series{}, nil
}

func (m *mockSonarrClient) RequestSeriesDownload(ctx context.Context, series Series, qualityProfileID int, rootFolderPath string) error {
	return nil
}

func (m *mockSonarrClient) RequestSeriesDelete(ctx context.Context, series Series) error {
	return nil
}

func (m *mockSonarrClient) SearchSeriesByGenre(ctx context.Context, genre string, similarTo string, limit int) ([]Series, error) {
	return []Series{}, nil
}

type mockRadarrClient struct{}

func (m *mockRadarrClient) LookupMovie(ctx context.Context, name string) ([]Movie, error) {
	return []Movie{}, nil
}

func (m *mockRadarrClient) RequestMovieDownload(ctx context.Context, movie Movie, qualityProfileID int, rootFolderPath string) error {
	return nil
}

func (m *mockRadarrClient) SearchMoviesByGenre(ctx context.Context, genre string, similarTo string, limit int) ([]Movie, error) {
	return []Movie{}, nil
}
