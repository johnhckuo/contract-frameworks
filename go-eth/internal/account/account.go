package account

import (
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/common"
)

type Account struct {
	privateKey *ecdsa.PrivateKey
	publicKey  common.Address
}

func NewAccount(publicKey common.Address, privateKey *ecdsa.PrivateKey) *Account {
	return &Account{
		privateKey: privateKey,
		publicKey:  publicKey,
	}
}

func (acc *Account) GetPublicKey() common.Address {
	return acc.publicKey
}

func (acc *Account) GetPrivateKey() *ecdsa.PrivateKey {
	return acc.privateKey
}
