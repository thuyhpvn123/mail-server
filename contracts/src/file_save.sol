// // SPDX-License-Identifier: MIT
// pragma solidity 0.8.19;
// import "./interfaces/IEmailStorage.sol";
// contract Files{

//     event FileAdded(bytes32 indexed fileKey, string name, uint64 contentLen);
//     event ChunkUploaded(bytes32 indexed fileKey, uint256 chunkIndex);
//     event FileDeleted(bytes32 indexed fileKey);
//     event FileLocked(bytes32 indexed fileKey);


//     mapping(bytes32 => FileInfo) public mKeyToFileInfo;
//     address public service;
//     address public owner;
//     mapping(address => bytes) public tokenImages;
//     constructor() payable {
//         owner = msg.sender;
//     }
//     modifier onlyOwner{
//         require(msg.sender == owner,"only owner can call");
//         _;
//     }

//     modifier onlyService{
//         require(msg.sender == service,"only service can call");
//         _;
//     }
//     function setService(address _service)external onlyOwner {
//         service = _service;
//     }
//     function pushFileInfos(
//         Info[] memory infos
//     )external returns(bytes32[] memory){
//         bytes32[] memory keys = new bytes32[](infos.length);
//         for(uint256 i; i< infos.length; i++){
//             keys[i] = pushFileInfo(infos[i]);
//         }
//         return keys;
//     }
//     function pushFileInfo(Info memory info) public returns (bytes32 fileKey) {

//         require(info.expireTime > block.timestamp + 24 * 60 * 60, "expire time must be at least 1 day in the future");

//         fileKey = keccak256(
//             abi.encodePacked(
//                 info.contentLen,
//                 info.expireTime,
//                 info.hash,
//                 info.name,
//                 info.ext,
//                 block.timestamp
//             )
//         );

//         require(mKeyToFileInfo[fileKey].info.hash == 0, "File already exists");

//         mKeyToFileInfo[fileKey].info = Info({
//             owner: msg.sender,
//             hash: info.hash,
//             contentLen: info.contentLen,
//             totalChunks: info.totalChunks,
//             expireTime: info.expireTime,
//             name: info.name,
//             ext: info.ext,
//             status: FileStatus.Processing,
//             contentDisposition: info.contentDisposition,
//             contentID: info.contentID

//         });

//         mKeyToFileInfo[fileKey].progress = FileProgress({
//             lastChunkHash: bytes32(0),
//             processedChunks: 0
//         });

//         emit FileAdded(fileKey, info.name, info.contentLen);

//         return fileKey;
//     }

//     // Upload chunk
//     function uploadChunk(bytes32 fileKey, bytes memory chunkData, bytes32 chunkHash) external {

//         // Kiểm tra kích thước chunk
//         require(chunkData.length > 0 && chunkData.length <= 10 * 1024, "Chunk size must be > 0 and <= 10KB");


//         FileInfo storage file = mKeyToFileInfo[fileKey];
//         require(file.info.owner == msg.sender, "Caller is not the owner");

//         require(file.progress.processedLength + chunkData.length <= file.info.contentLen, "Chunk exceeds file length");
//         require(file.info.status == FileStatus.Processing, "File upload not exists");

//         // Tính hash của chunk hiện tại
//         bytes32 computedChunkHash = keccak256(abi.encodePacked(file.progress.lastChunkHash, chunkData));

//         // So sánh với chunkHash được cung cấp
//         require(computedChunkHash == chunkHash, "Chunk hash mismatch");

//         // Cập nhật trạng thái
//         file.progress.lastChunkHash = computedChunkHash;

//         // Lưu chunk
//         file.chunks[file.progress.processedChunks] = chunkData;

//         file.progress.processedChunks++;
//         file.progress.processedLength += chunkData.length;

//         // Đánh dấu hoàn tất nếu upload đủ chunk
//         if (file.info.contentLen == file.progress.processedLength || file.progress.processedChunks == file.info.totalChunks) {
//             file.info.status = FileStatus.Active; // Đặt trạng thái file thành Active

