# MCParr

An MCP server to manage media library using Sonarr and Radarr.

⚠️WARNING: This server is still in active development and is not working.
This message will be removed as soon as I get this thing to work.

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
7. Chat with the AI 😃

⚠️WARNING: This server is still in active development and is not working.
This message will be removed as soon as I get this thing to work.
