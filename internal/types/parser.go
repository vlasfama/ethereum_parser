package types

// Parser defines the interface for blockchain transaction parsing
type Parser interface {
	GetCurrentBlock() (int64, error)
	Subscribe(address string) bool
	GetTransactions(address string) ([]Transaction, error)
}