//             // Xóa struct FileProgress
//             delete file.progress;

//             emit ChunkUploaded(fileKey, file.info.totalChunks);
//         } else {
//             emit ChunkUploaded(fileKey, file.progress.processedChunks - 1);
//         }

//     }

//     function getFileInfo(bytes32 fileKey) external view returns (Info memory) {
//         return mKeyToFileInfo[fileKey].info;
//     }

//     function getFilesInfo(bytes32[] memory fileKeys) external view returns(Info[] memory infos ){
//         infos = new Info[](fileKeys.length);
//         for(uint256 i=0; i< fileKeys.length; i++){
//             Info memory info = mKeyToFileInfo[fileKeys[i]].info;
//             infos[i] = info;
//         }
//     }

//     function getFileProgress(bytes32 fileKey) external view returns (FileProgress memory) {
//         FileInfo storage file = mKeyToFileInfo[fileKey];
//         require(file.info.status == FileStatus.Processing, "File upload not exists");

//         return mKeyToFileInfo[fileKey].progress;
//     }

//     // Hàm tải xuống file theo chunk
//     function downloadFile(
//         bytes32 fileKey,
//         uint256 start,
//         uint256 limit
//     ) external view returns (bytes[] memory) {
//         FileInfo storage file = mKeyToFileInfo[fileKey];
//         require(file.info.status == FileStatus.Active, "File upload not actived");
//         require(file.info.contentLen > 0, "File does not exist");
//         require(start < file.info.totalChunks, "Start index out of range");
//         require(limit > 0, "Limit must be greater than zero");
//         require(block.timestamp <= file.info.expireTime, "File has expired");

//         // Xác định số chunk tối đa có thể trả về
//         uint256 end = start + limit > file.info.totalChunks
//             ? file.info.totalChunks
//             : start + limit;

//         // Khởi tạo mảng trả về chunk
//         bytes[] memory chunkData = new bytes[](end - start);

//         for (uint256 i = start; i < end; i++) {
//             chunkData[i - start] = file.chunks[i];
//         }

//         return chunkData;
//     }

//     function lockFile(bytes32 fileKey) external {
//         FileInfo storage file = mKeyToFileInfo[fileKey];
//         require(file.info.owner == msg.sender, "Caller is not the owner");
//         require(file.info.status != FileStatus.Deleted, "File has been deleted");
//         require(file.info.status == FileStatus.Active, "File not active");
//         require(block.timestamp <= file.info.expireTime, "File has expired");

//         file.info.status = FileStatus.Deactive;

//         emit FileLocked(fileKey);
//     }

//     function deleteFile(bytes32 fileKey) external {
//         FileInfo storage file = mKeyToFileInfo[fileKey];
//         require(file.info.owner == msg.sender, "Caller is not the owner");
//         require(file.info.status != FileStatus.Deleted, "File has already been deleted");
//         require(block.timestamp <= file.info.expireTime, "File has expired");

//         file.info.status = FileStatus.Deleted;

//         for (uint256 i = 0; i < file.info.totalChunks; i++) {
//             delete file.chunks[i];
//         }

//         if (file.info.status != FileStatus.Processing) {
//             delete file.progress;
//         } else {

//             file.progress = FileProgress({
//                 lastChunkHash: bytes32(0),
//                 processedChunks: 0
//             });

//         }

//         emit FileDeleted(fileKey);
//     }

//     function renewTime(bytes32 fileKey, uint64 _newExpireTime) external {
//         FileInfo storage file = mKeyToFileInfo[fileKey];
//         require(file.info.owner == msg.sender, "Caller is not the owner");
//         require(file.info.status != FileStatus.Deleted, "File has been deleted");
//         require(_newExpireTime > block.timestamp + 24 * 60 * 60, "New expire time must be at least 1 day in the future");

//         file.info.expireTime = _newExpireTime;
//     }

// }
