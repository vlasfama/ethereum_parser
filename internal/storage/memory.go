package storage

import (
	"sync"

	"github.com/ethereum_parser/internal/types"
)

// Storage defines the interface for transaction storage
type Storage interface {
	StoreTransaction(address string, tx types.Transaction)
	GetTransactions(address string) ([]types.Transaction, error)
}

type MemoryStorage struct {
	transactions map[string][]types.Transaction
	mu           sync.RWMutex
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		transactions: make(map[string][]types.Transaction),
	}
}

func (ms *MemoryStorage) StoreTransaction(address string, tx types.Transaction) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	ms.transactions[address] = append(ms.transactions[address], tx)
}

func (ms *MemoryStorage) GetTransactions(address string) ([]types.Transaction, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	return ms.transactions[address], nil
}
