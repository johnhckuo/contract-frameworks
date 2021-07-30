package main

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	// for demo
)

func main() {
	/*
		client, err := ethclient.Dial("http://localhost:8545")
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("we have a connection")
		_ = client // we'll use this in the upcoming sections
	*/
	// create account

	// geth private chains id
	chainID := big.NewInt(1337)

	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatal(err)
	}

	// create a transaction signer from a single private key.
	auth, _ := bind.NewKeyedTransactorWithChainID(privateKey, chainID)

	//create a genesis account and assign it an initial balance
	balance := new(big.Int)
	balance.SetString("100000000000000000000", 10) // 100 eth in wei

	address := auth.From
	genesisAlloc := map[common.Address]core.GenesisAccount{
		address: {
			Balance: balance,
		},
	}

	blockGasLimit := uint64(4712388)
	client := backends.NewSimulatedBackend(genesisAlloc, blockGasLimit)

	// send a simple transaction

	fromAddress := auth.From
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}

	value := big.NewInt(1000000000000000000) // in wei (1 eth)
	gasLimit := uint64(21000)                // in units
	gasPrice := big.NewInt(875000000)

	toAddress := common.HexToAddress("0x4592d8f8d7b001e72cb26a73e4fa1806a51ac79d")
	var data []byte
	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, data)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatal(err)
	}

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("tx sent: %s\n", signedTx.Hash().Hex()) // tx sent: 0xec3ceb05642c61d33fa6c951b54080d1953ac8227be81e7b5e4e2cfed69eeb51

	client.Commit()

	// fetch transaction receipt

	receipt, err := client.TransactionReceipt(context.Background(), signedTx.Hash())
	if err != nil {
		log.Fatal(err)
	}
	if receipt == nil {
		log.Fatal("receipt is nil. Forgot to commit?")
	}

	fmt.Printf("status: %v\n", receipt.Status) // status: 1

	// deploy
	/*
		publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
		if !ok {
			log.Fatal("error casting public key to ECDSA")
		}

		fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
		nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
		if err != nil {
			log.Fatal(err)
		}

		gasPrice, err := client.SuggestGasPrice(context.Background())
		if err != nil {
			log.Fatal(err)
		}

		auth := bind.NewKeyedTransactor(privateKey)
		auth.Nonce = big.NewInt(int64(nonce))
		auth.Value = big.NewInt(0)     // in wei
		auth.GasLimit = uint64(300000) // in units
		auth.GasPrice = gasPrice

		//input := "1.0"
		address, tx, instance, err := api.DeployApi(auth, client)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(address.Hex())   // 0x147B8eb97fD247D06C4006D269c90C1908Fb5D54
		fmt.Println(tx.Hash().Hex()) // 0xdae8ba5444eefdc99f4d45cd0c4f24056cba6a02cefbf78066ef9f4188ff7dc0

		_ = instance
	*/

	// old
	// client, err := ethclient.Dial("http://127.0.0.1:8545")
	// if err != nil {
	// 	panic(err)
	// }

	// privateKey, err := crypto.HexToECDSA("0x701b615bbdfb9de65240bc28bd21bbc0d996645a3dd57e7b12bc2bdf6f192c82")
	// if err != nil {
	// 	panic(err)
	// }

	// publicKey := privateKey.Public()
	// publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	// if !ok {
	// 	panic("invalid key")
	// }

	// fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	// nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	// if err != nil {
	// 	panic(err)
	// }

	// chainID, err := client.ChainID(context.Background())
	// if err != nil {
	// 	panic(err)
	// }

	// auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	// if err != nil {
	// 	panic(err)
	// }
	// auth.Nonce = big.NewInt(int64(nonce))
	// auth.Value = big.NewInt(0)      // in wei
	// auth.GasLimit = uint64(3000000) // in units
	// auth.GasPrice = big.NewInt(1000000)

	// auth, err = config.Auth(client)
	// if err != nil {
	// 	panic(err)
	// }

	// address, tx, instance, err := api.DeployApi(auth, client)
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println(address.Hex())

	// _, _ = instance, tx
}
