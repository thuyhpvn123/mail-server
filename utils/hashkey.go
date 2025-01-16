package utils

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	// "encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"golang.org/x/net/html"
)

func GenerateFileHash(data []byte) [32]byte {
	// Prepend 0x00 to the data
	prefixedData := append([]byte{0x00}, data...)

	// Compute Keccak-256 hash
	hashBytes := crypto.Keccak256(prefixedData)

	// Convert the hash to a [32]byte array
	var hash [32]byte
	copy(hash[:], hashBytes)
	return hash
}
// Calculate keccak256 hash (equivalent to abi.encodePacked in Solidity)
func keccak256Hash(data []byte) [32]byte {
	return crypto.Keccak256Hash(data)
}

// Calculate chunk hash (keccak256(lastChunkHash + chunkData))
func CalculateChunkHash(lastChunkHash [32]byte, chunkData []byte) [32]byte {
	// Concatenate lastChunkHash and chunkData
	combinedData := append(lastChunkHash[:], chunkData...)
	return keccak256Hash(combinedData)
}

// Function to chunk a file into smaller pieces
func ChunkData(data []byte, chunkSize int) ([][]byte,int) {
	var chunks [][]byte
	for i := 0; i < len(data); i += chunkSize {
		end := i + chunkSize
		if end > len(data) {
			end = len(data)
		}
		chunks = append(chunks, data[i:end])
	}
	return chunks,len(chunks)
	// // Convert the image data to a hexadecimal string
	// hexData := hex.EncodeToString(data)

	// // Split the hexadecimal string into chunks
	// var hexChunks []string
	// for i := 0; i < len(hexData); i += chunkSize {
	// 	end := i + chunkSize
	// 	if end > len(hexData) {
	// 		end = len(hexData)
	// 	}
	// 	hexChunks = append(hexChunks, hexData[i:end])
	// }

	// return hexChunks, len(hexChunks)

}
func ExtractFileName(contentDisposition string) string {
	parts := strings.Split(contentDisposition, ";")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if strings.HasPrefix(part, "filename=") {
			fileName := strings.TrimPrefix(part, "filename=")
			fileName = strings.Trim(fileName, "\"") // Remove surrounding quotes
			return fileName
		}
	}
	return "unknown"
}
func ExtractFileExtension(fileName string) string {
	ext := filepath.Ext(fileName)
	if ext != "" {
		return ext[1:] // Remove the leading dot
	}
	return "unknown"
}
func CalculateChunks(contentLength, chunkSize int) int {
	if contentLength%chunkSize == 0 {
		return contentLength / chunkSize
	}
	return contentLength/chunkSize + 1
}
func GetMax200Characters(input string, maxLength int) string {
	// Convert the string to a rune slice to handle Unicode characters properly
	runes := []rune(input)

	if len(runes) > maxLength {
		return string(runes[:maxLength])
	}
	return input
}
// extractBodyContent extracts the content of the <body> tag from an HTML string.
func ExtractBodyContent(htmlString string) (string, error) {
	doc, err := html.Parse(strings.NewReader(htmlString))
	if err != nil {
		return "", fmt.Errorf("failed to parse HTML: %v", err)
	}

	var bodyBuffer bytes.Buffer
	var foundBody bool

	var traverse func(*html.Node)
	traverse = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "body" {
			foundBody = true
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				html.Render(&bodyBuffer, c)
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}
	}

	traverse(doc)

	if !foundBody {
		return "", fmt.Errorf("<body> tag not found in the HTML")
	}

	return bodyBuffer.String(), nil
}
// Hàm tạo password bằng ECDSA
func GeneratePassword(email string) (string, error) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return "", err
	}

	hash := sha256.Sum256([]byte(email))
	r, s, err := ecdsa.Sign(rand.Reader, privateKey, hash[:])
	if err != nil {
		return "", err
	}

	password := fmt.Sprintf("%x%x", r, s)
	return password, nil
}
// Hàm lưu email xuống ổ cứng tại server
func SaveEmailLocally(encryptedEmail []byte) error {
	// Define the directory path
	dirPath := "./data"

	// Check if the directory exists, and create it if not
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
	}

	// Generate the file name
	filename := fmt.Sprintf("email_%d.txt.gz", time.Now().UnixNano())
	filePath := fmt.Sprintf("%s/%s", dirPath, filename)

	// Create the file
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// Write the encrypted email to the file
	if _, err := file.Write(encryptedEmail); err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	return nil
}

