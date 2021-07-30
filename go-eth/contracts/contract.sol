pragma solidity ^0.8.0;

// Import Ownable from the OpenZeppelin Contracts library
import "./openzeppelin-contracts/contracts/access/Ownable.sol";

contract Amis is Ownable{
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
