package mynft

import (
	"context"
	"log"
	"math/big"
	"testing"
	"time"

	ethereum "github.com/ethereum/go-ethereum"
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
	instance        *Mynft
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

	//create a transaction signer function used for authorizing transactions in the simulated blockchain. 
	auth, _ = bind.NewKeyedTransactorWithChainID(testAccount.GetPrivateKey(), node.ChainId())
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)     // in wei
	auth.GasLimit = uint64(300000) // in units
	auth.GasPrice = big.NewInt(875000000)

	//input := "1.0"
	//var tx *types.Transaction
	contractAddress, _, instance, err = DeployMynft(auth, node.Client(), "john'nft", "john")
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
	close(logs)
})

var _ = Describe("contract test", func() {

	AfterEach(func() {
		// node.Client().Commit()
		time.Sleep(3 * time.Second)
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

	var (
		tokenOwner *account.Account
	)


	BeforeSuite(func(){
		tokenOwner = node.CreateAndFundAccount()
	})

	Context("Mint()", func() {
		It("mint should success", func() {
			val, err := instance.Mint(auth, tokenOwner.GetPublicKey())
			log.Println("Pending tx" + tx.Hash().Hex())
			node.Client().Commit()
			Expect(val).Should(Equal(big.NewInt(0)))

			Expect(err).Should(BeNil())
		})

		// It("Receive transfer event", func() {
		// 	// emit Transfer(address(0), to, tokenId);

		// })
		
		It("check owner", func() {
			tx, err := instance.OwnerOf(auth, testAccount.GetPublicKey())
			log.Println("Pending tx" + tx.Hash().Hex())
			node.Client().Commit()
			Expect(err).Should(BeNil())
		})

		It("burn should success", func() {
			tx, err := instance.Burn(auth, testAccount.GetPublicKey())
			log.Println("Pending tx" + tx.Hash().Hex())
			node.Client().Commit()
			Expect(err).Should(BeNil())
		})


		// It("retrieve value should not success since block yet to be mined", func() {
		// 	val, err := instance.Retrieve(&bind.CallOpts{Pending: false})
		// 	Expect(val).ShouldNot(Equal(big.NewInt(100)))
		// 	Expect(err).Should(BeNil())
		// })

		// It("retrieve value should success after block being mined", func() {
		// 	node.Client().Commit()
		// 	val, err := instance.Retrieve(&bind.CallOpts{Pending: false})
		// 	Expect(val).Should(Equal(big.NewInt(100)))
		// 	Expect(err).Should(BeNil())
		// })

	})

})
