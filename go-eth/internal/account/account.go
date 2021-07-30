package account

import (
	"crypto/ecdsa"
	"log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

type Account struct {
	privateKey *ecdsa.PrivateKey
}

func NewAccount() *Account {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatal(err)
	}
	return &Account{
		privateKey: privateKey,
	}
}

func (acc *Account) GetPublicKey() common.Address {
	return crypto.PubkeyToAddress(acc.privateKey.PublicKey)
}

func (acc *Account) GetPrivateKey() *ecdsa.PrivateKey {
	return acc.privateKey
}
