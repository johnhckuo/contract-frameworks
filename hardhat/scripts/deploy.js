// We require the Hardhat Runtime Environment explicitly here. This is optional
// but useful for running the script in a standalone fashion through `node <script>`.
//
// When running the script with `npx hardhat run <script>` you'll find the Hardhat
// Runtime Environment's members available in the global scope.
const hre = require("hardhat");

async function deploy() {
  // We get the contract to deploy
  const MyNFT = await hre.ethers.getContractFactory('MyNFT');
  console.log('Deploying MyNFT...');
  const myNFT = await MyNFT.deploy("john'nft", "john");
  await myNFT.deployed();
  console.log('MyNFT deployed to:', myNFT.address);
  process.env.NFT_ADDR = myNFT.address;

}


module.exports = {
  // scripts/deploy.js
  deploy: deploy
};

deploy()
.then(() => process.exit(0))
.catch(error => {
  console.error(error);
  process.exit(1);
});