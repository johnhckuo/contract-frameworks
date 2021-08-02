package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/johnhckuo/contract-frameworks/go-eth/api/mynft"
)

const key = `<PRIVATEKEY>`

func main() {

	// Create an IPC based RPC connection to a remote node and an authorized transactor
	conn, err := ethclient.Dial("<INFURA URL>")
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}
	// auth, err := bind.NewTransactor(strings.NewReader(key), "my awesome super secret password")
	// if err != nil {
	// 	log.Fatalf("Failed to create authorized transactor: %v", err)
	// }

	pk, _ := crypto.HexToECDSA(key)
	pub := crypto.PubkeyToAddress(pk.PublicKey)

	// get the latest nonce
	nonce, err := conn.PendingNonceAt(context.Background(), pub)
	if err != nil {
		log.Fatal(err)
	}
	//create a transaction signer function used for authorizing transactions in the simulated blockchain.
	auth, _ := bind.NewKeyedTransactorWithChainID(pk, big.NewInt(3))
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)             // in wei
	auth.GasLimit = uint64(6000000)        // in units
	auth.GasPrice = big.NewInt(7605625975) //block base fee

	// Deploy a new awesome contract for the binding demo
	address, tx, token, err := mynft.DeployMynft(auth, conn, "john's nft", "john")
	if err != nil {
		log.Fatalf("Failed to deploy new token contract: %v", err)
	}
	fmt.Printf("Contract pending deploy: 0x%x\n", address)
	fmt.Printf("Transaction waiting to be mined: 0x%x\n\n", tx.Hash())

	// Don't even wait, check its presence in the local pending state
	time.Sleep(250 * time.Millisecond) // Allow it to be processed by the local node :P

	receipt, err := bind.WaitMined(context.Background(), conn, tx)

	fmt.Println(receipt.Status)

	name, err := token.Name(&bind.CallOpts{Pending: true})
	if err != nil {
		log.Fatalf("Failed to retrieve pending name: %v", err)
	}
	fmt.Println("Pending name:", name)
}
