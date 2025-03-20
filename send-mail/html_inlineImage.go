package main

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"

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
	to := []string{"ooooo@payws.com"}
	// to := []string{"bbbb@payws.com"}
	// to := []string{"nguyenhuytam234@payws.com"}

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

	// Add second inline image (Referenced by CID)
	err = addInlineImage(relatedMP, "image2", "file3.jpg")
	if err != nil {
		log.Fatalf("Failed to add inline image2: %v", err)
	}

	_ = relatedMP.Close()

	// Attach multiple files
	files := []string{"file4.jpeg"}
	for _, file := range files {
		err := addAttachment(mixedWriter, file)
		if err != nil {
			log.Fatalf("Failed to attach file %s: %v", file, err)
		}
	}

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
func addInlineImage(mpWriter *multipart.Writer, cid, filePath string) error {
    file, err := os.Open(filePath)
    if err != nil {
        return fmt.Errorf("failed to open file: %v", err)
    }
    defer file.Close()

    data, err := io.ReadAll(file)
    if err != nil {
        return fmt.Errorf("failed to read file: %v", err)
    }
	writeDataToFile("data12_sendmail",hex.EncodeToString(data))
    base64Data := base64.StdEncoding.EncodeToString(data)

    partHeader := map[string][]string{
        "Content-Type":              {"image/png"},
        "Content-Disposition":       {`inline; filename="` + filepath.Base(filePath) + `"`},
        "Content-ID":                {"<" + cid + ">"},
        "Content-Transfer-Encoding": {"base64"},
    }

    imageWriter, err := mpWriter.CreatePart(partHeader)
    if err != nil {
        return fmt.Errorf("failed to create MIME part: %v", err)
    }

    _, err = imageWriter.Write([]byte(base64Data))
    if err != nil {
        return fmt.Errorf("failed to write base64 data: %v", err)
    }

    return nil
}
func writeDataToFile(filename string, data string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // Pretty-print JSON
	return encoder.Encode(data)
}
// addAttachment adds a file as an attachment using Base64 encoding
func addAttachment(mpWriter *multipart.Writer, filePath string) error {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	// Prepare the MIME part headers
	partHeader := map[string][]string{
		"Content-Type":              {`application/octet-stream`},
		"Content-Disposition":       {fmt.Sprintf(`attachment; filename="%s"`, filepath.Base(filePath))},
		"Content-Transfer-Encoding": {"base64"},
	}

	// Create a part for the attachment
	attachmentWriter, err := mpWriter.CreatePart(partHeader)
	if err != nil {
		return fmt.Errorf("failed to create MIME part: %v", err)
	}

	// Encode the file content in Base64 and write it
	buffer := make([]byte, 4096) // 4KB chunk size
	base64Encoder := base64.NewEncoder(base64.StdEncoding, attachmentWriter)
	defer base64Encoder.Close()

	for {
		n, err := file.Read(buffer)
		if err != nil && err != io.EOF {
			return fmt.Errorf("failed to read file: %v", err)
		}
		if n == 0 { // End of file
			break
		}

		_, err = base64Encoder.Write(buffer[:n])
		if err != nil {
			return fmt.Errorf("failed to write Base64 data: %v", err)
		}
	}

	return nil
}