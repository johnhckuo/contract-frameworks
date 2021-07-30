package erc721

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	_ "github.com/joho/godotenv/autoload"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestErc20(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Erc721 Suite")
}

var (
	instance *Erc721
)

var _ = BeforeSuite(func() {
	var err error

	client, err := ethclient.Dial(os.Getenv("INFURA_URL"))
	if err != nil {
		log.Fatal(err)
	}

	tokenAddress := common.HexToAddress("0xa74476443119A942dE498590Fe1f2454d7D4aC0d")
	instance, err = NewErc721(tokenAddress, client)
	if err != nil {
		log.Fatal(err)
	}
})

var _ = Describe("ERC721 contract test", func() {
	Context("Check an ERC721 contract on eth network", func() {
		It("should succeeed", func() {
			address := common.HexToAddress("0xdFea73FeEaA57dbde9091C07E4d7EC2BECb51125")
			bal, err := instance.BalanceOf(&bind.CallOpts{}, address)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Printf("wei: %s\n", bal) // "wei: 74605500647408739782407023"
		})
	})
})
