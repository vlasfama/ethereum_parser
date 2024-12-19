package ethereum

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"strconv"
	"time"

	"github.com/ethereum_parser/internal/types"
)

// JSONRPCRequest represents a standard JSON-RPC request
type JSONRPCRequest struct {
	JSONRPC string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	ID      int           `json:"id"`
}

// JSONRPCResponse represents a standard JSON-RPC response
type JSONRPCResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	Result  json.RawMessage `json:"result"`
	Error   *RPCError       `json:"error"`
	ID      int             `json:"id"`
}

// RPCError represents JSON-RPC error structure
type RPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Client handles Ethereum JSON-RPC interactions
type Client struct {
	rpcURL     string
	httpClient *http.Client
}

// NewClient creates a new Ethereum JSON-RPC client
func NewClient(rpcURL string) (*Client, error) {
	if rpcURL == "" {
		return nil, fmt.Errorf("RPC URL cannot be empty")
	}

	return &Client{
		rpcURL: rpcURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}, nil
}

// makeJSONRPCRequest sends a JSON-RPC request and returns the response
func (c *Client) makeJSONRPCRequest(method string, params []interface{}) (*JSONRPCResponse, error) {
	// Prepare request payload
	payload := JSONRPCRequest{
		JSONRPC: "2.0",
		Method:  method,
		Params:  params,
		ID:      1,
	}

	// Convert payload to JSON
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON-RPC request: %v", err)
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", c.rpcURL, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Send request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send HTTP request: %v", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	// Parse JSON-RPC response
	var rpcResp JSONRPCResponse
	if err := json.Unmarshal(body, &rpcResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON-RPC response: %v", err)
	}

	// Check for RPC error
	if rpcResp.Error != nil {
		return nil, fmt.Errorf("RPC error: code %d, message: %s",
			rpcResp.Error.Code, rpcResp.Error.Message)
	}

	return &rpcResp, nil
}

// GetBlockNumber retrieves the latest block number
func (c *Client) GetBlockNumber() (int64, error) {
	resp, err := c.makeJSONRPCRequest("eth_blockNumber", []interface{}{})
	if err != nil {
		return 0, err
	}

	// Remove "0x" prefix and convert hex to int64
	var hexBlockNum string
	if err := json.Unmarshal(resp.Result, &hexBlockNum); err != nil {
		return 0, fmt.Errorf("failed to parse block number: %v", err)
	}

	blockNum, err := strconv.ParseInt(hexBlockNum[2:], 16, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to convert block number: %v", err)
	}

	return blockNum, nil
}

// GetTransactionsForAddress retrieves transactions for a specific address
func (c *Client) GetTransactionsForAddress(
	ctx context.Context,
	address string,
	blockNumber int64,
) ([]types.Transaction, error) {
	// Validate address
	if !isValidEthereumAddress(address) {
		return nil, fmt.Errorf("invalid Ethereum address: %s", address)
	}

	// Convert block number to hex
	blockNumberHex := fmt.Sprintf("0x%x", blockNumber)

	// Fetch block details
	blockResp, err := c.makeJSONRPCRequest("eth_getBlockByNumber",
		[]interface{}{blockNumberHex, true})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch block: %v", err)
	}

	// Parse block transactions
	var block struct {
		Transactions []struct {
			Hash        string `json:"hash"`
			From        string `json:"from"`
			To          string `json:"to"`
			Value       string `json:"value"`
			BlockNumber string `json:"blockNumber"`
			Timestamp   string `json:"timestamp"`
		} `json:"transactions"`
	}

	if err := json.Unmarshal(blockResp.Result, &block); err != nil {
		return nil, fmt.Errorf("failed to parse block transactions: %v", err)
	}

	var transactions []types.Transaction

	for _, tx := range block.Transactions {
		// Filter transactions related to the address
		if tx.From != address && (tx.To == "" || tx.To != address) {
			continue
		}

		// Convert hex values
		value, _ := new(big.Int).SetString(tx.Value[2:], 16)
		blockNum, _ := strconv.ParseInt(tx.BlockNumber[2:], 16, 64)
		timestamp, _ := strconv.ParseInt(tx.Timestamp[2:], 16, 64)

		transactions = append(transactions, types.Transaction{
			Hash:        tx.Hash,
			From:        tx.From,
			To:          tx.To,
			Value:       value,
			BlockNumber: blockNum,
			Timestamp:   timestamp,
		})
	}

	return transactions, nil
}

// GetBalance retrieves the balance of an address
func (c *Client) GetBalance(address string) (*big.Int, error) {
	// Validate address
	if !isValidEthereumAddress(address) {
		return nil, fmt.Errorf("invalid Ethereum address: %s", address)
	}

	// Get balance
	resp, err := c.makeJSONRPCRequest("eth_getBalance",
		[]interface{}{address, "latest"})
	if err != nil {
		return nil, err
	}

	// Parse hex balance
	var hexBalance string
	if err := json.Unmarshal(resp.Result, &hexBalance); err != nil {
		return nil, fmt.Errorf("failed to parse balance: %v", err)
	}

	// Convert hex to big.Int (wei)
	balance, ok := new(big.Int).SetString(hexBalance[2:], 16)
	if !ok {
		return nil, fmt.Errorf("failed to convert balance")
	}

	return balance, nil
}

// Helper function to validate Ethereum address
func isValidEthereumAddress(address string) bool {
	// Basic validation
	if len(address) != 42 {
		return false
	}
	if address[:2] != "0x" {
		return false
	}
	return true
}
