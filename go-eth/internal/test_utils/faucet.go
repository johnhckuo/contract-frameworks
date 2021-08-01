package test_utils

import (
	"context"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/johnhckuo/contract-frameworks/go-eth/internal/account"
)

const chainId = 1337

var genesisAccount *account.Account

type Node struct {
	chainId *big.Int
	client  *backends.SimulatedBackend
}

func NewNode() *Node {
	return &Node{
		chainId: big.NewInt(chainId),
	}
}

func (node *Node) CreateAndFundAccount() (newAcc *account.Account) {

	defer node.Client().Rollback()
	// create new account
	newAcc = account.NewAccount()

	// construct new tx payload

	auth, _ := bind.NewKeyedTransactorWithChainID(genesisAccount.GetPrivateKey(), node.chainId)
	fromAddress := auth.From
	nonce, err := node.client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}

	value := big.NewInt(1000000000000000000) // in wei (1 eth)
	gasLimit := uint64(21000)                // in units
	gasPrice := big.NewInt(875000000)

	var data []byte
	tx := types.NewTransaction(nonce, newAcc.GetPublicKey(), value, gasLimit, gasPrice, data)

	// sign tx using genesis account private key
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(node.chainId), genesisAccount.GetPrivateKey())
	if err != nil {
		log.Fatal(err)
	}

	// send tx
	err = node.client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("tx sent: %s\n", signedTx.Hash().Hex()) // tx sent: 0xec3ceb05642c61d33fa6c951b54080d1953ac8227be81e7b5e4e2cfed69eeb51

	// mine this block
	node.client.Commit()

	// retrieve tx receipt
	receipt, err := node.client.TransactionReceipt(context.Background(), signedTx.Hash())
	if err != nil {
		log.Fatal(err)
	}
	if receipt == nil {
		log.Fatal("receipt is nil. Forgot to commit?")
	}

	log.Printf("status: %v\n", receipt.Status) // status: 1

	return
}

func (node *Node) Connect() {

	genesisAccount = account.NewAccount()

	auth, _ := bind.NewKeyedTransactorWithChainID(genesisAccount.GetPrivateKey(), node.chainId)
	//create a genesis account and assign it an initial balance
	balance := new(big.Int)
	balance.SetString("100000000000000000000", 10) // 100 eth in wei

	address := auth.From
	genesisAlloc := map[common.Address]core.GenesisAccount{
		address: {
			Balance: balance,
		},
	}

	blockGasLimit := uint64(4712388000)
	node.client = backends.NewSimulatedBackend(genesisAlloc, blockGasLimit)

}

func (node *Node) Client() *backends.SimulatedBackend {
	return node.client
}

func (node *Node) ChainId() *big.Int {
	return node.chainId
}

/*
func (node *Account) Fund(){

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

}
*/
