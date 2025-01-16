// SPDX-License-Identifier: MIT
pragma solidity ^0.8.17;

// Khai báo kiểu enum
enum FileStatus {
    Processing, // 0
    Active,     // 1
    Deactive,   // 2
    Deleted     // 3
}

struct Info {
    address owner;
    bytes32 hash;
    uint64 contentLen;
    uint64 totalChunks;
    uint64 expireTime;
    string name;
    string ext;
    FileStatus status;
}

struct FileProgress {
    bytes32 lastChunkHash;
    uint64 processedChunks;
    uint64 processedLength;
}
struct FileInfo {
    Info info;
    FileProgress progress;
    mapping(uint256 => bytes) chunks;
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
