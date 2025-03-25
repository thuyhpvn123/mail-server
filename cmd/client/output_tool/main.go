package main

import (
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

const (
	colorReset      = "\033[0m"
	colorRed        = "\033[31m"
	colorGreen      = "\033[32m"
	colorYellow     = "\033[33m"
	colorCyan       = "\033[36m"
	colorFieldName  = "\033[34m" // Blue for field names
	colorFieldValue = "\033[35m" // Magenta for field values
)

func ReadAbi(path string) (*abi.ABI, error) {
	reader, err := os.Open(path)
	if err != nil {
		log.Fatal("Error occurred while reading ABI file")
		return nil, err
	}
	defer reader.Close()

	newABI, err := abi.JSON(reader)
	if err != nil {
		log.Fatal("Error occurred while parsing ABI file")
		return nil, err
	}

	return &newABI, nil
}

func unpackData(ecomABI *abi.ABI, funcName string, data []byte) (interface{}, error) {
	var result interface{}
	err := ecomABI.UnpackIntoInterface(&result, funcName, data)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func unpackDataMap(ecomABI *abi.ABI, funcName string, data []byte) (map[string]interface{}, error) {
	// result := make(map[string]interface{})
	result := make(map[string]interface{})
	err := ecomABI.UnpackIntoMap(result, funcName, data)
	if err != nil {
		return nil, err
	}

	return result, nil
}

type InputData struct {
	ABI      string `json:"abi"`
	Function string `json:"function"`
	Input    string `json:"input"`
	Name     string `json:"name"`
}

func formatToJSON(value interface{}) (string, error) {
	// Convert the value to a reflect.Value
	v := reflect.ValueOf(value)

	// Recursively format the value
	formattedValue := formatValue(v)
	// Convert the value to a JSON string
	jsonBytes, err := json.MarshalIndent(formattedValue, "", "  ")
	if err != nil {
		return "", err
	}
	return string(jsonBytes), nil
}

func addressToHex(address common.Address) string {
	return common.Bytes2Hex(address.Bytes())
}

func array32ToHex(data [32]byte) string {
	return hex.EncodeToString(data[:])
}

// Recursively format values, converting []byte and common.Address to readable strings
func formatValue(v reflect.Value) interface{} {
	if v.Kind() == reflect.Interface {
		if v.IsNil() {
			return nil
		}
		v = v.Elem()
	}
	switch v.Kind() {
	case reflect.Slice:
		elemKind := v.Type().Elem().Kind()
		if elemKind == reflect.Uint8 {
			// Handle byte slices
			return hexToString(v.Interface().([]byte))
		}
		// Handle slices of other types
		formattedSlice := make([]interface{}, v.Len())
		for i := 0; i < v.Len(); i++ {
			formattedSlice[i] = formatValue(v.Index(i))
		}
		return formattedSlice
	case reflect.Array:
		// Handle arrays (similar to slices)
		if v.Type() == reflect.TypeOf(common.Address{}) {
			return addressToHex(v.Interface().(common.Address)) // Convert common.Address to hex string
		}
		if v.Len() == 32 && v.Type().Elem().Kind() == reflect.Uint8 {
			return array32ToHex(v.Interface().([32]byte)) // Convert [32]byte to hex string
		}

		formattedArray := make([]interface{}, v.Len())
		for i := 0; i < v.Len(); i++ {
			formattedArray[i] = formatValue(v.Index(i))
		}
		return formattedArray
	case reflect.Map:
		// Handle maps
		formattedMap := make(map[string]interface{})
		for _, key := range v.MapKeys() {
			formattedKey := formatValue(key)
			formattedValue := formatValue(v.MapIndex(key))
			if keyStr, ok := formattedKey.(string); ok {
				formattedMap[keyStr] = formattedValue
			}
		}
		return formattedMap
	case reflect.Struct:
		// Check if it's a common.Address

		// Handle other structs
		formattedStruct := make(map[string]interface{})
		for i := 0; i < v.NumField(); i++ {
			field := v.Type().Field(i)
			formattedStruct[field.Name] = formatValue(v.Field(i))
		}
		return formattedStruct
	case reflect.String:
		return v.String()
	default:
		// Handle default case
		return v.Interface()
	}
}

func hexToString(data []byte) string {
	decodedStr := string(data)
	decodedStr = strings.Trim(decodedStr, "\x00")
	return decodedStr
}

func writeJSONToFile(directory, filename, jsonData string) error {
	filePath := fmt.Sprintf("%s/%s", directory, filename)
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.WriteString(jsonData)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	jsonFilePath := flag.String("d", "", "Path to the JSON file")
	flag.Parse()

	// Get the current date
	currentDate := time.Now().Format("02-01-2006_15-04-05")
	dirPath := fmt.Sprintf("%s", currentDate)

	// Create a directory with the current date
	err := os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		log.Fatal("Failed to create directory:", err)
	}
	// Read JSON file
	file, err := os.Open(*jsonFilePath)
	if err != nil {
		log.Fatal("Failed to open JSON file:", err)
	}
	defer file.Close()

	var inputData []InputData
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&inputData)
	if err != nil {
		log.Fatal("Failed to decode JSON file:", err)
	}

	if len(inputData) == 0 {
		log.Fatal("No data in JSON file")
	}
	test := 0
	for _, data := range inputData {
		fmt.Println("\n" + colorYellow + "=== Processing Entry ===" + colorReset)
		fmt.Printf(colorGreen+"ABI File: %s\n"+colorReset, data.ABI)
		fmt.Printf(colorGreen+"Function: %s\n"+colorReset, data.Function)
		ecomABI, err := ReadAbi(data.ABI)
		if err != nil {
			log.Printf("Failed to read ABI from %s: %v", data.ABI, err)
			continue
		}
		fmt.Println(colorCyan + "Successfully read ABI from " + data.ABI + colorReset)
		// Decode hex data
		hexData := data.Input
		decodedData, err := hex.DecodeString(hexData)
		if err != nil {
			log.Printf("Invalid input data %s: %v", hexData, err)
			continue
		}
		fmt.Println(colorCyan + "Successfully decoded hex data " + data.Function + colorReset)

		result, err := ecomABI.Unpack(data.Function, decodedData)
		if err != nil {
			log.Printf("Error unpacking data: %v", err)
			continue
		}
		// result, err := unpackData(ecomABI, data.Function, decodedData)
		// if err != nil {
		// 	if strings.Contains(err.Error(), "cannot unmarshal tuple in to interface") {
		// 		result, err = unpackDataMap(ecomABI, data.Function, decodedData)
		// 		if err != nil {
		// 			log.Printf("Error unpacking data: %v", err)
		// 			continue
		// 		}
		// 	} else {
		// 		log.Printf("Error unpacking data: %v", err)
		// 		continue
		// 	}
		// }

		fmt.Println(colorCyan + "Successfully unpacked data" + colorReset)

		// Convert the unpacked result to JSON
		jsonResult, err := formatToJSON(result)
		if err != nil {
			log.Printf("Error formatting result to JSON: %v", err)
			continue
		}
		var outputFile string
		var test1 string
		// Generate a filename based on the Name field
		if data.Name == "" {
			test1 = fmt.Sprintf("%d.json", test)
			outputFile = fmt.Sprintf("%s/%s.json", dirPath, fmt.Sprint(test))
			test++
		} else {
			test1 = fmt.Sprintf("%s.json", data.Name)
			outputFile = fmt.Sprintf("%s/%s.json", dirPath, data.Name)
		}
		// Write JSON to file
		err = writeJSONToFile(dirPath, test1, jsonResult)
		if err != nil {
			log.Printf("Error writing JSON to file: %v", err)
			continue
		}
		fmt.Println(colorGreen + "Successfully wrote JSON to " + outputFile + colorReset)
	}
}
