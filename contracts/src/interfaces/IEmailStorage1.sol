// SPDX-License-Identifier: MIT
pragma solidity ^0.8.17;
struct File {
    string contentDisposition;
    string contentID;
    string contentType;
    bytes data;
}
struct Info {
    bytes32 hash;
    uint64 contentLen;
    uint64 expireTime;
}
struct FileInfo {
    Info info;
    File file;
}
struct EmailInfo {
    string subject;
    string from;
    uint64 createdAt; 
    uint8 isRead;
    bytes32[] fileKeys;
}
struct Email {
    EmailInfo info;
    string body;
}
struct UserReceiverEmail {
    uint256 emailID;
    address userAddress;
    bool status;
}
interface IEmailStorage{
    function setService(address _service)external;
    function setStorageOwner(address _storageOwner)external;
    function initialize(address _storageOwner, address _service, address _notification) external ;
    function createEmail(
        string memory sender,
        string memory subject,
        string memory body,
        bytes32[] memory _fileKeys,
        uint64 createdAt
    ) external  returns(uint256);
    function getEmail(uint256 emailID)
    external
    returns (
        Email memory
    );
    function getAllEmailInfos() external view returns(EmailInfo[] memory );
    function getEmailInfos(uint256 startIndex, uint256 count) external view returns (EmailInfo[] memory) ;
    function setNotificationSMC(
        address _notiSMCAddress
    ) external returns (bool);
}
