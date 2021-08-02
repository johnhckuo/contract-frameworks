package mynft

import (
	"bytes"
	"context"
	"encoding/hex"
	"log"
	"math/big"
	"testing"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/johnhckuo/contract-frameworks/go-eth/internal/account"
	"github.com/johnhckuo/contract-frameworks/go-eth/internal/test_utils"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestMyNFT(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "MyNFT Suite")
}

type LogTransfer struct {
	From   common.Address
	To     common.Address
	Tokens *big.Int
}

type LogApproval struct {
	TokenOwner common.Address
	Spender    common.Address
	Tokens     *big.Int
}

var (
	testAccount     *account.Account
	instance        *Mynft
	auth            *bind.TransactOpts
	node            *test_utils.Node
	contractAddress common.Address
	//logs              chan types.Log
	transferEventSink chan *MynftTransfer
	txStatus          chan int

	// erc721
	tokenOwner *account.Account
	tokenId    *big.Int
	newOwner   *account.Account
)

var _ = BeforeSuite(func() {

	// create a local private chain
	node = test_utils.NewNode()
	node.Connect()

	// create root account
	testAccount = node.CreateAndFundAccount()
	log.Println(testAccount)

	// get the latest nonce
	nonce, err := node.Client().PendingNonceAt(context.Background(), testAccount.GetPublicKey())
	if err != nil {
		log.Fatal(err)
	}
	//create a transaction signer function used for authorizing transactions in the simulated blockchain.
	auth, _ = bind.NewKeyedTransactorWithChainID(testAccount.GetPrivateKey(), node.ChainId())
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)            // in wei
	auth.GasLimit = uint64(300000000)     // in units
	auth.GasPrice = big.NewInt(765625975) //block base fee

	// var tx *types.Transaction
	contractAddress, _, instance, err = DeployMynft(auth, node.Client(), "john's nft", "john")
	if err != nil {
		log.Fatal(err)
	}
	node.Client().Commit()

	log.Println("contract created! :" + contractAddress.Hex())
	// log.Println(tx.Hash().Hex())

	// query := ethereum.FilterQuery{
	// 	FromBlock: nil,
	// 	Addresses: []common.Address{contractAddress},
	// }

	// logs = make(chan types.Log)

	// sub, err := node.Client().SubscribeFilterLogs(context.Background(), query, logs)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// go func() {
	// 	for {
	// 		select {
	// 		case err := <-sub.Err():
	// 			log.Fatal("receive error from contract: " + err.Error())
	// 		case vLog, received := <-logs:
	// 			if vLog.BlockNumber == 0 {
	// 				continue
	// 			}
	// 			// chan closed
	// 			if !received {
	// 				sub.Unsubscribe()
	// 				return
	// 			}
	// 			log.Printf("Received event from contract: %+v", vLog)
	// 		}
	// 	}
	// }()

	transferEventSink = make(chan *MynftTransfer)
	sub, err := instance.WatchTransfer(nil, transferEventSink, nil, nil, nil)
	Expect(err).Should(BeNil())
	go func() {
		for {
			select {
			case err := <-sub.Err():
				log.Fatal("receive error from contract: " + err.Error())
			case transferEvent, received := <-transferEventSink:

				// chan closed
				if !received {
					log.Println("chan close")
					sub.Unsubscribe()
					return
				}
				log.Println("[Transfer event]")
				log.Println("From: " + transferEvent.From.String())
				log.Println("To: " + transferEvent.To.String())
				log.Println("Value: " + transferEvent.TokenId.String())

			}
		}
	}()

	txStatus = make(chan int, 1)

	tokenOwner = node.CreateAndFundAccount()
	newOwner = node.CreateAndFundAccount()

})

// listen to new block event
// if new block is mined, we check the tx we've just sent and get it's status from tx receipt
// if status is 0 means tx failed, 1 means success, -1 means there's an error occured
func CheckTxStatus(hashToRead string) int {
	soc := make(chan *types.Header)
	newSub, err := node.Client().SubscribeNewHead(context.Background(), soc)
	if err != nil {
		return -1
	}

	go func() {
		defer newSub.Unsubscribe()
		for {
			select {
			case err := <-newSub.Err():
				log.Println("error while listening to new head", err)
				txStatus <- -1
			case header := <-soc:
				log.Println("checking the status of tx: " + header.TxHash.Hex())
				// 0 means failed, 1 means success
				txStatus <- checkTransactionReceipt(hashToRead)
				return
			}
		}
	}()
	return -1
}

func checkTransactionReceipt(_txHash string) int {
	txHash := common.HexToHash(_txHash)
	tx, err := node.Client().TransactionReceipt(context.Background(), txHash)
	if err != nil {
		return (-1)
	}
	return (int(tx.Status))
}

