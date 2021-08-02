// SPDX-License-Identifier: MIT

pragma solidity ^0.8.0;

import "@openzeppelin/contracts/access/AccessControl.sol";
import "@openzeppelin/contracts/token/ERC721/extensions/ERC721Enumerable.sol";

/**
 * @dev {ERC721} token, including:
 *
 *  - ability for holders to burn (destroy) their tokens
 *  - a minter role that allows for token minting (creation)
 *  - a pauser role that allows to stop all token transfers
 *  - token ID and URI autogeneration
 *
 * This contract uses {AccessControl} to lock permissioned functions using the
 * different roles - head to its documentation for details.
 *
 * The account that deploys the contract will be granted the minter and pauser
 * roles, as well as the default admin role, which will let it grant both minter
 * and pauser roles to other accounts.
 */
contract MyNFT is AccessControl, ERC721Enumerable {
    //using Counters for Counters.Counter;

    bytes32 public constant MINTER_ROLE = keccak256("MINTER_ROLE");

    // set init value to 1
    uint256 private _tokenIdTracker = 1;
    event tokenMinted(uint256, address);

    // for debugging
    event valueChanged(uint256 oldValue, uint256 newValue);

      //create an array of traits to keep track of different traits
    //no traits can be the same
    string[] public trait;
    mapping(string => bool) _traitExists;

    // // Mapping from token ID to owner address
    // mapping(uint256 => address) private _owners;

    // // Mapping owner address to token count
    // mapping(address => uint256) private _balances;

    /**
     * @dev Grants `DEFAULT_ADMIN_ROLE`, `MINTER_ROLE` and `PAUSER_ROLE` to the
     * account that deploys the contract.
     *
     * Token URIs will be autogenerated based on `baseURI` and their token IDs.
     * See {ERC721-tokenURI}.
     */
    constructor(string memory name, string memory symbol) ERC721(name, symbol) {   
        _setupRole(DEFAULT_ADMIN_ROLE, _msgSender());
        _setupRole(MINTER_ROLE, _msgSender());
    }

    function supportsInterface(bytes4 interfaceId) public view virtual override(ERC721Enumerable, AccessControl) returns (bool) {
        return super.supportsInterface(interfaceId);
    }

    /**
     * @dev Creates a new token for `to`. Its token ID will be automatically
     * assigned (and available on the emitted {IERC721-Transfer} event), and the token
     * URI autogenerated based on the base URI passed at construction.
     *
     * See {ERC721-_mint}.
     *
     * Requirements:
     *
     * - the caller must have the `MINTER_ROLE`.
     */
    function mint(address to) public {
        require(hasRole(MINTER_ROLE, _msgSender()), "MyNFT: must have minter role to mint");


        uint256 newTokenID = _tokenIdTracker;
        _safeMint(to, newTokenID);
        _tokenIdTracker = _tokenIdTracker + 1;
        emit tokenMinted(newTokenID, msg.sender);
    }


      //create a new token by calling the mint function
    //every token needs to be different
    //need to pass in a trait that does not exist
    function createNFT(string memory _trait) public {
        
        //check the mapping to determine if trait exists
        require(!_traitExists[_trait]);
        
        require(hasRole(MINTER_ROLE, _msgSender()), "MyNFT: must have minter role to mint");

        //if trait does not exist mint token and create a new token id
        //mint for the msg.sender
        //add trait to array and set trait in mapping to true

        trait.push(_trait);
        uint _id = trait.length - 1;
        _safeMint(msg.sender, _id);
        _traitExists[_trait] = true;
    }

}