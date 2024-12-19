package types

import "math/big"

// Transaction represents an Ethereum blockchain transaction
type Transaction struct {
	Hash           string
	From           string
	To             string
	Value          *big.Int
	BlockNumber    int64
	Timestamp      int64
	TransactionFee *big.Int
}