var _ = AfterSuite(func() {
	node.Client().Close()
	//close(logs)
	close(transferEventSink)
})

var _ = Describe("contract test", func() {

	AfterEach(func() {
		node.Client().Commit()
	})

	BeforeEach(func() {
		// update nonce
		nonce, err := node.Client().PendingNonceAt(context.Background(), testAccount.GetPublicKey())
		if err != nil {
			log.Fatal(err)
		}
		auth.Nonce = big.NewInt(int64(nonce))

		auth.GasPrice = big.NewInt(766599825)
	})

	Context("SupportInterface()", func() {
		It("Should success since it is ERC721", func() {
			byteArr, _ := hex.DecodeString("80ac58cd")
			var input [4]byte
			copy(input[:], byteArr[:4])
			bool, err := instance.SupportsInterface(&bind.CallOpts{
				From: tokenOwner.GetPublicKey(),
			}, input)
			Expect(err).Should(BeNil())
			Expect(bool).Should(BeTrue())
		})
	})

	Context("Meta data", func() {
		It("Name()", func() {
			name, err := instance.Name(&bind.CallOpts{})
			Expect(name).Should(Equal("john's nft"))
			Expect(err).Should(BeNil())
		})

		It("Symbol()", func() {
			symbol, err := instance.Symbol(&bind.CallOpts{})
			Expect(symbol).Should(Equal("john"))
			Expect(err).Should(BeNil())
		})
	})

	Context("Mint()", func() {

		It("mint should success", func() {
			tx, err := instance.Mint(auth, tokenOwner.GetPublicKey())
			node.Client().Commit()
			Expect(err).Should(BeNil())
			// should start from 1
			receipt, err := bind.WaitMined(context.Background(), node.Client(), tx)
			Expect(err).Should(BeNil())
			Expect(receipt.Status).Should(Equal(uint64(1)))
		})

		It("mint request from no-minter, should failed", func() {
			tempAccount := node.CreateAndFundAccount()

			// read for signing operations
			tempAuth, _ := bind.NewKeyedTransactorWithChainID(tempAccount.GetPrivateKey(), node.ChainId())
			tempNonce, _ := node.Client().PendingNonceAt(context.Background(), tempAccount.GetPublicKey())
			tempAuth.Nonce = big.NewInt(int64(tempNonce))
			tempAuth.Value = big.NewInt(0)            // in wei
			tempAuth.GasLimit = uint64(300000000)     // in units
			tempAuth.GasPrice = big.NewInt(765625975) //block base fee

			tx, err := instance.Mint(tempAuth, tempAccount.GetPublicKey())
			node.Client().Commit()
			Expect(err).Should(BeNil())
			receipt, err := bind.WaitMined(context.Background(), node.Client(), tx)
			Expect(err).Should(BeNil())
			Expect(receipt.Status).Should(Equal(uint64(0)))

			_, err = test_utils.GetFailingMessage(node.Client(), receipt.TxHash)
			Expect(err).ShouldNot(BeNil())
			Expect(err.Error()).Should(Equal("execution reverted: MyNFT: must have minter role to mint"))
		})

		It("Reading transfer event", func() {

			// reading event
			query := ethereum.FilterQuery{
				Addresses: []common.Address{
					contractAddress,
				},
			}

			logs, err := node.Client().FilterLogs(context.Background(), query)
			if err != nil {
				log.Fatal(err)
			}

			// contractAbi, err := abi.JSON(strings.NewReader(string(MynftABI)))
			// if err != nil {
			// 	log.Fatal(err)
			// }

			logTransferSig := []byte("Transfer(address,address,uint256)")
			LogApprovalSig := []byte("Approval(address,address,uint256)")
			logTransferSigHash := crypto.Keccak256Hash(logTransferSig)
			logApprovalSigHash := crypto.Keccak256Hash(LogApprovalSig)

			for _, vLog := range logs {
				// fmt.Printf("Log Block Number: %d\n", vLog.BlockNumber)
				// fmt.Printf("Log Index: %d\n", vLog.Index)

				if len(vLog.Topics) == 0 {
					continue
				}

				switch vLog.Topics[0].Hex() {
				case logTransferSigHash.Hex():

					var transferEvent LogTransfer

					transferEvent.From = common.HexToAddress(vLog.Topics[1].Hex())
					transferEvent.To = common.HexToAddress(vLog.Topics[2].Hex())
					transferEvent.Tokens = vLog.Topics[3].Big()
					tokenId = transferEvent.Tokens

					Expect(tokenId).Should(Equal(big.NewInt(1)))
					log.Printf("Transfer event received : %+v", transferEvent)
				case logApprovalSigHash.Hex():
					log.Println("approval")

				}
			}
		})

		It("check ownership", func() {

			// check ownership
			tx, err := instance.OwnerOf(&bind.CallOpts{
				From: tokenOwner.GetPublicKey(),
			}, tokenId)
			Expect(tx).Should(Equal(tokenOwner.GetPublicKey()))
			Expect(err).Should(BeNil())

		})

		It("mint should success2", func() {
			_, err := instance.Mint(auth, tokenOwner.GetPublicKey())
			node.Client().Commit()
			Expect(err).Should(BeNil())
		})

		It("check ownership after mint", func() {
			owner, err := instance.OwnerOf(&bind.CallOpts{}, big.NewInt(2))
			Expect(owner).Should(Equal(tokenOwner.GetPublicKey()))
			Expect(err).Should(BeNil())

		})

		It("check balance", func() {

			// check ownership
			balance, err := instance.BalanceOf(&bind.CallOpts{}, tokenOwner.GetPublicKey())
			Expect(balance).Should(Equal(big.NewInt(2)))
			Expect(err).Should(BeNil())

		})

		// It("burn should success", func() {
		// 	tx, err := instance.Burn(auth, tokenId)
		// 	log.Println("Pending tx" + tx.Hash().Hex())
		// 	node.Client().Commit()
		// 	Expect(err).Should(BeNil())
		// })

		// It("check ownership after burn", func() {
		// 	tx, err := instance.OwnerOf(&bind.CallOpts{
		// 		From:    tokenOwner.GetPublicKey(),
		// 		Pending: false,
		// 	}, tokenId)
		// 	Expect(tx).ShouldNot(Equal(tokenOwner.GetPublicKey()))
		// 	Expect(err).ShouldNot(BeNil())

	})

	Context("ERC721Enumerable", func() {

		var count *big.Int
		It("totalSupply()", func() {
			var err error
			count, err = instance.TotalSupply(&bind.CallOpts{})
			Expect(err).Should(BeNil())
			Expect(count).Should(Equal(big.NewInt(2)))
		})

		It("tokenByIndex() should success", func() {
			id, err := instance.TokenByIndex(&bind.CallOpts{}, big.NewInt(0))
			Expect(err).Should(BeNil())
			Expect(id).Should(Equal(tokenId))
		})

		It("tokenByIndex() should failed since we only mint 2 nft", func() {
			_, err := instance.TokenByIndex(&bind.CallOpts{}, big.NewInt(2))
			Expect(err).ShouldNot(BeNil())
		})

		It("tokenOfOwnerByIndex()", func() {
			start := big.NewInt(0)
			end := count

			var ids []*big.Int

			for i := new(big.Int).Set(start); i.Cmp(end) < 0; i.Add(i, big.NewInt(1)) {
				// check ownership
				id, err := instance.TokenOfOwnerByIndex(&bind.CallOpts{}, tokenOwner.GetPublicKey(), i)
				Expect(err).Should(BeNil())
				ids = append(ids, id)

			}

			Expect(len(ids)).Should(Equal(2))
			Expect(ids).Should(Equal([]*big.Int{big.NewInt(1), big.NewInt(2)}))

		})

	})

	Context("Approve()", func() {

		It("check approved address of this token", func() {
			addr, err := instance.GetApproved(&bind.CallOpts{}, tokenId)
			res := bytes.Equal(addr.Bytes(), make([]byte, 20))
			// instance.IsApprovedForAll()
			Expect(err).Should(BeNil())
			Expect(res).Should(BeTrue())
		})

		It("approved to transfer token with `tokenId` by real owner", func() {

			realAuth, _ := bind.NewKeyedTransactorWithChainID(tokenOwner.GetPrivateKey(), node.ChainId())
			realAuth.Value = big.NewInt(0)        // in wei
			realAuth.GasLimit = uint64(300000000) // in units
			nonce, _ := node.Client().PendingNonceAt(context.Background(), tokenOwner.GetPublicKey())
			realAuth.Nonce = big.NewInt(int64(nonce))
			realAuth.GasPrice = big.NewInt(766599825)

			tx, err := instance.Approve(realAuth, newOwner.GetPublicKey(), tokenId)
			CheckTxStatus(tx.Hash().Hex())
			node.Client().Commit()
			Expect(err).Should(BeNil())
			Expect(<-txStatus).Should(Equal(1))
		})

		It("check status after approve", func() {
			addr, err := instance.GetApproved(&bind.CallOpts{}, tokenId)
			res := bytes.Equal(addr.Bytes(), newOwner.GetPublicKey().Bytes())
			// instance.IsApprovedForAll()
			Expect(err).Should(BeNil())
			Expect(res).Should(BeTrue())
		})

		// It("Reading approval event", func() {

		// 	// reading event
		// 	query := ethereum.FilterQuery{
		// 		Addresses: []common.Address{
		// 			contractAddress,
		// 		},
		// 	}

		// 	logs, err := node.Client().FilterLogs(context.Background(), query)
		// 	if err != nil {
		// 		log.Fatal(err)
		// 	}

		// 	// contractAbi, err := abi.JSON(strings.NewReader(string(MynftABI)))
		// 	// if err != nil {
		// 	// 	log.Fatal(err)
		// 	// }

		// 	logTransferSig := []byte("Transfer(address,address,uint256)")
		// 	LogApprovalSig := []byte("Approval(address,address,uint256)")
		// 	logTransferSigHash := crypto.Keccak256Hash(logTransferSig)
		// 	logApprovalSigHash := crypto.Keccak256Hash(LogApprovalSig)

		// 	for _, vLog := range logs {

		// 		if len(vLog.Topics) == 0 {
		// 			continue
		// 		}

		// 		switch vLog.Topics[0].Hex() {
		// 		case logTransferSigHash.Hex():

		// 			log.Println("transfer event")
		// 		case logApprovalSigHash.Hex():

		// 			log.Println("Log Name: Approval")

		// 			var approvalEvent LogApproval

		// 			approvalEvent.TokenOwner = common.HexToAddress(vLog.Topics[1].Hex())
		// 			approvalEvent.Spender = common.HexToAddress(vLog.Topics[2].Hex())
		// 			approvalEvent.Tokens = vLog.Topics[3].Big()

		// 			Expect(approvalEvent).Should(Equal(LogApproval{
		// 				TokenOwner: tokenOwner.GetPublicKey(),
		// 				Spender:    receiver.GetPublicKey(),
		// 				Tokens:     tokenId,
		// 			}))
		// 			log.Printf("Token Owner: %s\n", approvalEvent.TokenOwner.Hex())
		// 			log.Printf("Spender: %s\n", approvalEvent.Spender.Hex())
		// 			log.Printf("Tokens: %s\n", approvalEvent.Tokens.String())

		// 		}
		// 	}
		// })
	})

	Context("Transfer", func() {
		var receiver *account.Account
		It("get allowance of this new account first", func() {
			receiver = node.CreateAndFundAccount()
			val, err := instance.BalanceOf(&bind.CallOpts{}, receiver.GetPublicKey())
			Expect(err).Should(BeNil())
			Expect(len(val.Bits())).Should(Equal(0))
		})

		It("transfer() by fake owner", func() {

			fakeOwner := node.CreateAndFundAccount()
			realAuth, _ := bind.NewKeyedTransactorWithChainID(fakeOwner.GetPrivateKey(), node.ChainId())
			realAuth.Value = big.NewInt(0)        // in wei
			realAuth.GasLimit = uint64(300000000) // in units
			nonce, _ := node.Client().PendingNonceAt(context.Background(), fakeOwner.GetPublicKey())
			realAuth.Nonce = big.NewInt(int64(nonce))
			realAuth.GasPrice = big.NewInt(766599825)

			tx, err := instance.TransferFrom(realAuth, fakeOwner.GetPublicKey(), receiver.GetPublicKey(), tokenId)
			Expect(err).Should(BeNil())

			// listen to new block event
			CheckTxStatus(tx.Hash().Hex())

			// mine the block
			node.Client().Commit()

			// check tx status after block is being mined
			Expect(<-txStatus).Should(Equal(0))

		})

		It("transfer() initiated by approved party, and transfer token from tokenOwner to receiver", func() {

			realAuth, _ := bind.NewKeyedTransactorWithChainID(newOwner.GetPrivateKey(), node.ChainId())
			realAuth.Value = big.NewInt(0)        // in wei
			realAuth.GasLimit = uint64(300000000) // in units
			nonce, _ := node.Client().PendingNonceAt(context.Background(), newOwner.GetPublicKey())
			realAuth.Nonce = big.NewInt(int64(nonce))
			realAuth.GasPrice = big.NewInt(766599825)

			tx, err := instance.TransferFrom(realAuth, tokenOwner.GetPublicKey(), receiver.GetPublicKey(), tokenId)
			Expect(err).Should(BeNil())

			// listen to new block event
			CheckTxStatus(tx.Hash().Hex())

			// mine the block
			node.Client().Commit()

			// check tx status after block is being mined
			Expect(<-txStatus).Should(Equal(1))

		})

		It("allowance should increased since it is from a real owner", func() {
			val, err := instance.BalanceOf(&bind.CallOpts{}, receiver.GetPublicKey())
			Expect(err).Should(BeNil())
			Expect(val).Should(Equal(big.NewInt(1)))
		})

	})

})
