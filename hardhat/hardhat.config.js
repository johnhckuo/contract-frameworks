const { task } = require("hardhat/config");

require("@nomiclabs/hardhat-waffle");
require('@nomiclabs/hardhat-ethers')
require("solidity-coverage");
require("hardhat-gas-reporter");
require("@nomiclabs/hardhat-solhint");
require('hardhat-log-remover');
require('hardhat-docgen');
require("hardhat-tracer");
require('dotenv').config()

// require("@nomiclabs/hardhat-truffle5");

// This is a sample Hardhat task. To learn how to create your own go to
// https://hardhat.org/guides/create-task.html
task("accounts", "Prints the list of accounts", async (taskArgs, hre) => {
  const accounts = await hre.ethers.getSigners();

  for (const account of accounts) {
    console.log(account.address);
  }
});

task("balance", "Prints an account's balance")
  .addParam("account", "The account's address")
  .setAction(async (taskArgs) => {
    const account = web3.utils.toChecksumAddress(taskArgs.account);
    const balance = await web3.eth.getBalance(account);

    console.log(web3.utils.fromWei(balance, "ether"), "ETH");
  });

task("task", "just testing", function(){
  console.log("hi")
})


// You need to export an object to set up your config
// Go to https://hardhat.org/config/ to learn more

/**
 * @type import('hardhat/config').HardhatUserConfig
 */
module.exports = {
  defaultNetwork: "hardhat",
  networks: {
    hardhat: {
      // mining: {
      //   auto: false,
      //   interval: 0
      // }
      forking: {
        url: "https://eth-ropsten.alchemyapi.io/v2/" + process.env.ALCHEMY_API_KEY,
        blockNumber: 10757709,
        enabled: true
      },
      loggingEnabled: false
    },
    rinkeby: {
      url: "https://eth-mainnet.alchemyapi.io/v2/" + process.env.ALCHEMY_API_KEY
    }
  },
  solidity: {
      version: "0.8.4",
      settings: {
        optimizer: {
          enabled: true,
          runs: 200,
        },
      },
  },
  paths: {
    sources: "./contracts",
    tests: "./test",
    cache: "./cache",
    artifacts: "./artifacts"
  },
  mocha: {
    timeout: 20000
  },
  gasReporter: {
    currency: "USD",
    enabled: true,
    excludeContracts: [],
    src: "./contracts",
  },
  docgen: {
    path: './docs',
    clear: true,
    runOnCompile: true,
  }
};
