package main

import (
	"bytes"
	"fmt"
	"log"
	"mime/multipart"
	"mime/quotedprintable"
	"net"
	"bufio"
	"os"
)

func main() {
	from := "your-email@example.com"
	to := []string{"ooooo@payws.com"}
	smtpHost := "127.0.0.1"
	// smtpHost := "54.184.97.95"
	smtpPort := "2025"

	// Email headers
	headers := make(map[string]string)
	headers["From"] = from
	headers["From-Header"] = "no-reply@tracking.example.com"
	headers["To"] = to[0]
	headers["Reply-To"] = "support@example.com"
	headers["Subject"] = "HTML Email Test"
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = `multipart/mixed; boundary="BOUNDARY"`
	headers["Message-ID"] = "CAF12345ABCD67@domain.com"
	// Create email body with attachment
	var emailBody bytes.Buffer
	mpWriter := multipart.NewWriter(&emailBody)
	_ = mpWriter.SetBoundary("BOUNDARY")

	// Add HTML part
	htmlWriter, _ := mpWriter.CreatePart(map[string][]string{
		"Content-Type": {"text/html; charset=UTF-8"},
	})
	htmlContent := `<html>
		<head><title>Test Email</title></head>
		<body>
			<h1>Hello, World!</h1>
			<p>This is an email with <b>HTML content</b>.</p>
		</body>
	</html>`
	_, _ = htmlWriter.Write([]byte(htmlContent))

	// Add attachment part
	attachmentWriter, _ := mpWriter.CreatePart(map[string][]string{
		"Content-Type":              {"text/plain; charset=UTF-8"},
		"Content-Disposition":       {`attachment; filename="hello.txt"`},
		"Content-Transfer-Encoding": {"quoted-printable"},
	})
	qpWriter := quotedprintable.NewWriter(attachmentWriter)
	_, _ = qpWriter.Write([]byte("hello trong file dinh kem la hello.txt"))
	_ = qpWriter.Close()

	_ = mpWriter.Close()

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

	// Send HELO command
	fmt.Fprintf(conn, "HELO payws.com\r\n")
	response, err = reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(response)

	// Send MAIL FROM command
	fmt.Fprintf(conn, "MAIL FROM:<%s>\r\n", from)
	response, err = reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(response)
	if response[0] == '5' {
		fmt.Fprintln(os.Stderr, "Error: MAIL FROM command failed -", response)
		return
	}

	// Send RCPT TO command
	fmt.Fprintf(conn, "RCPT TO:<%s>\r\n", to[0])
	response, err = reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(response)
	if response[0] == '5' {
		fmt.Fprintln(os.Stderr, "Error: RCPT TO command failed -", response)
		return
	}

	// Send DATA command
	fmt.Fprintf(conn, "DATA\r\n")
	response, err = reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(response)
	if response[0] == '5' {
		fmt.Fprintln(os.Stderr, "Error: DATA command failed -", response)
		return
	}

	// Send email content
	fmt.Fprintf(conn, "%s\r\n.\r\n", msg)
	response, err = reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(response)
	if response[0] == '5' {
		fmt.Fprintln(os.Stderr, "Error: Sending email content failed -", response)
		return
	}

	// Send QUIT command
	fmt.Fprintf(conn, "QUIT\r\n")
	response, err = reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(response)

	fmt.Println("Email sent successfully")
}
