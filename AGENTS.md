# MCParr Agent Guidelines

## Build Commands
- Build: `go build`
- Install: `go install`
- Run: `go run .`

## Test Commands
- Run all tests: `go test ./...`
- Run specific test: `go test -run TestName`
- Run tests with verbose output: `go test -v ./...`

## Lint/Format Commands
- Format code: `go fmt ./...`
- Lint code: `go vet ./...`

## Code Style Guidelines
- **Imports**: Group standard library imports first, then third-party packages, then local packages
- **Formatting**: Follow standard Go formatting (enforced by `go fmt`)
- **Types**: Use explicit type declarations for function parameters and returns
- **Naming**:
  - Use camelCase for variables and PascalCase for exported functions/types
  - Use descriptive names that reflect purpose
- **Error Handling**: Always check errors and return them to the caller
- **Comments**: Document exported functions, types, and constants
- **Config**: Environment variables are used for configuration (SONARR_URL, RADARR_URL, etc.)

## Project Structure
- Single package application with clear separation of concerns
- MCP tools defined in tools.go
- Configuration handled in config.go
