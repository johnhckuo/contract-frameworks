const { expect } = require("chai");
const { ethers } = require("hardhat");

describe("MyNFT contract", function () {

  let MyNFT;
  let contract;
  let owner;
  let accounts;

  before(async function () {
    MyNFT = await ethers.getContractFactory("MyNFT");
    
    [owner, ...accounts] = await hre.ethers.getSigners();

    // deploy and wait for being mined
    contract = await MyNFT.deploy("john's nft", "john");
    await contract.deployed()
    console.log('MyNFT deployed to:', contract.address);
   // contractAddr = contract.address;

    //contract = await MyNFT.attach(contractAddr);

  });

  describe("MyNFT", async function () {
    it("Should return the right name and symbol", async function () {
      expect(await contract.name()).to.equal("john's nft");
      expect(await contract.symbol()).to.equal("john");
    });

    it("is ERC721", async function(){
      expect(await contract.supportsInterface('0x80ac58cd')).to.equal(true)
    })
  });


  describe("Mint", function () {

    it("mint a token from different account", async function(){
      await expect(contract.connect(accounts[0]).mint(accounts[0].address))
      .to.be.revertedWith('MyNFT: must have minter role to mint'); //getting error here
    })

    it("mint a token from minter", async function(){
      expect(await contract.mint(accounts[0].address))
      .to.emit(contract, 'Transfer')
      .withArgs('0x0000000000000000000000000000000000000000', accounts[0].address, 1);
    })

    // If the callback function is async, Mocha will `await` it.
    it("check ownership", async function () {
      expect(await contract.ownerOf(1)).to.equal(accounts[0].address);
    });

    it("Should assign the total supply of tokens to the owner", async function () {
      const ownerBalance = await contract.balanceOf(accounts[0].address);
      expect(await contract.totalSupply()).to.equal(ownerBalance);
    });
  });

  describe("Approve", function () {

    it("check ownership", async function () {
      expect(await contract.ownerOf(1)).to.equal(accounts[0].address);
    });

    it("Approve an operator to transfer on owner's behalf", async function () {

      await expect(contract.connect(accounts[0]).approve(accounts[1].address, 1)).to.emit(contract, 'Approval')
      .withArgs(accounts[0].address, accounts[1].address, 1);
  
    });

    it("Approve from non owner", async function () {

      await expect(contract.connect(accounts[2]).approve(accounts[2].address, 1))
      .to.be.revertedWith('ERC721: approve caller is not owner nor approved for all'); //getting error here
  
    });
  });


  describe("Transfer", function () {

    it("check ownership, should be account[0]", async function () {
      expect(await contract.ownerOf(1)).to.equal(accounts[0].address);
    });

    it("check approved operator, should be account[1]", async function () {
      expect(await contract.getApproved(1)).to.equal(accounts[1].address);
    });

    it("transfer from non operator/owner to another account", async function () {
      await expect(contract.transferFrom(accounts[1].address, accounts[2].address, 1))
      .to.be.revertedWith('ERC721: transfer caller is not owner nor approved'); 
    });

    it("transfer from operator/owner to another account", async function () {
      await expect(contract.connect(accounts[1]).transferFrom(accounts[0].address, accounts[2].address, 1))
      .to.emit(contract, 'Transfer')
      .withArgs(accounts[0].address, accounts[2].address, 1);
    });

    it("check ownership", async function () {
      expect(await contract.ownerOf(1)).to.equal(accounts[2].address);
    });

    it("check balance", async function () {
      expect(await contract.balanceOf(accounts[2].address)).to.equal(1);
    });


    it("Fail so that you can see error message", async function () {
      await expect(contract.connect(accounts[1].address).fail())
      .to.emit(contract, 'Transfer')
      .withArgs(accounts[0].address, accounts[2].address, 1);
    });

  });
});