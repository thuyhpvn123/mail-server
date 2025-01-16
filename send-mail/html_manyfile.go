package main

// import (
// 	"bufio"
// 	"bytes"
// 	"encoding/base64"
// 	"fmt"
// 	"io"
// 	"log"
// 	"mime/multipart"
// 	// "mime/quotedprintable"
// 	"net"
// 	"os"
// 	"path/filepath"
// )

// func main() {
// 	from := "your-email@example.com"
// 	// to := []string{"ooooo@payws.com"}
// 	to := []string{"bbbb@payws.com"}
// 	smtpHost := "127.0.0.1"
// 	smtpPort := "2025"

// 	// Email headers
// 	headers := make(map[string]string)
// 	headers["From"] = from
// 	headers["From-Header"] = "no-reply@tracking.example.com"
// 	headers["To"] = to[0]
// 	headers["Reply-To"] = "support@example.com"
// 	headers["Subject"] = "HTML Email with Multiple Attachments"
// 	headers["MIME-Version"] = "1.0"
// 	headers["Content-Type"] = `multipart/mixed; boundary="BOUNDARY"`
// 	headers["Message-ID"] = "CAF12345ABCD67@domain.com"

// 	// Create email body with attachment
// 	var emailBody bytes.Buffer
// 	mpWriter := multipart.NewWriter(&emailBody)
// 	_ = mpWriter.SetBoundary("BOUNDARY")

// 	// Add HTML part
// 	htmlWriter, _ := mpWriter.CreatePart(map[string][]string{
// 		"Content-Type": {"text/html; charset=UTF-8"},
// 	})
// 	htmlContent := `<html>
// 		<head><title>Test Email</title></head>
// 		<body>
// 			<h1>Hello, World!</h1>
// 			<p>This is an email with <b>HTML content</b> and multiple attachments.</p>
// 		</body>
// 	</html>`
// 	_, _ = htmlWriter.Write([]byte(htmlContent))

// 	// Attach multiple files
// 	files := []string{"file1.txt", "file2.png"} // List of files to attach
// 	for _, file := range files {
// 		err := addAttachment(mpWriter, file)
// 		if err != nil {
// 			log.Fatalf("Failed to attach file %s: %v", file, err)
// 		}
// 	}

// 	_ = mpWriter.Close()

// 	// Construct the email message
// 	msg := ""
// 	for k, v := range headers {
// 		msg += fmt.Sprintf("%s: %s\r\n", k, v)
// 	}
// 	msg += "\r\n" + emailBody.String()

// 	// Connect to the SMTP server
// 	conn, err := net.Dial("tcp", smtpHost+":"+smtpPort)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer conn.Close()

// 	reader := bufio.NewReader(conn)

// 	// Read server greeting
// 	response, err := reader.ReadString('\n')
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Println(response)

// 	// SMTP Commands
// 	sendSMTPCommand(conn, reader, fmt.Sprintf("HELO payws.com\r\n"))
// 	sendSMTPCommand(conn, reader, fmt.Sprintf("MAIL FROM:<%s>\r\n", from))
// 	sendSMTPCommand(conn, reader, fmt.Sprintf("RCPT TO:<%s>\r\n", to[0]))
// 	sendSMTPCommand(conn, reader, "DATA\r\n")

// 	// Send email content
// 	fmt.Fprintf(conn, "%s\r\n.\r\n", msg)
// 	response, err = reader.ReadString('\n')
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Println(response)

// 	// Quit SMTP
// 	sendSMTPCommand(conn, reader, "QUIT\r\n")

// 	fmt.Println("Email sent successfully with multiple attachments!")
// }

// // // addAttachment adds a file as an attachment to the multipart writer
// // func addAttachment(mpWriter *multipart.Writer, filePath string) error {
// // 	file, err := os.Open(filePath)
// // 	if err != nil {
// // 		return err
// // 	}
// // 	defer file.Close()

// // 	partHeader := map[string][]string{
// // 		"Content-Type":              {`application/octet-stream`},
// // 		"Content-Disposition":       {fmt.Sprintf(`attachment; filename="%s"`, filepath.Base(filePath))},
// // 		"Content-Transfer-Encoding": {"quoted-printable"},
// // 	}

// // 	attachmentWriter, err := mpWriter.CreatePart(partHeader)
// // 	if err != nil {
// // 		return err
// // 	}

// // 	qpWriter := quotedprintable.NewWriter(attachmentWriter)
// // 	defer qpWriter.Close()

// // 	// Manually copy the file content to quotedprintable.Writer
// // 	_, err = io.Copy(qpWriter, file)
// // 	if err != nil {
// // 		return err
// // 	}

// // 	return nil
// // }

// // addAttachment adds a file as an attachment to the multipart writer
// func addAttachment(mpWriter *multipart.Writer, filePath string) error {
// 	file, err := os.Open(filePath)
// 	if err != nil {
// 		return fmt.Errorf("failed to open file %s: %w", filePath, err)
// 	}
// 	defer file.Close()

// 	// Determine content type based on file extension
// 	fileExtension := filepath.Ext(filePath)
// 	contentType := "application/octet-stream" // Default content type
// 	if fileExtension == ".png" {
// 		contentType = "image/png"
// 	}

// 	partHeader := map[string][]string{
// 		"Content-Type":              {fmt.Sprintf("%s; name=\"%s\"", contentType, filepath.Base(filePath))},
// 		"Content-Disposition":       {fmt.Sprintf(`attachment; filename="%s"`, filepath.Base(filePath))},
// 		"Content-Transfer-Encoding": {"base64"},
// 	}

// 	attachmentWriter, err := mpWriter.CreatePart(partHeader)
// 	if err != nil {
// 		return fmt.Errorf("failed to create attachment part for %s: %w", filePath, err)
// 	}

// 	// Encode the file content as base64
// 	encoder := base64.NewEncoder(base64.StdEncoding, attachmentWriter)
// 	defer encoder.Close()

// 	_, err = io.Copy(encoder, file)
// 	if err != nil {
// 		return fmt.Errorf("failed to encode file content for %s: %w", filePath, err)
// 	}

// 	return nil
// }
// // sendSMTPCommand sends an SMTP command and reads the response
// func sendSMTPCommand(conn net.Conn, reader *bufio.Reader, cmd string) {
// 	fmt.Fprintf(conn, cmd)
// 	response, err := reader.ReadString('\n')
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Println(response)
// 	if response[0] == '5' {
// 		log.Fatalf("SMTP Error: %s", response)
// 	}
// }
