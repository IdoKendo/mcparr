package tools

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// MediaTools holds all the MCP tools for media management.
type MediaTools struct {
	config       Config
	sonarrClient SonarrClient
	radarrClient RadarrClient
	logger       *log.Logger
}

// Config is a simplified interface for the configuration.
type Config interface {
	SonarrURL() string
	SonarrAPIKey() string
	RadarrURL() string
	RadarrAPIKey() string
	ShowsRootPath() string
	MoviesRootPath() string
	DefaultQualityProfileID() int
}

// SonarrClient is a simplified interface for the Sonarr client.
type SonarrClient interface {
	LookupSeries(ctx context.Context, name string) ([]Series, error)
	RequestSeriesDownload(ctx context.Context, series Series, qualityProfileID int, rootFolderPath string) error
	SearchSeriesByGenre(ctx context.Context, genre string, similarTo string, limit int) ([]Series, error)
}

// RadarrClient is a simplified interface for the Radarr client.
type RadarrClient interface {
	LookupMovie(ctx context.Context, name string) ([]Movie, error)
	RequestMovieDownload(ctx context.Context, movie Movie, qualityProfileID int, rootFolderPath string) error
	SearchMoviesByGenre(ctx context.Context, genre string, similarTo string, limit int) ([]Movie, error)
}

// Series represents a TV series.
type Series struct {
	ID       int      `json:"tvdbId"`
	Title    string   `json:"title"`
	Overview string   `json:"overview,omitempty"`
	Genres   []string `json:"genres,omitempty"`
}

// Movie represents a movie.
type Movie struct {
	ID       int      `json:"tmdbId"`
	Title    string   `json:"title"`
	Overview string   `json:"overview,omitempty"`
	Genres   []string `json:"genres,omitempty"`
}

// New creates a new MediaTools instance.
func New(cfg Config, sonarrClient SonarrClient, radarrClient RadarrClient) *MediaTools {
	return &MediaTools{
		config:       cfg,
		sonarrClient: sonarrClient,
		radarrClient: radarrClient,
		logger:       log.Default(),
	}
}

// Tools returns all the MCP tools.
func (m *MediaTools) Tools() []server.ServerTool {
	return []server.ServerTool{
		m.SearchMediaID(),
		m.SearchByGenre(),
		m.RequestDownload(),
	}
}

// SearchMediaID returns a tool for searching media by ID.
func (m *MediaTools) SearchMediaID() server.ServerTool {
	tool := mcp.NewTool(
		"search_media_id",
		mcp.WithDescription("Search for media ID by name"),
		mcp.WithString(
			"type",
			mcp.Required(),
			mcp.Description("The type of media to download"),
			mcp.Enum("movie", "series"),
		),
		mcp.WithString(
			"name",
			mcp.Required(),
			mcp.Description("The name of media to find"),
		),
	)

	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		mediaType, err := request.RequireString("type")
		if err != nil {
			m.logger.Printf("Error getting media type: %v", err)
			return mcp.NewToolResultError(fmt.Sprintf("Invalid media type: %v", err)), nil
		}

		mediaName, err := request.RequireString("name")
		if err != nil {
			m.logger.Printf("Error getting media name: %v", err)
			return mcp.NewToolResultError(fmt.Sprintf("Invalid media name: %v", err)), nil
		}

		m.logger.Printf("Searching for %s with name: %s", mediaType, mediaName)

		var result string
		switch mediaType {
		case "series":
			series, err := m.sonarrClient.LookupSeries(ctx, mediaName)
			if err != nil {
				m.logger.Printf("Error looking up series: %v", err)
				return mcp.NewToolResultError(fmt.Sprintf("Failed to fetch data from Sonarr: %v", err)), nil
			}

			if len(series) > 0 {
				m.logger.Printf("Found series: %s with ID: %d", series[0].Title, series[0].ID)
				result = fmt.Sprintf("Found Sonarr series with ID: %d", series[0].ID)
			} else {
				m.logger.Printf("No series found for: %s", mediaName)
				result = "No matching series found in Sonarr."
			}
		case "movie":
			movies, err := m.radarrClient.LookupMovie(ctx, mediaName)
			if err != nil {
				m.logger.Printf("Error looking up movie: %v", err)
				return mcp.NewToolResultError(fmt.Sprintf("Failed to fetch data from Radarr: %v", err)), nil
			}

			if len(movies) > 0 {
				m.logger.Printf("Found movie: %s with ID: %d", movies[0].Title, movies[0].ID)
				result = fmt.Sprintf("Found Radarr movie with ID: %d", movies[0].ID)
			} else {
				m.logger.Printf("No movie found for: %s", mediaName)
				result = "No matching movie found in Radarr."
			}
		default:
			m.logger.Printf("Unsupported media type: %s", mediaType)
			result = fmt.Sprintf("Unsupported media type: %s. Must be 'movie' or 'series'.", mediaType)
		}

		return mcp.NewToolResultText(result), nil
	}

	return server.ServerTool{
		Tool:    tool,
		Handler: handler,
	}
}

