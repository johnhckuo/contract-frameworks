// We require the Hardhat Runtime Environment explicitly here. This is optional
// but useful for running the script in a standalone fashion through `node <script>`.
//
// When running the script with `npx hardhat run <script>` you'll find the Hardhat
// Runtime Environment's members available in the global scope.
const hre = require("hardhat");

var test = "john"
// scripts/deploy.js
async function main () {
  // We get the contract to deploy
  const Box = await hre.ethers.getContractFactory('Box');
  console.log('Deploying Box...');
  const box = await Box.deploy();
  await box.deployed();
  console.log('Box deployed to:', box.address);
  //console.log('box struct: ', box)
}

main()
  .then(() => process.exit(0))
  .catch(error => {
    console.error(error);
    process.exit(1);
  });