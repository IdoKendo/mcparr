package main

import (
	"io"
	"log"
	"os"
	"os/user"
	"path/filepath"

	"github.com/mark3labs/mcp-go/server"

	"github.com/IdoKendo/mcparr/internal/config"
	"github.com/IdoKendo/mcparr/internal/tools"
	"github.com/IdoKendo/mcparr/pkg/client"
)

func initLogger() error {
	usr, err := user.Current()
	if err != nil {
		return err
	}

	logPath := filepath.Join(usr.HomeDir, ".cache", "mcparr", "history.log")
	err = os.MkdirAll(filepath.Dir(logPath), 0755)
	if err != nil {
		return err
	}

	logFile, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	multiWriter := io.MultiWriter(os.Stdout, logFile)

	log.SetOutput(multiWriter)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetPrefix("[MCParr] ")

	return nil
}

func main() {
	err := initLogger()
	if err != nil {
		log.Panic(err)
	}

	log.Println("Starting MCParr server...")

	s := server.NewMCPServer(
		"MCP Arr",
		"1.0.0",
		server.WithToolCapabilities(false),
		server.WithRecovery(),
	)
	log.Println("MCP server initialized")

	cfg := config.New()
	log.Println("Configuration loaded")

	log.Println("Initializing API clients...")
	sonarrClient := client.NewSonarrClient(cfg.SonarrURL(), cfg.SonarrAPIKey())
	radarrClient := client.NewRadarrClient(cfg.RadarrURL(), cfg.RadarrAPIKey())
	log.Println("API clients initialized")

	sonarrAdapter := tools.NewSonarrClientAdapter(sonarrClient)
	radarrAdapter := tools.NewRadarrClientAdapter(radarrClient)
	log.Println("Client adapters created")

	log.Println("Initializing MCP tools...")
	mediaTools := tools.New(cfg, sonarrAdapter, radarrAdapter)
	log.Println("MCP tools initialized")

	s.AddTools(mediaTools.Tools()...)
	log.Println("Tools added to server")

	log.Println("Starting server...")
	if err := server.ServeStdio(s); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