// SearchByGenre returns a tool for searching media by genre.
func (m *MediaTools) SearchByGenre() server.ServerTool {
	tool := mcp.NewTool(
		"search_by_genre",
		mcp.WithDescription("Search for media by genre or similar content"),
		mcp.WithString(
			"type",
			mcp.Required(),
			mcp.Description("The type of media to search for"),
			mcp.Enum("movie", "series"),
		),
		mcp.WithString(
			"genre",
			mcp.Required(),
			mcp.Description("The genre to search for (e.g. action, comedy, drama)"),
		),
		mcp.WithString(
			"similar_to",
			mcp.Description("Find content similar to this title (optional)"),
		),
		mcp.WithNumber(
			"limit",
			mcp.Description("Maximum number of results to return (default: 5)"),
		),
	)

	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		mediaType, err := request.RequireString("type")
		if err != nil {
			m.logger.Printf("Error getting media type: %v", err)
			return mcp.NewToolResultError(fmt.Sprintf("Invalid media type: %v", err)), nil
		}

		genre, err := request.RequireString("genre")
		if err != nil {
			m.logger.Printf("Error getting genre: %v", err)
			return mcp.NewToolResultError(fmt.Sprintf("Invalid genre: %v", err)), nil
		}

		similarTo := request.GetString("similar_to", "")
		limit := request.GetInt("limit", 5)

		m.logger.Printf("Searching for %s with genre: %s, similar to: %s, limit: %d",
			mediaType, genre, similarTo, limit)

		var result string
		switch mediaType {
		case "series":
			series, err := m.sonarrClient.SearchSeriesByGenre(ctx, genre, similarTo, limit)
			if err != nil {
				m.logger.Printf("Error searching series by genre: %v", err)
				return mcp.NewToolResultError(fmt.Sprintf("Failed to search series by genre: %v", err)), nil
			}

			if len(series) > 0 {
				m.logger.Printf("Found %d series matching genre: %s", len(series), genre)
				var resultBuilder strings.Builder
				resultBuilder.WriteString(fmt.Sprintf("Found %d series matching genre '%s':\n", len(series), genre))

				for i, s := range series {
					resultBuilder.WriteString(fmt.Sprintf("%d. %s (ID: %d)\n", i+1, s.Title, s.ID))
				}

				result = resultBuilder.String()
			} else {
				m.logger.Printf("No series found for genre: %s", genre)
				result = fmt.Sprintf("No series found matching genre '%s'. Try a different genre such as 'drama', 'comedy', 'action', or 'thriller'.", genre)
			}

		case "movie":
			movies, err := m.radarrClient.SearchMoviesByGenre(ctx, genre, similarTo, limit)
			if err != nil {
				m.logger.Printf("Error searching movies by genre: %v", err)
				return mcp.NewToolResultError(fmt.Sprintf("Failed to search movies by genre: %v", err)), nil
			}

			if len(movies) > 0 {
				m.logger.Printf("Found %d movies matching genre: %s", len(movies), genre)
				var resultBuilder strings.Builder
				resultBuilder.WriteString(fmt.Sprintf("Found %d movies matching genre '%s':\n", len(movies), genre))

				for i, m := range movies {
					resultBuilder.WriteString(fmt.Sprintf("%d. %s (ID: %d)\n", i+1, m.Title, m.ID))
				}

				result = resultBuilder.String()
			} else {
				m.logger.Printf("No movies found for genre: %s", genre)
				result = fmt.Sprintf("No movies found matching genre '%s'. Try a different genre such as 'drama', 'comedy', 'action', or 'thriller'.", genre)
			}

		default:
			m.logger.Printf("Unsupported media type: %s", mediaType)
			result = fmt.Sprintf("Unsupported media type: %s. Must be 'movie' or 'series'.", mediaType)
		}

		return mcp.NewToolResultText(result), nil
	}

	return server.ServerTool{
		Tool:    tool,
		Handler: handler,
	}
}

