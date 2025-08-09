package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

// SonarrClient is a client for interacting with the Sonarr API.
type SonarrClient struct {
	client *Client
}

// NewSonarrClient creates a new Sonarr API client.
func NewSonarrClient(baseURL, apiKey string) *SonarrClient {
	return &SonarrClient{
		client: NewClient(baseURL, apiKey),
	}
}

// LookupSeries searches for series by name.
func (s *SonarrClient) LookupSeries(ctx context.Context, name string) ([]Series, error) {
	term := url.QueryEscape(name)

	params := map[string]string{"term": term}
	data, err := s.client.Get(ctx, "series/lookup", params)
	if err != nil {
		return nil, fmt.Errorf("failed to lookup series: %w", err)
	}

	var series []Series
	if err := json.Unmarshal(data, &series); err != nil {
		return nil, fmt.Errorf("failed to parse series response: %w", err)
	}

	return series, nil
}

// RequestSeriesDownload requests a series to be downloaded.
func (s *SonarrClient) RequestSeriesDownload(ctx context.Context, series Series, qualityProfileID int, rootFolderPath string) error {
	data := map[string]any{
		"title":            series.Title,
		"tvdbId":           series.ID,
		"qualityProfileId": qualityProfileID,
		"rootFolderPath":   rootFolderPath,
	}

	_, err := s.client.Post(ctx, "series", data)
	if err != nil {
		return fmt.Errorf("failed to request series download: %w", err)
	}

	return nil
}

// SearchSeriesByGenre searches for series by genre.
func (s *SonarrClient) SearchSeriesByGenre(ctx context.Context, genre string, similarTo string, limit int) ([]Series, error) {
	var allSeries []Series

	if similarTo != "" {
		series, err := s.LookupSeries(ctx, similarTo)
		if err != nil {
			return nil, err
		}
		allSeries = append(allSeries, series...)
	} else {
		popularShows := []string{"Breaking Bad", "Game of Thrones", "Stranger Things", "The Office", "Friends"}
		for _, show := range popularShows {
			series, err := s.LookupSeries(ctx, show)
			if err != nil {
				continue
			}
			allSeries = append(allSeries, series...)

			if len(allSeries) >= 20 {
				break
			}
		}
	}

	var matchingSeries []Series
	genreLower := strings.ToLower(genre)

	for _, s := range allSeries {
		for _, g := range s.Genres {
			if strings.ToLower(g) == genreLower {
				matchingSeries = append(matchingSeries, s)
				break
			}
		}
		if len(matchingSeries) >= limit {
			break
		}
	}

	return matchingSeries, nil
}

func (s *SonarrClient) RequestSeriesDelete(ctx context.Context, series Series) error {
	data := map[string]any{
		"title":  series.Title,
		"tvdbId": series.ID,
	}
	_, err := s.client.Delete(ctx, fmt.Sprintf("series/%d", series.ID), data)
	if err != nil {
		return fmt.Errorf("failed to request series delete: %w", err)
	}
	return nil
}
