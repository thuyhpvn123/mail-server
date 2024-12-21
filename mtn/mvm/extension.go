package mvm

/*
#cgo CFLAGS: -w
#cgo CXXFLAGS: -std=c++17 -w
#cgo LDFLAGS: -L./linker/build/lib/static -lmvm_linker -L./c_mvm/build/lib/static -lmvm -lstdc++
#cgo CPPFLAGS: -I./linker/build/include
#include "mvm_linker.hpp"
#include <stdlib.h>
*/
import "C"
import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"unsafe"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"gomail/mtn/bls"
	"gomail/mtn/logger"
	"gomail/mtn/common"
	"gomail/mtn/smart_contract/argument_encode"
)

//export ExtensionCallGetApi
func ExtensionCallGetApi(
	bytes *C.uchar,
	size C.int,
) (
	data_p *C.uchar,
	data_size C.int,
) {
	bCallData := C.GoBytes(unsafe.Pointer(bytes), size)
	logger.Debug("Calling get api data ", hex.EncodeToString(bCallData))
	url := argument_encode.DecodeStringInput(bCallData[4:], 0)
	response, err := http.Get(url)
	if err != nil {
		logger.Warn("Error when call get api to ", url, err)
		return
	}
	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		logger.Warn("Error when call get api to ", url, err)
		return
	}
	encodedRespone := argument_encode.EncodeSingleString(string(responseData))
	logger.Debug("Extension call get api result ", encodedRespone)
	data_size = C.int(len(encodedRespone))
	data_p = (*C.uchar)(C.CBytes(encodedRespone))
	return
}

//export ExtensionExtractJsonField
func ExtensionExtractJsonField(
	bytes *C.uchar,
	size C.int,
) (
	data_p *C.uchar,
	data_size C.int,
) {
	bCallData := C.GoBytes(unsafe.Pointer(bytes), size)
	logger.Debug("Extension extract json field ", hex.EncodeToString(bCallData))
	jsonMap := make(map[string]interface{})
	jsonStr := argument_encode.DecodeStringInput(bCallData[4:], 0)
	field := argument_encode.DecodeStringInput(bCallData[4:], 1)
	var fieldData interface{}
	err := json.Unmarshal([]byte(jsonStr), &jsonMap)
	var data string
	// process json map
	if err == nil {
		fieldData = jsonMap[field]
	} else {
		// process json array
		jsonArr := []interface{}{}
		err = json.Unmarshal([]byte(jsonStr), &jsonArr)
		if err != nil {
			logger.Warn("Error when extract json field ", jsonStr, field, err)
			return
		}
		intField, err := strconv.Atoi(field)
		if err != nil {
			logger.Warn("Error when extract json field ", jsonStr, field, err)
			return
		}
		fieldData = jsonArr[intField]
	}

	if reflect.ValueOf(fieldData).Kind() == reflect.Map || reflect.ValueOf(fieldData).Kind() == reflect.Array {
		bData, _ := json.Marshal(fieldData)
		data = string(bData)
	} else {
		data = fmt.Sprintf("%v", fieldData)
		// reformat boolean
		if data == "false" {
			data = "0"
		}
		if data == "true" {
			data = "1"
		}
	}

	encodedData := argument_encode.EncodeSingleString(data)
	data_size = C.int(len(encodedData))
	data_p = (*C.uchar)(C.CBytes(encodedData))
	return
}


//export ExtensionBlst
func ExtensionBlst(
	bytes *C.uchar,
	size C.int,
) (
	data_p *C.uchar,
	data_size C.int,
) {
  blstAbiStr :=  strings.NewReader(`
[
	{
		"inputs": [
			{
				"internalType": "bytes[]",
				"name": "publicKey",
				"type": "bytes[]"
			},
			{
				"internalType": "bytes",
				"name": "sign",
				"type": "bytes"
			},
			{
				"internalType": "bytes[]",
				"name": "message",
				"type": "bytes[]"
			}
		],
		"name": "verifyAggregateSign",
		"outputs": [
			{
				"internalType": "bool",
				"name": "",
				"type": "bool"
			}
		],
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"inputs": [
			{
				"internalType": "bytes",
				"name": "publicKey",
				"type": "bytes"
			},
			{
				"internalType": "bytes",
				"name": "sign",
				"type": "bytes"
			},
			{
				"internalType": "bytes",
				"name": "message",
				"type": "bytes"
			}
		],
		"name": "verifySign",
		"outputs": [
			{
				"internalType": "bool",
				"name": "",
				"type": "bool"
			}
		],
		"stateMutability": "nonpayable",
		"type": "function"
	}
]
  `)
  blstAbi, err := abi.JSON(blstAbiStr)
  if err != nil {
    logger.Error("Error ", err)
    return
  }

	bCallData := C.GoBytes(unsafe.Pointer(bytes), size)
	logger.Debug("Calling extention blst", hex.EncodeToString(bCallData))
  
  method, err := blstAbi.MethodById(bCallData[0:4])
  if err != nil {
    logger.Warn("Error when get method by id", err)
    return
  }

  if method.RawName == "verifySign" {
    mapInput := make(map[string]interface{})
    method.Inputs.UnpackIntoMap(mapInput, bCallData[4:])
    outputs, err := method.Outputs.Pack(
      bls.VerifySign(
        common.PubkeyFromBytes(mapInput["publicKey"].([]byte)), 
        common.SignFromBytes(mapInput["sign"].([]byte)), 
        mapInput["message"].([]byte),
       ),
    )
    if err != nil {
      logger.Warn("Error when pack output", err)
    }
	  data_size = C.int(len(outputs))
	  data_p = (*C.uchar)(C.CBytes(outputs))
  }

  if method.RawName == "verifyAggregateSign" {
// VerifyAggregateSign(bPubs [][]byte, bSig []byte, bMsgs [][]byte) bool
    mapInput := make(map[string]interface{})
    method.Inputs.UnpackIntoMap(mapInput, bCallData[4:])
    outputs, err := method.Outputs.Pack(
      bls.VerifyAggregateSign(
        mapInput["publicKey"].([][]byte), 
        mapInput["sign"].([]byte), 
        mapInput["message"].([][]byte),
       ),
    )
    if err != nil {
      logger.Warn("Error when pack output", err)
    }
	  data_size = C.int(len(outputs))
	  data_p = (*C.uchar)(C.CBytes(outputs))
    logger.Debug("Extension blst result ", hex.EncodeToString(outputs))
  }
  return
}

func WrapExtensionBlst(
  data []byte,
) ( []byte) {
  cData := C.CBytes(data)
  cReturn, returnSize:= ExtensionBlst((*C.uchar)(cData), C.int(len(data)))
  return C.GoBytes(unsafe.Pointer(cReturn), returnSize)
}
