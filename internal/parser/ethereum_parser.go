package parser

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ethereum_parser/internal/config"
	"github.com/ethereum_parser/internal/ethereum"
	"github.com/ethereum_parser/internal/storage"
	"github.com/ethereum_parser/internal/types"
)

// EthereumParser implements the Parser interface
type EthereumParser struct {
	client      *ethereum.Client
	storage     storage.Storage
	subscribers map[string]bool
	config      *config.Config
}

func NewEthereumParser(storage storage.Storage, cfg *config.Config) (*EthereumParser, error) {
	client, err := ethereum.NewClient(cfg.EthereumRPCURL)
	if err != nil {
		return nil, err
	}

	return &EthereumParser{
		client:      client,
		storage:     storage,
		subscribers: make(map[string]bool),
		config:      cfg,
	}, nil
}

func (p *EthereumParser) GetCurrentBlock() (int64, error) {
	return p.client.GetBlockNumber()
}

func (p *EthereumParser) Subscribe(address string) bool {
	if _, exists := p.subscribers[address]; exists {
		return false
	}
	p.subscribers[address] = true
	return true
}

func (p *EthereumParser) GetTransactions(address string) ([]types.Transaction, error) {
	return p.storage.GetTransactions(address)
}

func (p *EthereumParser) startBlockPolling() {
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	var lastProcessedBlock int64

	for range ticker.C {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		currentBlock, err := p.GetCurrentBlock()
		if err != nil {
			log.Printf("Failed to get current block: %v", err)
			cancel()
			continue
		}

		// Ensure sequential block processing
		if lastProcessedBlock == 0 {
			lastProcessedBlock = currentBlock
		}

		for blockNumber := lastProcessedBlock + 1; blockNumber <= currentBlock; blockNumber++ {
			p.processBlock(ctx, blockNumber)
		}

		lastProcessedBlock = currentBlock
		cancel()
	}
}

func (p *EthereumParser) processBlock(ctx context.Context, blockNumber int64) {
	for address := range p.subscribers {
		var txs []types.Transaction
		var err error

		// Retry logic
		for retries := 0; retries < 3; retries++ {
			txs, err = p.client.GetTransactionsForAddress(ctx, address, blockNumber)
			if err == nil {
				break
			}
			log.Printf("Retry %d: Failed to get transactions for %s: %v", retries+1, address, err)
		}

		if err != nil {
			log.Printf("Failed to get transactions for %s after retries: %v", address, err)
			continue
		}

		for _, tx := range txs {
			p.storage.StoreTransaction(address, tx)

			p.notifyTransaction(tx, address, p.config.WebhookURL)
		}
	}
}

func (p *EthereumParser) notifyTransaction(tx types.Transaction, address string, webhookURL string) {
	payload, err := json.Marshal(map[string]interface{}{
		"address":      address,
		"transaction":  tx,
		"notification": "New transaction detected",
	})
	if err != nil {
		log.Printf("Failed to marshal notification payload: %v", err)
		return
	}

	if err := sendToWebhook(webhookURL, payload); err != nil {
		log.Printf("Failed to send notification for transaction %s: %v", tx.Hash, err)
		return
	}
	log.Printf("Notification sent for transaction: %s", tx.Hash)
}

func sendToWebhook(url string, payload []byte) error {
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("webhook returned non-200 status: %d", resp.StatusCode)
	}
	return nil
}
