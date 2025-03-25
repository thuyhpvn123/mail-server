// SPDX-License-Identifier: MIT
pragma solidity 0.8.19;
import "./interfaces/IEmailStorage.sol";
contract Files {
    event FileAdded(bytes32 indexed fileKey, string name, uint64 contentLen);
    event ChunkUploaded(bytes32 indexed fileKey, uint256 chunkIndex);
    event FileDeleted(bytes32 indexed fileKey);
    event FileLocked(bytes32 indexed fileKey);

    mapping(bytes32 => FileInfo) public mKeyToFileInfo;
    mapping(string => bytes32) public mNameToFileKey;
    address public service;
    address public owner;

    modifier onlyOwner {
        require(msg.sender == owner, "Caller is not the owner");
        _;
    }

    modifier onlyService {
        require(msg.sender == service, "Caller is not the service");
        _;
    }

    constructor() payable {
        owner = msg.sender;
    }

    function setService(address _service) external onlyOwner {
        service = _service;
    }
    
    function pushFileInfos(
        Info[] memory infos
    )external returns(bytes32[] memory){
        bytes32[] memory keys = new bytes32[](infos.length);
        for(uint256 i; i< infos.length; i++){
            keys[i] = pushFileInfo(infos[i]);
        }
        return keys;
    }

    function pushFileInfo(Info memory info) public returns (bytes32 fileKey) {
        require(info.expireTime > block.timestamp + 1 days, "Expire time must be at least 1 day in the future");
        fileKey = keccak256(
            abi.encodePacked(info.contentLen, info.expireTime, info.hash, info.name, info.ext, block.timestamp)
        );
        mNameToFileKey[info.name] = fileKey;
        require(mKeyToFileInfo[fileKey].info.hash == 0, "File already exists");

        mKeyToFileInfo[fileKey].info = Info({
            owner: msg.sender,
            hash: info.hash,
            contentLen: info.contentLen,
            totalChunks: info.totalChunks,
            expireTime: info.expireTime,
            name: info.name,
            ext: info.ext,
            status: FileStatus.Processing,
            contentDisposition: info.contentDisposition,
            contentID: info.contentID
        });

        mKeyToFileInfo[fileKey].progress = FileProgress({
            lastChunkHash: bytes32(0),
            processedChunks: 0,
            processedLength: 0
        });

        emit FileAdded(fileKey, info.name, info.contentLen);
        return fileKey;
    }
    
    function getFileKeyFromName(string[] memory names) view external returns(bytes32[] memory ){
        bytes32[] memory filekeys = new bytes32[](names.length);
        for (uint256 i;i < names.length; i++){
            filekeys[i] = mNameToFileKey[names[i]];
        }
        return filekeys;
    }
    function uploadChunks(bytes32 fileKey, bytes[] memory chunkDatas, bytes32[] memory chunkHashes) external {
        require(chunkHashes.length == chunkDatas.length,"length arrays should be equal" );
        for(uint256 i; i<chunkDatas.length; i++){
            uploadChunk(fileKey,chunkDatas[i], chunkHashes[i]);
        }
    }
    function uploadChunk(bytes32 fileKey, bytes memory chunkData, bytes32 chunkHash) public returns(bool){
        require(chunkData.length > 0 && chunkData.length <= 10 * 1024, "Invalid chunk size");
        
        FileInfo storage file = mKeyToFileInfo[fileKey];
        require(file.info.owner == msg.sender, "Caller is not the owner");
        require(file.progress.processedLength + chunkData.length <= file.info.contentLen, "Chunk exceeds file length");
        require(file.info.status == FileStatus.Processing, "File is not being processed");
        bytes32 computedChunkHash = keccak256(abi.encodePacked(file.progress.lastChunkHash, chunkData));
        require(computedChunkHash == chunkHash, "Chunk hash mismatch");

        file.progress.lastChunkHash = computedChunkHash;
        file.chunks[file.progress.processedChunks] = chunkData;
        file.progress.processedChunks++;
        file.progress.processedLength += uint64(chunkData.length);
        emit ChunkUploaded(fileKey, file.progress.processedChunks - 1);

        if (file.info.contentLen == file.progress.processedLength || file.progress.processedChunks == file.info.totalChunks) {

            file.info.status = FileStatus.Active;
            delete file.progress;
        }
    }

    function lockFile(bytes32 fileKey) external {
        FileInfo storage file = mKeyToFileInfo[fileKey];
        require(file.info.owner == msg.sender, "Caller is not the owner");
        require(file.info.status == FileStatus.Active, "File is not active");
        require(block.timestamp <= file.info.expireTime, "File has expired");

        file.info.status = FileStatus.Deactive;
        emit FileLocked(fileKey);
    }

    function deleteFile(bytes32 fileKey) external {
        FileInfo storage file = mKeyToFileInfo[fileKey];
        require(file.info.owner == msg.sender, "Caller is not the owner");
        require(file.info.status != FileStatus.Deleted, "File already deleted");

        file.info.status = FileStatus.Deleted;
        for (uint256 i = 0; i < file.info.totalChunks; i++) {
            delete file.chunks[i];
        }
        delete file.progress;
        emit FileDeleted(fileKey);
    }

    function renewTime(bytes32 fileKey, uint64 _newExpireTime) external {
        FileInfo storage file = mKeyToFileInfo[fileKey];
        require(file.info.owner == msg.sender, "Caller is not the owner");
        require(file.info.status != FileStatus.Deleted, "File has been deleted");
        require(_newExpireTime > block.timestamp + 1 days, "New expire time must be at least 1 day in the future");

        file.info.expireTime = _newExpireTime;
    }
    function getFileInfo(bytes32 fileKey) external view returns (Info memory) {
        return mKeyToFileInfo[fileKey].info;
    }
    function getFilesInfo(bytes32[] memory fileKeys) external view returns(Info[] memory infos ){
        infos = new Info[](fileKeys.length);
        for(uint256 i=0; i< fileKeys.length; i++){
            Info memory info = mKeyToFileInfo[fileKeys[i]].info;
            infos[i] = info;
        }
    }

    function getFileProgress(bytes32 fileKey) external view returns (FileProgress memory) {
        FileInfo storage file = mKeyToFileInfo[fileKey];
        require(file.info.status == FileStatus.Processing, "File upload not exists");

        return mKeyToFileInfo[fileKey].progress;
    }
    //     // Hàm tải xuống file theo chunk
    function downloadFile(
        bytes32 fileKey,
        uint256 start,
        uint256 limit
    ) external view returns (bytes[] memory) {
        FileInfo storage file = mKeyToFileInfo[fileKey];
        require(file.info.status == FileStatus.Active, "File upload not actived");
        require(file.info.contentLen > 0, "File does not exist");
        require(start < file.info.totalChunks, "Start index out of range");
        require(limit > 0, "Limit must be greater than zero");
        require(block.timestamp <= file.info.expireTime, "File has expired");

        // Xác định số chunk tối đa có thể trả về
        uint256 end = start + limit > file.info.totalChunks
            ? file.info.totalChunks
            : start + limit;

        // Khởi tạo mảng trả về chunk
        bytes[] memory chunkData = new bytes[](end - start);

        for (uint256 i = start; i < end; i++) {
            chunkData[i - start] = file.chunks[i];
        }

        return chunkData;
    }



}
