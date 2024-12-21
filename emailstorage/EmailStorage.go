// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package emailstorage

import (
	"math/big"
)

type Email struct {
	Id         *big.Int
	Subject    string
	From       string
	FromHeader string
	ReplyTo    string
	MessageID  string
	Body       string
	Html       string
	CreatedAt  *big.Int
	Files      []File
}

// File is an auto generated low-level Go binding around an user-defined struct.
type File struct {
	ContentDisposition string
	ContentID          string
	ContentType        string
	Data               []byte 
}

