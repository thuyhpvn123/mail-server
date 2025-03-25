// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;
// import "@openzeppelin/contracts@v4.9.0/access/Ownable.sol";
import "./interfaces/IEmailStorage.sol";
import "./interfaces/INoti.sol";
contract EmailStorage {
    mapping(uint256 => Email) private mIDToEmail;
    address public notiSMC;
    address public service;
    address public storageOwner;
    address public fileSc;
    mapping(address => bool) public isOwner;
    string public REPO_NOTI_SMC = "MAIL";

    // Event for email creation
    event EmailCreated(uint256 emailID, string subject, string indexed creator);
    event NotificationFailed(string reason);
    uint256 public emailCounter;
    modifier onlyOwner {
        require(isOwner[msg.sender],"only owner");
        _;
    }
    constructor()payable{}
    function initialize(address _storageOwner, address _service, address _notification) external {
        require(storageOwner == address(0), "Already initialized"); // Prevent reinitialization
        storageOwner = _storageOwner;
        service = _service;
        notiSMC = _notification;
        isOwner[_storageOwner] = true;
        isOwner[msg.sender] = true;
    }

    modifier onlyService{
        require(msg.sender == service,"only service can call");
        _;
    }
    modifier onlyEmailStorageOwner{
        require(msg.sender == storageOwner,"only Email Storage Owner can call");
        _;
    }
    function setService(address _service)external onlyOwner {
        service = _service;
    }
    function setStorageOwner(address _storageOwner)external onlyOwner {
        storageOwner = _storageOwner;
    }
    function setNotificationSMC(
        address _notiSMCAddress
    ) external onlyOwner {
        notiSMC = _notiSMCAddress;
    }
    function setFileSMC(address _fileSc)external onlyOwner {
        fileSc = _fileSc;
    }

    function createEmail(
        string memory sender,
        string memory subject,
        string memory body,
        bytes32[] memory _fileKeys,
        uint64 createdAt,
        string memory discription
    ) external onlyService returns(uint256) {
        // Store the email
        Email storage  newEmail = mIDToEmail[emailCounter];
        newEmail.info.subject = subject;
        newEmail.info.from = sender;
        newEmail.body = body;
        newEmail.info.createdAt = createdAt;
        newEmail.info.discription = discription;
        // Add files to the email
        for (uint256 i = 0; i < _fileKeys.length; i++) {
           newEmail.info.fileKeys.push(_fileKeys[i]);
        }
        emailCounter++;
        emit EmailCreated(emailCounter, subject, sender);
        if (address(notiSMC) != address(0)) {
            NotiParams memory notiParams;
            notiParams.title = "Metanode-Email";
            notiParams.body = subject;
        try INoti(notiSMC).AddNoti(notiParams, storageOwner) {
            // Notification sent successfully+
        } catch Error(string memory reason) {
            // Handle known errors
            emit NotificationFailed(reason);
        } catch (bytes memory /*lowLevelData*/) {
            // Handle unknown errors
            emit NotificationFailed("Notification failed due to unknown error.");
        }
    
    }
        return (emailCounter-1);
    }
    function getEmail(uint256 emailID)
        external
        onlyEmailStorageOwner
        returns (
            Email memory
        )
    {
        require(emailID >= 0 && emailID <= emailCounter, "Email does not exist");
        Email storage email = mIDToEmail[emailID];
        if (email.info.isRead == 0) {
            email.info.isRead = 1;
        }
        return mIDToEmail[emailID];
    }
    function getAllEmailInfos() external view onlyEmailStorageOwner returns(EmailInfo[] memory ){
        require(emailCounter>0,"there is not any email");
        EmailInfo[] memory emails = new EmailInfo[](emailCounter);
        for(uint256 i; i<emailCounter; i++){
            emails[i] = mIDToEmail[i].info;
        }
        return emails;
    }
    function getEmailInfos(uint256 startIndex, uint256 count) external view returns (EmailInfo[] memory) {
        require(startIndex < emailCounter, "Start index out of range");
        require(count < 20, "Count out of range");
        uint256 endIndex = startIndex + count;
        if (endIndex > emailCounter) {
            endIndex = emailCounter;
        }

        EmailInfo[] memory result = new EmailInfo[](endIndex - startIndex);
        for (uint256 i = startIndex; i < endIndex; i++) {
            result[i - startIndex] = mIDToEmail[i].info;
        }
        return result;
    }

}