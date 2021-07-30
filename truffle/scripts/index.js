// scripts/index.js
// Set up a Truffle contract, representing our deployed Box instance
const Box = artifacts.require('Box');

module.exports = async function main (callback) {
    try {
        // Retrieve accounts from the local node
        const accounts = await web3.eth.getAccounts();
        console.log(accounts)

        const box = await Box.deployed();

        // Send a transaction to store() a new value in the Box
        await box.store(23);

        // Call the retrieve() function of the deployed Box contract
        const value = await box.retrieve();
        console.log('Box value is', value.toString());
        
        callback(0);
    } catch (error) {
        console.error(error);
        callback(1);
    }
  };