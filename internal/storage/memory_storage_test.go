package storage

import (
	"math/big"
	"testing"

	"github.com/ethereum_parser/internal/types"
)

func TestMemoryStorage(t *testing.T) {
	// Create a new memory storage
	storage := NewMemoryStorage()

	// Test transaction storage
	testTx := types.Transaction{
		Hash:        "0x123",
		From:        "0xSender",
		To:          "0xReceiver",
		Value:       big.NewInt(1000),
		BlockNumber: 100,
	}

	// Store transaction
	storage.StoreTransaction("0xSender", testTx)

	// Retrieve transactions
	txs, err := storage.GetTransactions("0xSender")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(txs) != 1 {
		t.Errorf("Expected 1 transaction, got %d", len(txs))
	}

	if txs[0].Hash != "0x123" {
		t.Errorf("Unexpected transaction hash: %s", txs[0].Hash)
	}

	// Test multiple transactions
	anotherTx := types.Transaction{
		Hash:        "0x456",
		From:        "0xSender",
		To:          "0xAnotherReceiver",
		Value:       big.NewInt(2000),
		BlockNumber: 101,
	}
	storage.StoreTransaction("0xSender", anotherTx)

	txs, err = storage.GetTransactions("0xSender")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(txs) != 2 {
		t.Errorf("Expected 2 transactions, got %d", len(txs))
	}
}
