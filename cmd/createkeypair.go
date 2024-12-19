package main

import (
	"crypto/ecdsa"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/crypto"
)

func handleCreateKey() {
	address, privateKey, err := generateAddress()
	if err != nil {
		log.Fatalf("Failed to generate address: %v", err)
	}

	fmt.Printf("Generated Address: %s\n", address)
	fmt.Printf("Private Key: %x\n", crypto.FromECDSA(privateKey))
}

func generateAddress() (string, *ecdsa.PrivateKey, error) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return "", nil, err
	}

	address := crypto.PubkeyToAddress(privateKey.PublicKey)
	return address.Hex(), privateKey, nil
}
