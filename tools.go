package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func SearchMediaID(cfg Config) server.ServerTool {
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
			return mcp.NewToolResultError(err.Error()), nil
		}

		mediaName, err := request.RequireString("name")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		term := strings.ReplaceAll(mediaName, " ", "+")
		var result string
		switch mediaType {
		case "series":
			url := fmt.Sprintf("%s/api/v3/series/lookup?apikey=%s&term=%s", cfg.SonarrUrl, cfg.SonarrApiKey, term)
			resp, err := http.Get(url)
			if err != nil {
				return mcp.NewToolResultError("Failed to fetch data from Sonarr: " + err.Error()), nil
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				return mcp.NewToolResultError("Non-OK HTTP status: " + resp.Status), nil
			}

			var series []struct {
				ID    int    `json:"tvdbId"`
				Title string `json:"title"`
			}
			if err := json.NewDecoder(resp.Body).Decode(&series); err != nil {
				return mcp.NewToolResultError("Failed to parse Sonarr response: " + err.Error()), nil
			}

			if len(series) > 0 {
				result = fmt.Sprintf("Found Sonarr series with ID: %d", series[0].ID)
			} else {
				result = "No matching series found in Sonarr."
			}
		case "movie":
			url := fmt.Sprintf("%s/api/v3/movie/lookup?apikey=%s&term=%s", cfg.RadarrUrl, cfg.RadarrApiKey, term)
			resp, err := http.Get(url)
			if err != nil {
				return mcp.NewToolResultError("Failed to fetch data from Radarr: " + err.Error()), nil
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				return mcp.NewToolResultError("Non-OK HTTP status: " + resp.Status), nil
			}

			var movies []struct {
				ID    int    `json:"tmdbId"`
				Title string `json:"title"`
			}
			if err := json.NewDecoder(resp.Body).Decode(&movies); err != nil {
				return mcp.NewToolResultError("Failed to parse Radarr response: " + err.Error()), nil
			}

			if len(movies) > 0 {
				result = fmt.Sprintf("Found Radarr movie with ID: %d", movies[0].ID)
			} else {
				result = "No matching movie found in Radarr."
			}
		default:
			result = "Unsupported media type for search."
		}

		return mcp.NewToolResultText(result), nil
	}

	return server.ServerTool{
		Tool:    tool,
		Handler: handler,
	}
}

func RequestDownload(cfg Config) server.ServerTool {
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
			return mcp.NewToolResultError(err.Error()), nil
		}

		mediaName, err := request.RequireString("name")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		mediaID, err := request.RequireInt("id")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		var result string
		switch mediaType {
		case "series":
			url := fmt.Sprintf("%s/api/v3/series?apikey=%s", cfg.SonarrUrl, cfg.SonarrApiKey)
			data := map[string]any{
				"title":            mediaName,
				"tvdbId":           mediaID,
				"qualityProfileId": 6,                      // TODO: not hard coded
				"rootFolderPath":   "/media/library/shows", // TODO: not hard coded
			}
			body, err := json.Marshal(data)
			if err != nil {
				return mcp.NewToolResultError("Failed to marhsal series ID for Radarr: " + err.Error()), nil
			}
			resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
			if err != nil {
				return mcp.NewToolResultError("Failed to request download from Sonarr: " + err.Error()), nil
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				return mcp.NewToolResultError("Non-OK HTTP status: " + resp.Status), nil
			}

			result = fmt.Sprintf("Download requested for Sonarr series with ID: %d", mediaID)
		case "movie":
			url := fmt.Sprintf("%s/api/v3/movie?apikey=%s", cfg.RadarrUrl, cfg.RadarrApiKey)
			data := map[string]any{
				"title":            mediaName,
				"tmdbId":           mediaID,
				"qualityProfileId": 6,                       // TODO: not hard coded
				"rootFolderPath":   "/media/library/movies", // TODO: not hard coded
			}
			body, err := json.Marshal(data)
			if err != nil {
				return mcp.NewToolResultError("Failed to marhsal movie ID for Radarr: " + err.Error()), nil
			}
			resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
			if err != nil {
				return mcp.NewToolResultError("Failed to request download from Radarr: " + err.Error()), nil
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				return mcp.NewToolResultError("Non-OK HTTP status: " + resp.Status), nil
			}

			result = fmt.Sprintf("Download requested for Radarr movie with ID: %d", mediaID)
		default:
			result = "Unsupported media type for download request."
		}

		return mcp.NewToolResultText(result), nil
	}

	return server.ServerTool{
		Tool:    tool,
		Handler: handler,
	}
}
