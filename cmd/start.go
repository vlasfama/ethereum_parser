package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/ethereum_parser/internal/api"
	"github.com/ethereum_parser/internal/config"
	"github.com/ethereum_parser/internal/parser"
	"github.com/ethereum_parser/internal/storage"
	"github.com/ethereum_parser/internal/types"
)

func handleStart(cfg *config.Config) {

	flag.StringVar(&cfg.EthereumRPCURL, "rpc-url", cfg.EthereumRPCURL, "Ethereum RPC endpoint URL")
	flag.IntVar(&cfg.HTTPPort, "port", cfg.HTTPPort, "HTTP server port")
	flag.StringVar(&cfg.LogLevel, "log-level", cfg.LogLevel, "Logging level (debug, info, warn, error)")

	cfg.LoadEnvironmentVariables()

	flag.Parse()

	setupLogging(cfg.LogLevel)

	// Initialize in-memory storage
	memoryStore := storage.NewMemoryStorage()

	ethParser, err := parser.NewEthereumParser(memoryStore, cfg)
	if err != nil {
		log.Fatalf("Failed to initialize parser: %v", err)
	}

	log.Printf("Starting Ethereum Transaction Parser")
	log.Printf("RPC URL: %s", cfg.EthereumRPCURL)
	log.Printf("HTTP Port: %d", cfg.HTTPPort)

	// Start HTTP server
	if err := startHTTPServer(ethParser, cfg.HTTPPort); err != nil {
		log.Fatalf("Failed to start HTTP server: %v", err)
	}
}

func setupLogging(level string) {
	switch level {
	case "debug":
		log.SetOutput(os.Stdout)
	case "error":
		log.SetOutput(os.Stderr)
	default:

		log.SetOutput(os.Stdout)
	}
}

func startHTTPServer(parser types.Parser, port int) error {
	// Implement HTTP server startup with configurable port
	server := api.NewHTTPServer(parser)
	serverAddr := fmt.Sprintf(":%d", port)
	log.Printf("Starting HTTP server on %s", serverAddr)
	return server.Start(serverAddr)
}
