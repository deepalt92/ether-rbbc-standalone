pragma solidity >=0.4.24;

contract NameContract {

    string private name = "Ire";

    function getName() public view returns (string memory)
    {
        return name;
    }

    function setName(string memory newName) public
    {
        name = newName;
    }

}