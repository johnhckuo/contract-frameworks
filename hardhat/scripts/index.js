// scripts/index.js

var tools = require('./deploy');
const hre = require("hardhat");

async function main () {

    await tools.deploy()
    // Retrieve accounts from the local node
    const accounts = await hre.ethers.provider.listAccounts();
    console.log(accounts);
    const address = process.env.NFT_ADDR;
    const MyNFT = await hre.ethers.getContractFactory('MyNFT');
    const nft = await MyNFT.attach(address);

    // Mint
    await nft.mint();

    console.log('MyNFT value is', value.toString());
  }
  
  main()
    .then(() => process.exit(0))
    .catch(error => {
      console.error(error);
      process.exit(1);
    });