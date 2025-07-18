package tools

import (
	"context"

	"github.com/IdoKendo/mcparr/pkg/client"
)

// SonarrClientAdapter adapts the client.SonarrClient to tools.SonarrClient.
type SonarrClientAdapter struct {
	client *client.SonarrClient
}

// NewSonarrClientAdapter creates a new SonarrClientAdapter.
func NewSonarrClientAdapter(client *client.SonarrClient) *SonarrClientAdapter {
	return &SonarrClientAdapter{client: client}
}

// LookupSeries adapts the client.SonarrClient.LookupSeries method.
func (a *SonarrClientAdapter) LookupSeries(ctx context.Context, name string) ([]Series, error) {
	clientSeries, err := a.client.LookupSeries(ctx, name)
	if err != nil {
		return nil, err
	}

	series := make([]Series, len(clientSeries))
	for i, s := range clientSeries {
		series[i] = Series{
			ID:       s.ID,
			Title:    s.Title,
			Overview: s.Overview,
			Genres:   s.Genres,
		}
	}

	return series, nil
}

// RequestSeriesDownload adapts the client.SonarrClient.RequestSeriesDownload method.
func (a *SonarrClientAdapter) RequestSeriesDownload(ctx context.Context, series Series, qualityProfileID int, rootFolderPath string) error {
	clientSeries := client.Series{
		ID:       series.ID,
		Title:    series.Title,
		Overview: series.Overview,
		Genres:   series.Genres,
	}

	return a.client.RequestSeriesDownload(ctx, clientSeries, qualityProfileID, rootFolderPath)
}

// SearchSeriesByGenre adapts the client.SonarrClient.SearchSeriesByGenre method.
func (a *SonarrClientAdapter) SearchSeriesByGenre(ctx context.Context, genre string, similarTo string, limit int) ([]Series, error) {
	clientSeries, err := a.client.SearchSeriesByGenre(ctx, genre, similarTo, limit)
	if err != nil {
		return nil, err
	}

	series := make([]Series, len(clientSeries))
	for i, s := range clientSeries {
		series[i] = Series{
			ID:       s.ID,
			Title:    s.Title,
			Overview: s.Overview,
			Genres:   s.Genres,
		}
	}

	return series, nil
}

// RadarrClientAdapter adapts the client.RadarrClient to tools.RadarrClient.
type RadarrClientAdapter struct {
	client *client.RadarrClient
}

// NewRadarrClientAdapter creates a new RadarrClientAdapter.
func NewRadarrClientAdapter(client *client.RadarrClient) *RadarrClientAdapter {
	return &RadarrClientAdapter{client: client}
}

// LookupMovie adapts the client.RadarrClient.LookupMovie method.
func (a *RadarrClientAdapter) LookupMovie(ctx context.Context, name string) ([]Movie, error) {
	clientMovies, err := a.client.LookupMovie(ctx, name)
	if err != nil {
		return nil, err
	}

	movies := make([]Movie, len(clientMovies))
	for i, m := range clientMovies {
		movies[i] = Movie{
			ID:       m.ID,
			Title:    m.Title,
			Overview: m.Overview,
			Genres:   m.Genres,
		}
	}

	return movies, nil
}

// RequestMovieDownload adapts the client.RadarrClient.RequestMovieDownload method.
func (a *RadarrClientAdapter) RequestMovieDownload(ctx context.Context, movie Movie, qualityProfileID int, rootFolderPath string) error {
	clientMovie := client.Movie{
		ID:       movie.ID,
		Title:    movie.Title,
		Overview: movie.Overview,
		Genres:   movie.Genres,
	}

	return a.client.RequestMovieDownload(ctx, clientMovie, qualityProfileID, rootFolderPath)
}

// SearchMoviesByGenre adapts the client.RadarrClient.SearchMoviesByGenre method.
func (a *RadarrClientAdapter) SearchMoviesByGenre(ctx context.Context, genre string, similarTo string, limit int) ([]Movie, error) {
	clientMovies, err := a.client.SearchMoviesByGenre(ctx, genre, similarTo, limit)
	if err != nil {
		return nil, err
	}

	movies := make([]Movie, len(clientMovies))
	for i, m := range clientMovies {
		movies[i] = Movie{
			ID:       m.ID,
			Title:    m.Title,
			Overview: m.Overview,
			Genres:   m.Genres,
		}
	}

	return movies, nil
}
