# MCParr

⚠️WARNING: This server is still in active development and is not reliable.
This message will be removed as soon as I get this thing to work.

## Overview

MCParr is an MCP (Multi-Agent Conversational Program) server that integrates
with Sonarr and Radarr to help you manage your media library through natural
language conversations. It allows you to search for movies and TV shows,
browse by genre, and request downloads using simple commands.

## Features

- Search for movies and TV shows by name
- Browse media by genre
- Request downloads for specific media
- Integration with Sonarr (for TV shows) and Radarr (for movies)

## Prerequisites

In order to get this up and running I am using `ollama` and `mcphost`.

1. Install them and run `ollama pull qwen2.5`.
2. Clone this repository.
3. Run `ollama serve`.
4. Add `mcparr` to `$HOME/.mcp.json`:
```json
{
  "mcpServers": {
    "mcparr": {
      "command": "mcparr",
      "args": []
    }
  }
}
```
5. Run `go install .`
6. Apply the Radarr and Sonarr env variables: `SONARR_URL`, `RADARR_URL`, `SONARR_API_KEY`, `RADARR_API_KEY`.
7. Run `mcphost -m ollama:qwen2.5`
8. Chat with the AI 😃

## Configuration

MCParr requires the following environment variables to be set:

- `SONARR_API_KEY`: Your Sonarr API key
- `RADARR_API_KEY`: Your Radarr API key

Optional environment variables:

- `SONARR_URL`: The URL of your Sonarr instance (default: "http://localhost:8989")
- `RADARR_URL`: The URL of your Radarr instance (default: "http://localhost:7878")
- `SHOWS_ROOT_PATH`: The root path for TV shows (default: "/media/library/shows")
- `MOVIES_ROOT_PATH`: The root path for movies (default: "/media/library/movies")
- `DEFAULT_QUALITY_PROFILE_ID`: The default quality profile ID to use (default: 6)

## Project Structure

- `main.go`: Entry point of the application
- `internal/config`: Configuration management
- `internal/tools`: MCP tools implementation
- `pkg/client`: API clients for Sonarr and Radarr

## License

See the LICENSE file for details.

⚠️WARNING: This server is still in active development and is not reliable.
This message will be removed as soon as I get this thing to work.
