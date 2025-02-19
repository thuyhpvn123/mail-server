// SPDX-License-Identifier: MIT
// pragma solidity 0.8.19;
// import "./interfaces/IEmailStorage.sol";
// contract Files{
//     mapping(bytes32 => FileInfo) public mKeyToFileInfo;
//     address public service;
//     bytes32[] private allKeys; // Array to store keys
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
//             bytes32 key = pushFileInfo(infos[i]);
//             keys[i] = key;
//         }
//         return keys;
//     }
//     function pushFileInfo(
//         Info memory info
//     ) public onlyService returns(bytes32 key){
//         key = keccak256(abi.encodePacked(info.contentLen,info.expireTime,info.hash,block.timestamp));
//         if (mKeyToFileInfo[key].info.hash == 0) {
//             allKeys.push(key);
//         }
//         FileInfo storage fileinfo = mKeyToFileInfo[key];
//         fileinfo.info.hash = info.hash;
//         fileinfo.info.contentLen = info.contentLen;
//         fileinfo.info.expireTime = info.expireTime;
//     }
//     function GetInfoFile(bytes32[] memory keys) external view returns(Info[] memory infos ){
//         infos = new Info[](keys.length);
//         for(uint256 i=0; i< keys.length; i++){
//             Info memory info = mKeyToFileInfo[keys[i]].info;
//             infos[i] = info;
//         }
//     }
//     function UploadFiles(
//         bytes32[] memory keys,
//         File[] memory files
//     ) external onlyService {
//         require(keys.length == files.length,"keys and files should be equal");
//         for (uint256 i = 0; i < files.length; i++) {
//             //check hash
//             bytes32 hash = hashAttachment(files[i]);
//             FileInfo storage fileinfo = mKeyToFileInfo[keys[i]];
//             require(hash == fileinfo.info.hash,"hash file not equal");
//             fileinfo.file = files[i];
//         }
//     }
//     function hashAttachment(File memory file) internal pure returns (bytes32) {
//         // Concatenate all fields into a single bytes array
//         bytes memory serialized = abi.encodePacked(
//             file.contentDisposition,
//             file.contentID,
//             file.contentType,
//             file.data
//         );

//         // Compute the hash using keccak256
//         return keccak256(serialized);
//     }
//     function DownloadFiles(bytes32[] memory keys) external view returns(FileInfo[] memory){
//         FileInfo[] memory fileInfos = new FileInfo[](keys.length);
//         for(uint256 i; i<keys.length; i++){
//             fileInfos[i] = mKeyToFileInfo[keys[i]];
//         }
//         return fileInfos;
//     }
//     function DeleteFile() external onlyService returns(bool){
//         for (uint256 i = 0; i < allKeys.length; i++) {
//             if(mKeyToFileInfo[allKeys[i]].info.expireTime < block.timestamp){
//                 delete mKeyToFileInfo[allKeys[i]];
//                 allKeys[i] = allKeys[allKeys.length -1];
//                 allKeys.pop();
//             }
//         }        
//         return true;
//     }
//     function RenewTime(bytes32 key, uint64 _newExpireTime) external {
//         mKeyToFileInfo[key].info.expireTime = _newExpireTime;
//     }

//     // Function to update an image
//     function updateImage(bytes calldata imageData,address token) external {
//         require(imageData.length > 0, "Image data cannot be empty");
//         tokenImages[token] = imageData;
//     }

//     // Function to retrieve an image
//     function getImage(address token) external view returns (bytes memory) {
//         bytes memory imageData = tokenImages[token];
//         require(imageData.length > 0, "No image found for this token");
//         return imageData;
//     }
// }