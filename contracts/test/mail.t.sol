// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "forge-std/Test.sol";
import "../src/mailFactory.sol";
import "../src/mailStorage.sol";
import "../src/interfaces/IEmailStorage.sol";
import "../src/file.sol";
import "../src/Noti.sol";
contract FactoryMailTest is Test {
    FactoryMail factoryMail;
    EmailStorage emailStorage;
    EmailStorage emailStorageMaster;
    Files FILESC;
    NotificationManager NOTI;
    address deployer;
    address service = address(0x123);
    // address NOTI = address(0x124);
    address user1 = address(0x1);
    address user2 = address(0x2);
   
    function setUp() public {
        vm.startPrank(deployer);
        emailStorageMaster = new EmailStorage();
        factoryMail = new FactoryMail(address(emailStorageMaster));
        FILESC = new Files();
        NOTI = new NotificationManager();
        
        factoryMail.setService(service);
        factoryMail.setNoti(address(NOTI));
        FILESC.setService(service);

        vm.stopPrank();
        
    }

    function testCreateMultipleEmailStorages() public {
        // Create first EmailStorage
        vm.startPrank(user1);
        factoryMail.createEmailStorage();

        // Attempt to create a second EmailStorage for the same sender
        vm.expectRevert("EmailStorage already exists for sender");
        factoryMail.createEmailStorage();
        vm.stopPrank();
    }

    function testSetService() public {

        // Verify that service is set correctly
        assertEq(factoryMail.service(), service, "Service address should be set correctly");
        vm.stopPrank();
        vm.startPrank(user1);
        // Create an EmailStorage contract and check if service is set in it
        address emailStorageAddress = factoryMail.createEmailStorage();
        factoryMail.init(emailStorageAddress);
        // Check if service is set in the created EmailStorage contract
        (bool success, bytes memory data) = emailStorageAddress.call(abi.encodeWithSignature("service()"));
        require(success, "Call to get service failed");
        
        assertEq(abi.decode(data, (address)), service, "Service in EmailStorage should match the Factory's service");
        vm.stopPrank();
    }
    function testCreateEmailStorage() public {
        vm.startPrank(user1);
        address emailStorageAddress = factoryMail.createEmailStorage();
        factoryMail.init(emailStorageAddress);
        assertNotEq(emailStorageAddress, address(0), "EmailStorage address should not be zero");
        
        assertEq(factoryMail.getEmailStorageBySender(user1), emailStorageAddress, "Sender should be mapped to their EmailStorage");
        
        vm.stopPrank();
        //add Noti
        // vm.startPrank(deployer);
        // NOTI.addSystemDApp(emailStorageAddress);
        // vm.stopPrank();
        createEmailStorage(emailStorageAddress);
        GetByteCode(emailStorageAddress);
    }
    function createEmailStorage(address emailStorageAddress)public{
        // File[] memory sampleFiles = new File[](1);
        string memory sender = "sender@example.com";
        string memory subject = "Test Subject";
        bytes32[] memory _fileKeys = new bytes32[](1);
        _fileKeys[0] = bytes32(0);
        vm.startPrank(service);
        
        uint256 emailID = IEmailStorage(emailStorageAddress).createEmail(
            sender,
            subject,
            "<p>This is the HTML body of the email</p>",
            _fileKeys,
            1733994682
        );
        bytes32 hashh = bytes32(uint256(123));
        uint64 contentLen = 111111;
        uint64 expireTime = 999999999;
        //push fileinfo
        Info memory info = Info(hashh,contentLen,expireTime);
        bytes32 key = bytes32(uint256(456));
        FILESC.pushFileInfo(info,key);
        //upload file
        bytes32[] memory keys = new bytes32[](1);
        keys[0] = key;
        File memory file = File({
            contentDisposition: "attachment",
            contentID: "12345",
            contentType: "application/pdf",
            data: bytes("day la file pdf nhe")
        });
        File[] memory files = new File[](1);
        files[0] = file;
        FILESC.UploadFiles(keys,files);
        vm.stopPrank();
        //getEmail
        vm.prank(user1);
        // Verify the email creation
        Email memory createdEmail = IEmailStorage(emailStorageAddress).getEmail(emailID);
        assertEq(createdEmail.info.subject, subject, "Subject does not match.");
        assertEq(createdEmail.info.from, sender, "Sender does not match.");
        //get info file
        Info memory infokq = FILESC.GetInfoFile(key);
        assertEq(infokq.hash, hashh);
        //Download file
        File[] memory fileskq = FILESC.DownloadFiles(keys);
        assertEq(fileskq.length, 1, "Files were not added correctly.");
        assertEq(fileskq[0].contentDisposition, "attachment", "File contentDisposition mismatch.");
        vm.stopPrank();
        //Delete file
        //renew time

    }
    function testPushFileInfos() public {
        vm.startPrank(service);

        Info[] memory infos = new Info[](2);
        infos[0] = Info({hash: keccak256("file1"), contentLen: 100, expireTime: block.timestamp + 1 days});
        infos[1] = Info({hash: keccak256("file2"), contentLen: 200, expireTime: block.timestamp + 2 days});

        bytes32[] memory keys = FILESC.pushFileInfos(infos);
        assertEq(keys.length, 2);

        // Verify stored info
        for (uint256 i = 0; i < keys.length; i++) {
            Info memory storedInfo = FILESC.mKeyToFileInfo(keys[i]).info;
            assertEq(storedInfo.hash, infos[i].hash);
            assertEq(storedInfo.contentLen, infos[i].contentLen);
            assertEq(storedInfo.expireTime, infos[i].expireTime);
        }

        vm.stopPrank();
    }

    function testUploadFiles() public {
        vm.startPrank(service);

        Info memory info = Info({hash: keccak256("file1"), contentLen: 100, expireTime: block.timestamp + 1 days});
        bytes32 key = FILESC.pushFileInfo(info);

        File[] memory filesArray = new File[](1);
        filesArray[0] = File({
            contentDisposition: "attachment",
            contentID: "12345",
            contentType: "application/pdf",
            data: bytes("file content")
        });

        bytes32[] memory keys = new bytes32[](1);
        keys[0] = key;

        FILESC.UploadFiles(keys, filesArray);

        // Verify stored file
        File memory storedFile = FILESC.mKeyToFileInfo(key).file;
        assertEq(keccak256(storedFile.data), keccak256(filesArray[0].data));

        vm.stopPrank();
    }

    function testDeleteFile() public {
        vm.startPrank(service);

        Info memory info = Info({hash: keccak256("file1"), contentLen: 100, expireTime: block.timestamp - 1 days});
        FILESC.pushFileInfo(info);

        assert(FILESC.DeleteFile());

        // Verify deletion
        bytes32[] memory keys = new bytes32[](1);
        keys[0] = keccak256(abi.encodePacked(info.contentLen, info.expireTime, info.hash, block.timestamp));
        Info memory deletedInfo = FILESC.mKeyToFileInfo(keys[0]).info;
        assertEq(deletedInfo.hash, 0);

        vm.stopPrank();
    }

    function testUpdateAndRetrieveImage() public {
        bytes memory imageData = bytes("image data");
        address token = address(0x123);

        vm.startPrank(service);
        FILESC.updateImage(imageData, token);
        bytes memory retrievedImage = FILESC.getImage(token);
        assertEq(keccak256(retrievedImage), keccak256(imageData));
        vm.stopPrank();
    }

    function GetByteCode(address  emailStorageAddress)public{
        string memory sender = "sender@example.com";
        string memory subject = "Test Subject";
        bytes32[] memory _fileKeys = new bytes32[](1);
        _fileKeys[0] = bytes32(uint256(123));
        bytes memory bytesCodeCall = abi.encodeCall(
        emailStorage.createEmail,
            (sender,
            subject,
            "<p>This is the HTML body of the email</p>",
            _fileKeys,
            1733994682
            )
        );
        console.log("createEmail:");
        console.logBytes(bytesCodeCall);
        console.log(
            "-----------------------------------------------------------------------------"
        );  
    }
}
