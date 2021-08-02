const { expect } = require("chai");
const Box = artifacts.require("Box");
const { ethers } = require("hardhat");

//const { ethereum } = require("hardhat");
// describe("Greeter", function () {
//   it("Should return the new greeting once it's changed", async function () {
//     const Greeter = await ethers.getContractFactory("Greeter");
//     const greeter = await Greeter.deploy("Hello, world!");
//     await greeter.deployed();

//     expect(await greeter.greet()).to.equal("Hello, world!");

//     const setGreetingTx = await greeter.setGreeting("Hola, mundo!");

//     // wait until the transaction is mined
//     await setGreetingTx.wait();

//     expect(await greeter.greet()).to.equal("Hola, mundo!");
//   });
// });

// async function snapshot () {
//   return ethereum.send('evm_snapshot', [])
// }

// async function restore (snapshotId) {
//   return ethereum.send('evm_revert', [snapshotId])
// }

// let snapshotId


// Start test block
describe('Box', function () {
  before(async function () {
    this.Box = await ethers.getContractFactory('Box');
  });

  beforeEach(async function () {
    //snapshotId = await snapshot()
    this.box = await this.Box.deploy();
    await this.box.deployed();
  });

  afterEach(async () => {
    //await restore(snapshotId)
  })

  // Test case
  it('retrieve returns a value previously stored', async function () {
    // Store a value
    await this.box.store(42);

    // Test if the returned value is the same one
    // Note that we need to use strings to compare the 256 bit integers
    expect((await this.box.retrieve()).toString()).to.equal('42');
  });
});



// Traditional Truffle test
contract("Box", (accounts) => {
  it("Should return the new greeting once it's changed", async function () {
    const box = await Box.new();
    await box.store(32);
    assert.equal(await box.retrieve(), 32);
  });
});
