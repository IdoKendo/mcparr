package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

// RadarrClient is a client for interacting with the Radarr API.
type RadarrClient struct {
	client *Client
}

// NewRadarrClient creates a new Radarr API client.
func NewRadarrClient(baseURL, apiKey string) *RadarrClient {
	return &RadarrClient{
		client: NewClient(baseURL, apiKey),
	}
}

// LookupMovie searches for movies by name.
func (r *RadarrClient) LookupMovie(ctx context.Context, name string) ([]Movie, error) {
	term := url.QueryEscape(name)

	params := map[string]string{"term": term}
	data, err := r.client.Get(ctx, "movie/lookup", params)
	if err != nil {
		return nil, fmt.Errorf("failed to lookup movie: %w", err)
	}

	var movies []Movie
	if err := json.Unmarshal(data, &movies); err != nil {
		return nil, fmt.Errorf("failed to parse movie response: %w", err)
	}

	return movies, nil
}

// RequestMovieDownload requests a movie to be downloaded.
func (r *RadarrClient) RequestMovieDownload(ctx context.Context, movie Movie, qualityProfileID int, rootFolderPath string) error {
	data := map[string]any{
		"title":            movie.Title,
		"tmdbId":           movie.ID,
		"qualityProfileId": qualityProfileID,
		"rootFolderPath":   rootFolderPath,
	}

	_, err := r.client.Post(ctx, "movie", data)
	if err != nil {
		return fmt.Errorf("failed to request movie download: %w", err)
	}

	return nil
}

// SearchMoviesByGenre searches for movies by genre.
func (r *RadarrClient) SearchMoviesByGenre(ctx context.Context, genre string, similarTo string, limit int) ([]Movie, error) {
	var allMovies []Movie

	if similarTo != "" {
		movies, err := r.LookupMovie(ctx, similarTo)
		if err != nil {
			return nil, err
		}
		allMovies = append(allMovies, movies...)
	} else {
		popularMovies := []string{"Inception", "The Shawshank Redemption", "The Dark Knight", "Pulp Fiction", "Forrest Gump"}
		for _, movie := range popularMovies {
			movies, err := r.LookupMovie(ctx, movie)
			if err != nil {
				continue
			}
			allMovies = append(allMovies, movies...)

			if len(allMovies) >= 20 {
				break
			}
		}
	}

	var matchingMovies []Movie
	genreLower := strings.ToLower(genre)

	for _, m := range allMovies {
		for _, g := range m.Genres {
			if strings.ToLower(g) == genreLower {
				matchingMovies = append(matchingMovies, m)
				break
			}
		}
		if len(matchingMovies) >= limit {
			break
		}
	}

	return matchingMovies, nil
}
