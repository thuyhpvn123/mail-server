// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;
import {EmailStorage} from "./mailStorage.sol";
import {IEmailStorage} from "./interfaces/IEmailStorage.sol";
import "@openzeppelin/contracts@v4.9.0/access/Ownable.sol";
import "@openzeppelin/contracts@v4.9.0/proxy/Clones.sol";
contract FactoryMail is Ownable{
    // Mapping to store the EmailStorage contract for each sender
    mapping(address => address) private senderToEmailStorage;
    // Mapping to store the sender for each EmailStorage contract
    mapping(uint256 => address)public mEmailStorage;
    uint64 public count;
    address public service;
    address public noti;
    address public masterEmailStorage;
    // address public fileSc;
    event EmailStorageCreated(address indexed sender, address emailStorage);
    constructor(address _masterEmailStorage)payable{
        require(_masterEmailStorage != address(0), "Master contract cannot be zero address");
        masterEmailStorage = _masterEmailStorage;
    }
    function setService(address _service)external onlyOwner {
        require(_service != address(0), "Service address cannot be zero");
        service = _service;
        for(uint256 i=1; i<= count;i++){
            _setService(mEmailStorage[i],_service);
        }
    }
    function setNoti(address _noti)external onlyOwner {
        require(_noti != address(0), "Notification address cannot be zero");
        noti = _noti;
        for(uint256 i=1; i< count;i++){
            _setNoti(mEmailStorage[i],_noti);
        }
    }
    function createEmailStorage() public returns(address) {
        require(senderToEmailStorage[msg.sender] == address(0), "EmailStorage already exists for sender");
         // Create a clone proxy
        address clone = Clones.clone(masterEmailStorage);
        senderToEmailStorage[msg.sender] = clone;
        count++;
        mEmailStorage[count] = clone;
        emit EmailStorageCreated(msg.sender, clone);
        return clone;

    }
    //không gọi được trong createEmailStorage vì related address là chính nó
    function init(address clone)external{
        require(service != address(0),"service not set yet");
        require(noti != address(0),"noti not set yet");
        IEmailStorage(clone).initialize(msg.sender, service, noti);
    }

    function getEmailStorageBySender(address sender) external view returns (address) {
        return senderToEmailStorage[sender];
    }
    function _getRevertMsg(bytes memory _returnData) internal pure returns (string memory) {
    // Nếu không có dữ liệu trả về, tức là lỗi không có thông báo
    if (_returnData.length < 68) return "Transaction reverted silently";

    assembly {
        // Bỏ phần prefix 4 bytes của Error Selector
        _returnData := add(_returnData, 0x04)
    }

    // Decode và trả về thông báo lỗi
    return abi.decode(_returnData, (string));
    }

    function _setService(address _emailStorage, address _service) internal {
        require(_emailStorage != address(0), "Invalid EmailStorage address");
        require(_service != address(0), "Invalid service address");

        // Encode the function call to setService
        (bool success, bytes memory data) = _emailStorage.call(
            abi.encodeWithSignature("setService(address)", _service)
        );

        // Nếu call thất bại, parse thông báo lỗi
        if (!success) {
            string memory errorMessage = _getRevertMsg(data);
            revert(errorMessage);
        }
    }

    function _setNoti(address _emailStorage, address _notiSMC) internal {
        require(_emailStorage != address(0), "Invalid EmailStorage address");
        require(_notiSMC != address(0), "Invalid notification address");

        // Encode the function call to setNotificationSMC
        (bool success, bytes memory data) = _emailStorage.call(
            abi.encodeWithSignature("setNotificationSMC(address)", _notiSMC)
        );

        // Nếu call thất bại, parse thông báo lỗi
        if (!success) {
            string memory errorMessage = _getRevertMsg(data);
            revert(errorMessage);
        }
    }


}
