package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ethereum_parser/internal/config"
)

// Generated Address: 0x97c5aBe06209123987392D4489b54B8b213E0Dac
// Private Key: 2c672b55cb8b98ddb0237326d1c0ab4d9a2a82fa02947afb29b741f1b9505a75

// Generated Address: 0xC15683bC491872ff122A11eDB9a2b038f8BA15AD
// Private Key: c739efc78f5cef4812c13ca6bb351c69e2c7f107eae7eba8695911e66f588ce9

// https://ethereum-sepolia-rpc.publicnode.com

func main() {
	// Define subcommands
	startCmd := flag.NewFlagSet("start", flag.ExitOnError)
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
	createKeyCmd := flag.NewFlagSet("create_key", flag.ExitOnError)

	// Create a configuration object
	cfg := config.NewConfig()

	// Define flags for the "start" subcommand
	startCmd.StringVar(&cfg.EthereumRPCURL, "rpc-url", cfg.EthereumRPCURL, "Ethereum RPC URL")
	startCmd.IntVar(&cfg.HTTPPort, "port", cfg.HTTPPort, "HTTP server port")
	startCmd.StringVar(&cfg.LogLevel, "log-level", cfg.LogLevel, "Logging level (debug, info, warn, error)")

	// Define flags for the "send" subcommand
	privateKey := sendCmd.String("private-key", "", "Sender's private key")
	toAddress := sendCmd.String("to-address", "", "Recipient's Ethereum address")
	value := sendCmd.String("value", "", "Value to send (in wei)")
	sendCmd.StringVar(&cfg.EthereumRPCURL, "rpc-url", cfg.EthereumRPCURL, "Ethereum RPC URL")

	// Parse the top-level command
	if len(os.Args) < 2 {
		fmt.Println("Expected 'start', 'send', or 'create_key' subcommands")
		return
	}

	// Load environment variables into the config
	cfg.LoadEnvironmentVariables()

	switch os.Args[1] {
	case "start":
		startCmd.Parse(os.Args[2:])
		handleStart(cfg)

	case "send":
		sendCmd.Parse(os.Args[2:])
		handleSend(*privateKey, *toAddress, *value, cfg.EthereumRPCURL)

	case "create_key":
		createKeyCmd.Parse(os.Args[2:])
		handleCreateKey()

	default:
		fmt.Println("Unknown command. Expected 'start', 'send', or 'create_key'")
	}
}
