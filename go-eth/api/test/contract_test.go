package test

import (
	"context"
	"log"
	"math/big"
	"strings"
	"testing"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/johnhckuo/contract-frameworks/go-eth/internal/account"
	"github.com/johnhckuo/contract-frameworks/go-eth/internal/test_utils"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestApi(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Api Suite")
}

var (
	testAccount     *account.Account
	instance        *Test
	auth            *bind.TransactOpts
	node            *test_utils.Node
	contractAddress common.Address
	logs            chan types.Log
)

var _ = BeforeSuite(func() {
	node = test_utils.NewNode()
	node.Connect()
	testAccount = node.CreateAndFundAccount()
	log.Println(testAccount)

	nonce, err := node.Client().PendingNonceAt(context.Background(), testAccount.GetPublicKey())
	if err != nil {
		log.Fatal(err)
	}

	auth, _ = bind.NewKeyedTransactorWithChainID(testAccount.GetPrivateKey(), node.ChainId())
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)     // in wei
	auth.GasLimit = uint64(300000) // in units
	auth.GasPrice = big.NewInt(875000000)

	//input := "1.0"
	//var tx *types.Transaction
	contractAddress, _, instance, err = DeployTest(auth, node.Client())
	if err != nil {
		log.Fatal(err)
	}
	node.Client().Commit()

	// fmt.Println(address.Hex())   // 0x147B8eb97fD247D06C4006D269c90C1908Fb5D54
	// fmt.Println(tx.Hash().Hex()) // 0xdae8ba5444eefdc99f4d45cd0c4f24056cba6a02cefbf78066ef9f4188ff7dc0

	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(1),
		Addresses: []common.Address{contractAddress},
	}

	logs = make(chan types.Log)

	sub, err := node.Client().SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		defer sub.Unsubscribe()
		for {
			select {
			case err := <-sub.Err():
				Fail("receive error from contract: " + err.Error())
			case vLog := <-logs:
				if vLog.BlockNumber == 0 {
					continue
				}
				log.Printf("Received event from contract: %+v", vLog)
			}
		}
	}()

})

var _ = AfterSuite(func() {
	node.Client().Close()
	//close(logs)
})

var _ = Describe("contract test", func() {

	AfterEach(func() {
		// node.Client().Commit()
	})

	BeforeEach(func() {
		// update nonce
		nonce, err := node.Client().PendingNonceAt(context.Background(), testAccount.GetPublicKey())
		if err != nil {
			log.Fatal(err)
		}
		auth.Nonce = big.NewInt(int64(nonce))

		// auth.GasPrice, err = node.Client().SuggestGasPrice(context.Background())
		// if err != nil {
		// 	log.Fatal(err)
		// }
		auth.GasPrice = big.NewInt(766599825)
	})

	Context("Store()", func() {
		It("store value should success", func() {
			tx, err := instance.Store(auth, big.NewInt(100))
			log.Println("Pending tx" + tx.Hash().Hex())
			Expect(err).Should(BeNil())
		})

		It("retrieve value should not success since block yet to be mined", func() {
			val, err := instance.Retrieve(&bind.CallOpts{Pending: false})
			Expect(val).ShouldNot(Equal(big.NewInt(100)))
			Expect(err).Should(BeNil())
		})

		It("retrieve value should success after block being mined", func() {
			node.Client().Commit()
			val, err := instance.Retrieve(&bind.CallOpts{Pending: false})
			Expect(val).Should(Equal(big.NewInt(100)))
			Expect(err).Should(BeNil())
		})

		It("reading event log", func() {
			query := ethereum.FilterQuery{
				FromBlock: nil,
				ToBlock:   nil,
				Addresses: []common.Address{
					contractAddress,
				},
			}

			logs, err := node.Client().FilterLogs(context.Background(), query)
			if err != nil {
				log.Fatal(err)
			}

			contractAbi, err := abi.JSON(strings.NewReader(string(TestABI)))
			if err != nil {
				log.Fatal(err)
			}

			for _, vLog := range logs {

				if len(vLog.Data) == 0 {
					continue
				}

				data, err := contractAbi.Unpack("ValueChanged", vLog.Data)
				if err != nil {
					log.Fatal(err)
				}
				log.Println(data[0].(*big.Int))
				log.Println(vLog.BlockNumber)
				log.Println(vLog.TxHash.Hex())

			}

		})
	})

})
