package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func handleSend(privateKeyStr, toAddress, valueStr, rpcURL string) {
	if privateKeyStr == "" || toAddress == "" || valueStr == "" {
		log.Fatalf("Private key, to-address, and value are required for the 'send' command")
	}

	privateKey, err := crypto.HexToECDSA(privateKeyStr)
	if err != nil {
		log.Fatalf("Invalid private key: %v", err)
	}

	value := new(big.Int)
	if _, ok := value.SetString(valueStr, 10); !ok {
		log.Fatalf("Invalid value: %s", valueStr)
	}

	// Send transaction
	err = sendTransaction(privateKey, toAddress, value, rpcURL)
	if err != nil {
		log.Fatalf("Failed to send transaction: %v", err)
	}

	log.Printf("Transaction sent to %s with value %s wei", toAddress, valueStr)
}

func sendTransaction(privateKey *ecdsa.PrivateKey, toAddress string, value *big.Int, rpcURL string) error {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return fmt.Errorf("failed to connect to Ethereum RPC: %w", err)
	}
	defer client.Close()

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return fmt.Errorf("invalid public key type")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return fmt.Errorf("failed to get nonce: %w", err)
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get gas price: %w", err)
	}

	to := common.HexToAddress(toAddress)
	tx := types.NewTransaction(nonce, to, value, uint64(21000), gasPrice, nil)

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get chain ID: %w", err)
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		return fmt.Errorf("failed to sign transaction: %w", err)
	}

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return fmt.Errorf("failed to send transaction: %w", err)
	}

	log.Printf("Transaction sent: %s", signedTx.Hash().Hex())
	return nil
}
