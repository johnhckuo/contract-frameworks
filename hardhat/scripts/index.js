// scripts/index.js
async function main () {
    // Retrieve accounts from the local node
    const accounts = await ethers.provider.listAccounts();
    console.log(accounts);
    const address = '0x9fE46736679d2D9a65F0992F2272dE9f3c7fa6e0';
    const Box = await ethers.getContractFactory('Box');
    const box = await Box.attach(address);

    // Send a transaction to store() a new value in the Box
    await box.store(23);

    // Call the retrieve() function of the deployed Box contract
    const value = await box.retrieve();
    console.log('Box value is', value.toString());
  }
  
  main()
    .then(() => process.exit(0))
    .catch(error => {
      console.error(error);
      process.exit(1);
    });