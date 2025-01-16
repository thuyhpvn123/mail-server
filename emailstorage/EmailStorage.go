// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package emailstorage

import (
	// "math/big"
	"github.com/ethereum/go-ethereum/common"
)
// type Chunk struct {
// 	FileKey [32]byte
// 	ChunkData []byte
// 	ChunkHash [32]byte
// }
type Info struct {
	Owner common.Address
	Hash [32]byte
	ContentLen uint64
	TotalChunks uint64
	ExpireTime uint64
	Name string
	Ext string
	Status uint8
	ContentDisposition string
	ContentID string
}
// type Email struct {
// 	Id         *big.Int
// 	Subject    string
// 	From       string
// 	FromHeader string
// 	ReplyTo    string
// 	MessageID  string
// 	Body       string
// 	Html       string
// 	CreatedAt  *big.Int
// 	Files      []File
// }

// File is an auto generated low-level Go binding around an user-defined struct.
// type File struct {
// 	ContentDisposition string
// 	ContentID          string
// 	ContentType        string
// 	Data               []byte 
// }

