package api

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"testing"

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
	testAccount *account.Account
	instance    *Api
	auth        *bind.TransactOpts
	node        *test_utils.Node
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
	var address common.Address
	var tx *types.Transaction
	address, tx, instance, err = DeployApi(auth, node.Client())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(address.Hex())   // 0x147B8eb97fD247D06C4006D269c90C1908Fb5D54
	fmt.Println(tx.Hash().Hex()) // 0xdae8ba5444eefdc99f4d45cd0c4f24056cba6a02cefbf78066ef9f4188ff7dc0
})

var _ = Describe("API Key Authorizer Test", func() {

	AfterEach(func() {
		// mockServer.AssertExpectations(GinkgoT())
		// mockStorage.AssertExpectations(GinkgoT())
	})

	BeforeEach(func() {
		// update nonce
		nonce, err := node.Client().PendingNonceAt(context.Background(), testAccount.GetPublicKey())
		if err != nil {
			log.Fatal(err)
		}
		auth.Nonce = big.NewInt(int64(nonce))
	})

	Context("Store()", func() {
		It("should success", func() {
			log.Println(instance)
			instance.Store(auth, big.NewInt(100))
		})
	})

})
