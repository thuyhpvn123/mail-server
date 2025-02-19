// SPDX-License-Identifier: MIT
pragma solidity ^0.8.17;

import "forge-std/Test.sol";
import "../src/file.sol";

contract FilesTest is Test {
    Files private files;
    address private owner;
    address private service;
    bytes32 private fileHash;
    
    function setUp() public {
        owner = address(this);
        service = address(0x123);
        files = new Files();
        files.setService(service);
        GetByteCode();
    }

    function testPushFileInfo() public {
        Info memory info = Info({
            owner: owner,
            hash: keccak256(abi.encodePacked("file1")),
            contentLen: 1024,
            totalChunks: 2,
            expireTime: uint64(block.timestamp + 2 days),
            name: "testFile",
            ext: "txt",
            status: FileStatus.Processing,
            contentDisposition: "inline",
            contentID: "12345"
        });

        bytes32 fileKey = files.pushFileInfo(info);
        Info memory storedInfo = files.getFileInfo(fileKey);

        assertEq(storedInfo.name, "testFile");
        assertEq(storedInfo.ext, "txt");
        assertEq(storedInfo.contentLen, 1024);
    }
    function testPushFileInfos() public {
        Info[] memory infos = new Info[](2);
        
        infos[0] = Info({
            owner: address(this),
            hash: keccak256("file1"),
            contentLen: 1024,
            totalChunks: 4,
            expireTime: uint64(block.timestamp + 2 days),
            name: "file1",
            ext: "txt",
            status: FileStatus.Processing,
            contentDisposition: "inline",
            contentID: "cid1"
        });
        
        infos[1] = Info({
            owner: address(this),
            hash: keccak256("file2"),
            contentLen: 2048,
            totalChunks: 8,
            expireTime: uint64(block.timestamp + 3 days),
            name: "file2",
            ext: "jpg",
            status: FileStatus.Processing,
            contentDisposition: "attachment",
            contentID: "cid2"
        });
        
        bytes32[] memory fileKeys = files.pushFileInfos(infos);

        assertEq(fileKeys.length, 2);
        assertTrue(fileKeys[0] != bytes32(0));
        assertTrue(fileKeys[1] != bytes32(0));
    }
    function testUploadChunk() public {
        Info memory info = Info({
            owner: owner,
            hash: keccak256(abi.encodePacked("file1")),
            contentLen: 2048,
            totalChunks: 2,
            expireTime: uint64(block.timestamp + 2 days),
            name: "testFile",
            ext: "txt",
            status: FileStatus.Processing,
            contentDisposition: "inline",
            contentID: "12345"
        });

        bytes32 fileKey = files.pushFileInfo(info);
        bytes memory chunkData = "Chunk1Data";
        bytes32 chunkHash = keccak256(abi.encodePacked(bytes32(0), chunkData));

        files.uploadChunk(fileKey, chunkData, chunkHash);
        FileProgress memory progress = files.getFileProgress(fileKey);

        assertEq(progress.processedChunks, 1);
    }

    function testDeleteFile() public {
        Info memory info = Info({
            owner: owner,
            hash: keccak256(abi.encodePacked("file1")),
            contentLen: 2048,
            totalChunks: 2,
            expireTime: uint64(block.timestamp + 2 days),
            name: "testFile",
            ext: "txt",
            status: FileStatus.Processing,
            contentDisposition: "inline",
            contentID: "12345"
        });

        bytes32 fileKey = files.pushFileInfo(info);
        files.deleteFile(fileKey);
        Info memory storedInfo = files.getFileInfo(fileKey);

        assertEq(uint8(storedInfo.status), uint8(FileStatus.Deleted));
    }

    function testRenewTime() public {
        Info memory info = Info({
            owner: owner,
            hash: keccak256(abi.encodePacked("file1")),
            contentLen: 2048,
            totalChunks: 2,
            expireTime: uint64(block.timestamp + 2 days),
            name: "testFile",
            ext: "txt",
            status: FileStatus.Processing,
            contentDisposition: "inline",
            contentID: "12345"
        });

        bytes32 fileKey = files.pushFileInfo(info);
        uint64 newExpireTime = uint64(block.timestamp + 5 days);
        files.renewTime(fileKey, newExpireTime);
        Info memory storedInfo = files.getFileInfo(fileKey);

        assertEq(storedInfo.expireTime, newExpireTime);
    }
    function GetByteCode()public{
         Info memory info = Info({
            owner: owner,
            hash: keccak256(abi.encodePacked("file1")),
            contentLen: 1024,
            totalChunks: 2,
            expireTime: uint64(1739935067 + 2 days),
            name: "testFile",
            ext: "txt",
            status: FileStatus.Processing,
            contentDisposition: "inline",
            contentID: "12345"
        });
        console.log("expire time:",info.expireTime);
        bytes memory bytesCodeCall = abi.encodeCall(
        files.pushFileInfo,
            (info)
        );
        console.log("pushFileInfo:");
        console.logBytes(bytesCodeCall);
        console.log(
            "-----------------------------------------------------------------------------"
        );  
        Info[] memory infos = new Info[](2);
        infos[0] = Info({
            owner: address(this),
            hash: keccak256("file1"),
            contentLen: 1024,
            totalChunks: 4,
            expireTime: uint64(1739935067 + 2 days),
            name: "file1",
            ext: "txt",
            status: FileStatus.Processing,
            contentDisposition: "inline",
            contentID: "cid1"
        });
        
        infos[1] = Info({
            owner: address(this),
            hash: keccak256("file2"),
            contentLen: 2048,
            totalChunks: 8,
            expireTime: uint64(1739935067 + 3 days),
            name: "file2",
            ext: "jpg",
            status: FileStatus.Processing,
            contentDisposition: "attachment",
            contentID: "cid2"
        });
        bytesCodeCall = abi.encodeCall(
        files.pushFileInfos,
            (infos)
        );
        console.log("pushFileInfos:");
        console.logBytes(bytesCodeCall);
        console.log(
            "-----------------------------------------------------------------------------"
        );  
    }

}
