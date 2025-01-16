package main

import (
	"bufio"
	"bytes"
	"encoding/hex"

	// "encoding/base64"
	"fmt"
	"io"
	"log"
	"mime/multipart"

	// "mime/quotedprintable"
	"net"
	"os"
	"path/filepath"
)

func main() {
	from := "your-email@example.com"
	// to := []string{"ooooo@payws.com"}
	// to := []string{"bbbb@payws.com"}
	to := []string{"aaabbb@payws.com"}

	smtpHost := "127.0.0.1"
	smtpPort := "2025"

	// Email headers
	headers := make(map[string]string)
	headers["From"] = from
	headers["From-Header"] = "no-reply@tracking.example.com"
	headers["To"] = to[0]
	headers["Reply-To"] = "support@example.com"
	headers["Subject"] = "HTML Email with Inline Image and Attachments"
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = `multipart/mixed; boundary="BOUNDARY"`
	headers["Message-ID"] = "CAF12345ABCD67@domain.com"

	// Create email body with inline image and attachment
	var emailBody bytes.Buffer
	mixedWriter := multipart.NewWriter(&emailBody)
	_ = mixedWriter.SetBoundary("BOUNDARY")

	// Create a related part for inline images
	relatedWriter, _ := mixedWriter.CreatePart(map[string][]string{
		"Content-Type": {"multipart/related; boundary=RELATED"},
	})

	relatedMP := multipart.NewWriter(relatedWriter)
	_ = relatedMP.SetBoundary("RELATED")

	// Add HTML part with inline image reference
	htmlPart, _ := relatedMP.CreatePart(map[string][]string{
		"Content-Type": {"text/html; charset=UTF-8"},
	})
	htmlContent := `<html>
		<head><title>Test Email</title></head>
		<body>
			<h1>Gui em thuong nho!</h1>
			<p>This is an email with two inline images and an attachment.</p>
			<img src="cid:image1" alt="Inline Image 1" />
			<img src="cid:image2" alt="Inline Image 2" />
		</body>
	</html>`
	_, _ = htmlPart.Write([]byte(htmlContent))

	// Add first inline image (Referenced by CID)
	err := addInlineImage(relatedMP, "image1", "file2.png")
	if err != nil {
		log.Fatalf("Failed to add inline image1: %v", err)
	}

	// // Add second inline image (Referenced by CID)
	// err = addInlineImage(relatedMP, "image2", "file3.jpg")
	// if err != nil {
	// 	log.Fatalf("Failed to add inline image2: %v", err)
	// }

	_ = relatedMP.Close()

	// // Attach multiple files
	// files := []string{"file1.txt"}
	// for _, file := range files {
	// 	err := addAttachment(mixedWriter, file)
	// 	if err != nil {
	// 		log.Fatalf("Failed to attach file %s: %v", file, err)
	// 	}
	// }

	_ = mixedWriter.Close()

	// Construct the email message
	msg := ""
	for k, v := range headers {
		msg += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	msg += "\r\n" + emailBody.String()

	// Connect to the SMTP server
	conn, err := net.Dial("tcp", smtpHost+":"+smtpPort)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	reader := bufio.NewReader(conn)

	// Read server greeting
	response, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(response)

	// SMTP Commands
	sendSMTPCommand(conn, reader, fmt.Sprintf("HELO payws.com\r\n"))
	sendSMTPCommand(conn, reader, fmt.Sprintf("MAIL FROM:<%s>\r\n", from))
	sendSMTPCommand(conn, reader, fmt.Sprintf("RCPT TO:<%s>\r\n", to[0]))
	sendSMTPCommand(conn, reader, "DATA\r\n")

	// Send email content
	fmt.Fprintf(conn, "%s\r\n.\r\n", msg)
	response, err = reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(response)

	// Quit SMTP
	sendSMTPCommand(conn, reader, "QUIT\r\n")

	fmt.Println("Email sent successfully with inline image and attachment!")
}
// sendSMTPCommand sends an SMTP command and reads the response
func sendSMTPCommand(conn net.Conn, reader *bufio.Reader, cmd string) {
	fmt.Fprintf(conn, cmd)
	response, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(response)
	if response[0] == '5' {
		log.Fatalf("SMTP Error: %s", response)
	}
}
// addInlineImage embeds an image inline using CID with hex encoding
func addInlineImage(mpWriter *multipart.Writer, cid, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	partHeader := map[string][]string{
		"Content-Type":              {"image/png"},
		"Content-Disposition":       {`inline; filename="` + filepath.Base(filePath) + `"`},
		"Content-ID":                {"<" + cid + ">"},
		"Content-Transfer-Encoding": {"x-hex"},
	}

	imageWriter, err := mpWriter.CreatePart(partHeader)
	if err != nil {
		return err
	}

	// Convert file content to hexadecimal
	buffer := make([]byte, 1024) // 4KB chunk size
	for {
		n, err := file.Read(buffer)
		fmt.Println("da chuyen email total chunk:",n)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}
		// hexData := fmt.Sprintf("%x", buffer[:n])
		chunkData := buffer[:n]
		hexData := hex.EncodeToString(chunkData)
		_, err = imageWriter.Write([]byte(hexData))
		if err != nil {
			return err
		}
	}
	
	return nil
}

// addAttachment adds a file as an attachment using hex encoding
func addAttachment(mpWriter *multipart.Writer, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	partHeader := map[string][]string{
		"Content-Type":              {`application/octet-stream`},
		"Content-Disposition":       {fmt.Sprintf(`attachment; filename="%s"`, filepath.Base(filePath))},
		"Content-Transfer-Encoding": {"x-hex"},
	}

	attachmentWriter, err := mpWriter.CreatePart(partHeader)
	if err != nil {
		return err
	}

	// Convert file content to hexadecimal
	buffer := make([]byte, 4096) // 4KB chunk size
	for {
		n, err := file.Read(buffer)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}
		hexData := fmt.Sprintf("%x", buffer[:n])
		_, err = attachmentWriter.Write([]byte(hexData))
		if err != nil {
			return err
		}
	}

	return nil
}