// RequestDownload returns a tool for requesting media downloads.
func (m *MediaTools) RequestDownload() server.ServerTool {
	tool := mcp.NewTool(
		"request_download",
		mcp.WithDescription("Request a download for a movie or TV show"),
		mcp.WithString(
			"type",
			mcp.Required(),
			mcp.Description("The type of media to download"),
			mcp.Enum("movie", "series"),
		),
		mcp.WithString(
			"name",
			mcp.Required(),
			mcp.Description("The name of media to download"),
		),
		mcp.WithNumber(
			"id",
			mcp.Required(),
			mcp.Description("The ID of media to download"),
		),
	)

	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		mediaType, err := request.RequireString("type")
		if err != nil {
			m.logger.Printf("Error getting media type: %v", err)
			return mcp.NewToolResultError(fmt.Sprintf("Invalid media type: %v", err)), nil
		}

		mediaName, err := request.RequireString("name")
		if err != nil {
			m.logger.Printf("Error getting media name: %v", err)
			return mcp.NewToolResultError(fmt.Sprintf("Invalid media name: %v", err)), nil
		}

		mediaID, err := request.RequireInt("id")
		if err != nil {
			m.logger.Printf("Error getting media ID: %v", err)
			return mcp.NewToolResultError(fmt.Sprintf("Invalid media ID: %v", err)), nil
		}

		m.logger.Printf("Requesting download for %s: %s (ID: %d)", mediaType, mediaName, mediaID)

		var result string
		switch mediaType {
		case "series":
			series := Series{
				ID:    mediaID,
				Title: mediaName,
			}

			qualityProfileID := m.config.DefaultQualityProfileID()
			rootFolderPath := m.config.ShowsRootPath()

			m.logger.Printf("Using quality profile ID: %d and root folder: %s",
				qualityProfileID, rootFolderPath)

			err := m.sonarrClient.RequestSeriesDownload(
				ctx,
				series,
				qualityProfileID,
				rootFolderPath,
			)

			if err != nil {
				m.logger.Printf("Error requesting series download: %v", err)
				return mcp.NewToolResultError(fmt.Sprintf("Failed to request download from Sonarr: %v", err)), nil
			}

			m.logger.Printf("Successfully requested download for series: %s", mediaName)
			result = fmt.Sprintf("Download requested for Sonarr series with ID: %d", mediaID)
		case "movie":
			movie := Movie{
				ID:    mediaID,
				Title: mediaName,
			}

			qualityProfileID := m.config.DefaultQualityProfileID()
			rootFolderPath := m.config.MoviesRootPath()

			m.logger.Printf("Using quality profile ID: %d and root folder: %s",
				qualityProfileID, rootFolderPath)

			err := m.radarrClient.RequestMovieDownload(
				ctx,
				movie,
				qualityProfileID,
				rootFolderPath,
			)

			if err != nil {
				m.logger.Printf("Error requesting movie download: %v", err)
				return mcp.NewToolResultError(fmt.Sprintf("Failed to request download from Radarr: %v", err)), nil
			}

			m.logger.Printf("Successfully requested download for movie: %s", mediaName)
			result = fmt.Sprintf("Download requested for Radarr movie with ID: %d", mediaID)
		default:
			m.logger.Printf("Unsupported media type: %s", mediaType)
			result = fmt.Sprintf("Unsupported media type: %s. Must be 'movie' or 'series'.", mediaType)
		}

		return mcp.NewToolResultText(result), nil
	}

	return server.ServerTool{
		Tool:    tool,
		Handler: handler,
	}
}
