// contracts/Box.sol
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

// Import Ownable from the OpenZeppelin Contracts library
import "@openzeppelin/contracts/access/Ownable.sol";


contract Box is Ownable{
    uint256 private _value;

    event ValueChanged(uint256 value);

    function store(uint256 value) public onlyOwner {
        _value = value;
        emit ValueChanged(value);
    }

    function retrieve() public view returns (uint256) {
        return _value;
    }
}

// contract Greeter {
//   string greeting;

//   constructor(string memory _greeting) {
//     greeting = _greeting;
//   }

//   function greet() public view returns (string memory) {
//     return greeting;
//   }

//   function setGreeting(string memory _greeting) public {
//     greeting = _greeting;
//   }
// }
