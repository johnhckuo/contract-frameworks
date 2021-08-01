// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

// Import Ownable from the OpenZeppelin Contracts library
import "./openzeppelin-contracts/contracts/access/Ownable.sol";
import "./parent.sol";
import "./openzeppelin-contracts/contracts/token/ERC721/ERC721.sol";


contract Amis is Ownable, Parent, ERC721{

    uint256 private _tokenIdTracker = 1;

    constructor(string memory name, string memory symbol) ERC721(name, symbol) {   

        // AccessControl
        // _setupRole(DEFAULT_ADMIN_ROLE, _msgSender());

        // _setupRole(MINTER_ROLE, _msgSender());
    }

    function mint(address to) external onlyOwner returns (uint256)  {
        //require(hasRole(MINTER_ROLE, _msgSender()), "MyNFT: must have minter role to mint");

        uint256 newTokenID = _tokenIdTracker;
        _safeMint(to, newTokenID);
        _tokenIdTracker = _tokenIdTracker + 1;
        return newTokenID;
    }

}
