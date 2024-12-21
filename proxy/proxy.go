package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"
	"bytes"
	"io"
	"strings"
)

const (
	GuerrillaAddr     = "127.0.0.1:2025" // worker chạy trên cổng 2525
	ProcessingTimeout = 30 * time.Second // Timeout tối đa cho mỗi kết nối
)

// Khởi chạy server proxy cho cả IPv4 và IPv6
func main() {
	// Tạo channel để đồng bộ hóa các goroutines
	errCh := make(chan error, 2)

	// Khởi chạy listener IPv4
	go func() {
		errCh <- startServer("tcp4", "0.0.0.0:25")
	}()

	// Khởi chạy listener IPv6
	go func() {
		errCh <- startServer("tcp6", "[::]:25")
	}()

	// Đợi lỗi từ một trong các server
	if err := <-errCh; err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

// Hàm khởi chạy server
func startServer(network, address string) error {
	listener, err := net.Listen(network, address)
	if err != nil {
		return fmt.Errorf("failed to start SMTP proxy server on %s (%s): %w", network, address, err)
	}
	defer listener.Close()

	log.Printf("SMTP Proxy Server running on %s (%s)", address, network)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Connection error on %s (%s): %v", address, network, err)
			continue
		}

		// Xử lý mỗi kết nối trong một goroutine
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {

	log.Printf("handleConnection 1")
	defer conn.Close()

	// Thiết lập timeout cho toàn bộ xử lý kết nối
	ctx, cancel := context.WithTimeout(context.Background(), ProcessingTimeout)
	defer cancel()

	// Kênh để nhận kết quả xử lý
	done := make(chan error, 1)

	go func() {
		// Xử lý logic SMTP
		done <- processSMTP(conn)
	}()

	select {
	case <-ctx.Done():
		// Timeout xảy ra
		log.Println("Connection timeout")
		conn.Write([]byte("421 Connection timeout exceeded\r\n"))
		time.Sleep(1 * time.Second) // Đợi để Gmail nhận mã lỗi
	case err := <-done:
		if err != nil {
			// Lỗi trong xử lý SMTP
			log.Printf("Error during connection handling: %v", err)
			response := fmt.Sprintf("550 Permanent error: backend connection failed (%s)\r\n", err.Error())
			conn.Write([]byte(response))
			time.Sleep(1 * time.Second) // Đợi để Gmail nhận mã lỗi

		}
	}
}

func processSMTP(conn net.Conn) error {
	// Gửi thông điệp chào mừng
	if _, err := conn.Write([]byte("220 Proxy SMTP Server Ready\r\n")); err != nil {
		return fmt.Errorf("failed to send greeting: %w", err)
	}

	// Kết nối đến worker với timeout
	workerConn, err := net.DialTimeout("tcp", GuerrillaAddr, 3*time.Second)
	if err != nil {
		log.Printf("Failed to connect to worker: %v", err)

		// Gửi mã lỗi ngay lập tức
		conn.Write([]byte("550 Permanent error: backend connection failed\r\n"))
		return fmt.Errorf("backend connection failed")
	}
	defer workerConn.Close()

	buffer := make([]byte, 4096)
	dataBuffer := make([]byte, 0) // Buffer lưu trữ dữ liệu trong chế độ DATA
    inDataMode := false          // Đánh dấu đang trong chế độ xử lý lệnh DATA

	for {
		conn.SetReadDeadline(time.Now().Add(ProcessingTimeout))

		// Đọc dữ liệu từ Gmail
		n, err := conn.Read(buffer)
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				log.Println("Timeout occurred while waiting for client")
				conn.Write([]byte("421 Connection timeout exceeded\r\n")) // Gửi lỗi ngay lập tức
				return fmt.Errorf("timeout exceeded while waiting for client: %w", netErr)
			}
			if err == io.EOF {
				log.Println("Client closed connection")
				break
			}
			return fmt.Errorf("failed to read input: %w", err)
		}



		// clientInput := string(buffer[:n])
		// log.Printf("Received from client: %s - %d", clientInput[:6], len(clientInput))


		// Forward tới worker
		// Đọc phản hồi từ worker
		// workerConn.SetReadDeadline(time.Now().Add(ProcessingTimeout))
		// workerConn.SetWriteDeadline(time.Now().Add(ProcessingTimeout))
		 if inDataMode {
            // Xử lý trong chế độ DATA
            dataBuffer = append(dataBuffer, buffer[:n]...)
            if n >= 3 && string(buffer[n-3:n]) == ".\r\n" {
			    log.Println("End of DATA detected.")
			    inDataMode = false

			    // Gửi dữ liệu và kiểm tra lỗi
			    if err := forwardWorker(dataBuffer, workerConn, conn); err != nil {
			        return err
			    }

			    // Reset buffer sau khi gửi dữ liệu
			    dataBuffer = nil
			}
            continue
        }



        // Kiểm tra các lệnh SMTP
        if !inDataMode && bytes.Equal(buffer[:n], []byte("DATA\r\n")) {
            dataBuffer = append(dataBuffer, buffer[:n]...)

            log.Println("Entering DATA mode.")
            conn.Write([]byte("354 Start mail input; end with <CRLF>.<CRLF>\r\n"))
            inDataMode = true
            continue
		}
        // if !inDataMode && clientInput == "DATA\r\n" {
        // }

		if err := forwardWorker(buffer[:n], workerConn, conn); err != nil {
            return err
        }
	}

	return nil
}

func forwardWorker(data []byte, workerConn net.Conn, conn net.Conn) error {
	if _, err := workerConn.Write(data); err != nil {
		log.Printf("Failed to forward to worker: %v", err)
		conn.Write([]byte("550 Permanent error: forwarding failed\r\n")) // Gửi lỗi ngay lập tức
		return fmt.Errorf("failed to forward to worker: %w", err)
	}


	workerBuffer := make([]byte, 1024)
	workerN, err := workerConn.Read(workerBuffer)
	if err != nil {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			conn.Write([]byte("550 Permanent error: worker timeout exceeded\r\n")) // Gửi lỗi ngay lập tức
			return fmt.Errorf("worker timeout exceeded: %w", netErr)
		}
		conn.Write([]byte("550 Permanent error: failed to read from worker\r\n")) // Gửi lỗi ngay lập tức
		return fmt.Errorf("failed to read from worker: %w", err)
	}


	// Gửi phản hồi từ worker về Gmail
	if _, err := conn.Write(workerBuffer[:workerN]); err != nil {
		return fmt.Errorf("failed to send response to client: %w", err)
	}

	workerResponse := string(workerBuffer[:workerN])
	// log.Printf("Worker response 11: %s - %d", workerResponse, len(workerResponse))

	if strings.HasPrefix(workerResponse, "4") || strings.HasPrefix(workerResponse, "5") {
	    conn.Write([]byte("550 Permanent error: worker returned error\r\n")) // Báo lỗi lại cho Gmail

	    parts := strings.Split(workerResponse, "\r\n")
	    return fmt.Errorf("worker returned error: %s", parts[0])
	}

	return nil
}

