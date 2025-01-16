package main

import (
	"bytes"
	"compress/gzip"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime"
	"net"
	"net/http"
	"net/mail"
	"net/url"
	"os"
	"os/exec"
	"reflect"

	// "regexp"
	"strings"
	"time"

	"github.com/emersion/go-msgauth/dmarc"
	"github.com/ethereum/go-ethereum/common"

	// "github.com/microcosm-cc/bluemonday"
	"github.com/miekg/dns"
	"github.com/phires/go-guerrilla"
	"github.com/phires/go-guerrilla/backends"
	guerrillaMail "github.com/phires/go-guerrilla/mail"
	"github.com/toorop/go-dkim"
	"golang.org/x/net/html"

	// "math/big"
	"gomail/cmd/client"
	c_config "gomail/cmd/client/pkg/config"
	"gomail/config"
	"gomail/emailstorage"
	"gomail/services"
	"gomail/utils"
	"mime/multipart"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi"
)
const chunkSize = 1 * 1024
// File to monitor
const monitoredFile = "/home/ubuntu/gitMain" // Change to the path of your executable

// Time interval for polling
const pollingInterval = 3 * time.Second

const maxDesLength = 200

var lastModTime time.Time

var (
	ChainClient     *client.Client
	emailStorageMap = sync.Map{}
)

type ParsedEmail struct {
	Subject     string
	Body        string       // Nội dung body
	Attachments []Attachment // Mảng các attachment
	Discription string
	// ReplyTo 	string
	// MessageId   string
	// FromHeader  string
	// Html		string
}

type Attachment struct {
	ContentDisposition string // Content-Disposition (ví dụ: "attachment")
	ContentID          string // Content-ID nếu có
	ContentType        string // Content-Type của file (ví dụ: "application/pdf")
	Data               []byte // Dữ liệu của attachment
	FileName string
	FileExtension string
	FileHash [32]byte
	ContentLength uint64
	TotalChunks uint64
}

// func sanitizeEmailHTML(html string) string {
// 	// Create a Bluemonday policy based on UGCPolicy
// 	policy := bluemonday.UGCPolicy()

// 	// Allow basic HTML email structure
// 	policy.AllowElements("html", "head", "body", "label", "input", "font", "main", "nav", "header", "footer", "kbd", "legend", "map", "title", "div", "span")

// 	// Allow basic style attributes globally
// 	policy.AllowAttrs("style").Globally()

// 	// Custom attributes for specific tags
// 	policy.AllowAttrs("face", "size").OnElements("font")
// 	policy.AllowAttrs("name", "content", "http-equiv").OnElements("meta")
// 	policy.AllowAttrs("name", "data-id").OnElements("a")
// 	policy.AllowAttrs("for").OnElements("label")
// 	policy.AllowAttrs("type").OnElements("input")
// 	policy.AllowAttrs("rel", "href").OnElements("link")
// 	policy.AllowAttrs("topmargin", "leftmargin", "marginwidth", "marginheight", "yahoo").OnElements("body")
// 	policy.AllowAttrs("xmlns").OnElements("html")
// 	policy.AllowAttrs("class", "id", "style").OnElements("div", "span")
// 	policy.AllowAttrs("style", "vspace", "hspace", "face", "bgcolor", "color", "border", "cellpadding", "cellspacing").Globally()

// 	// Support common email tags and attributes
// 	policy.AllowAttrs("bgcolor", "color", "align").OnElements("basefont", "font", "hr", "table", "td")
// 	policy.AllowAttrs("border").OnElements("img", "table", "basefont", "font", "hr", "td")
// 	policy.AllowAttrs("cellpadding", "cellspacing", "valign", "halign").OnElements("table")

// 	// Handle <img> tags and inline images
// 	policy.AllowAttrs("src", "alt", "width", "height", "style").OnElements("img")

// 	// Allow `cid:` URIs specifically for inline images
// 	policy.AllowAttrs("src").Matching(regexp.MustCompile(`^cid:[a-zA-Z0-9._%+-]+$`)).OnElements("img")

// 	// Allow safe `data:image` URIs
// 	policy.AllowAttrs("src").Matching(regexp.MustCompile(`^data:image\/(png|jpeg|gif|webp);base64,`)).OnElements("img")

// 	// Allow `http://` and `https://` for external images (optional)
// 	policy.AllowAttrs("src").Matching(regexp.MustCompile(`^https?://`)).OnElements("img")

// 	// Add "nofollow" and "target=_blank" to links
// 	policy.RequireNoFollowOnLinks(true)
// 	policy.RequireNoFollowOnFullyQualifiedLinks(true)
// 	policy.AddTargetBlankToFullyQualifiedLinks(true)

// 	// Sanitize and return the cleaned HTML
// 	return policy.Sanitize(html)
// }
func sanitizeEmailHTML(htmlString string) (string, error) {
	// Parse the HTML string
	doc, err := html.Parse(strings.NewReader(htmlString))
	if err != nil {
		return "", fmt.Errorf("failed to parse HTML: %v", err)
	}

	// Traverse the HTML and process nodes
	var traverse func(*html.Node)
	traverse = func(n *html.Node) {
		// Check for <img> tags and process their src attributes
		if n.Type == html.ElementNode && n.Data == "img" {
			for i, attr := range n.Attr {
				if attr.Key == "src" && strings.HasPrefix(attr.Val, "cid:") {
					// Preserve cid attributes as-is
					n.Attr[i].Val = attr.Val
				}
			}
		}

		// Recursively process child nodes
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}
	}

	traverse(doc)

	// Render the sanitized HTML back to a string
	var buffer bytes.Buffer
	if err := html.Render(&buffer, doc); err != nil {
		return "", fmt.Errorf("failed to render sanitized HTML: %v", err)
	}

	return buffer.String(), nil
}

// Hàm parseEmail trả về body và attachments
func parseEmail(emailData string) (*ParsedEmail, error) {
	// Chuyển email thành một Reader
	reader := strings.NewReader(emailData)

	// Phân tích email bằng mail.ReadMessage (thuộc thư viện net/mail)
	msg, err := mail.ReadMessage(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read MIME message: %w", err)
	}

	// Trích xuất Subject từ header
	subject := msg.Header.Get("Subject")

	// Kiểm tra Content-Type để xác định kiểu email
	contentType := msg.Header.Get("Content-Type")
	// replyTo := msg.Header.Get("Reply-To")
	// messageId := msg.Header.Get("Message-ID")
	// fromHeader := msg.Header.Get("From-Header")
	mediaType, params, err := mime.ParseMediaType(contentType)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Content-Type: %w", err)
	}

	// Tạo ParsedEmail để lưu thông tin
	parsedEmail := &ParsedEmail{
		Subject: subject,
		// ReplyTo : replyTo,
		// MessageId: messageId,
		// FromHeader: fromHeader,
	}

	// Nếu email là multipart, xử lý các phần của nó
	if strings.HasPrefix(mediaType, "multipart/") {
		multipartReader := multipart.NewReader(msg.Body, params["boundary"])
		return parseMultipartEmail(multipartReader, parsedEmail)
	}

	// Nếu không phải multipart, đọc nội dung Body (text/plain hoặc text/html)
	bodyContent, err := io.ReadAll(msg.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading email body: %w", err)
	}

	if mediaType == "text/html" {
		parsedEmail.Body = string(bodyContent)
	}
	return parsedEmail, nil
}
func parseMultipartEmail(multipartReader *multipart.Reader, parsedEmail *ParsedEmail) (*ParsedEmail, error) {
	// Iterate through each part of the email
	for {
		part, err := multipartReader.NextPart()
		if err == io.EOF {
			break // No more parts
		}
		if err != nil {
			return nil, fmt.Errorf("error reading multipart part: %w", err)
		}

		// Get Content-Type and Content-Disposition
		contentType := part.Header.Get("Content-Type")
		contentDisposition := part.Header.Get("Content-Disposition")

		// Handle nested multipart/related
		if strings.HasPrefix(contentType, "multipart/related") {
			relatedReader := multipart.NewReader(part, getBoundary(contentType))
			_, err := parseMultipartEmail(relatedReader, parsedEmail)
			if err != nil {
				return nil, fmt.Errorf("error parsing multipart/related: %w", err)
			}
			continue
		}

		// Handle text parts (text/plain or text/html)
		if strings.HasPrefix(contentType, "text/") {
			bodyContent, err := io.ReadAll(part)
			if err != nil {
				return nil, fmt.Errorf("error reading body content: %w", err)
			}

			if strings.HasPrefix(contentType, "text/plain") || strings.HasPrefix(contentType, "text/html"){
				parsedEmail.Body = string(bodyContent)
				if strings.HasPrefix(contentType, "text/html"){
					extractBody, err := utils.ExtractBodyContent(string(bodyContent))
					if err != nil {
						return nil, fmt.Errorf("error reading attachment: %w", err)
					}
					parsedEmail.Discription = utils.GetMax200Characters(extractBody,maxDesLength)
				}else{
					parsedEmail.Discription = utils.GetMax200Characters(string(bodyContent),maxDesLength)
				}
			}
		}

		// Handle inline attachments (e.g., images in multipart/related)
		if strings.HasPrefix(contentDisposition, "inline") || strings.HasPrefix(contentDisposition, "attachment") {
			attachmentData, err := io.ReadAll(part)
			if err != nil {
				return nil, fmt.Errorf("error reading attachment: %w", err)
			}

			contentID := part.Header.Get("Content-ID")
			// Extract file details
			fileName := utils.ExtractFileName(contentDisposition)
			fileExt := utils.ExtractFileExtension(fileName)
			
			fileHash := utils.GenerateFileHash(attachmentData)

			contentLength := len(attachmentData)
			totalChunks := utils.CalculateChunks(contentLength, chunkSize) // Assuming 1 MB chunks

			parsedEmail.Attachments = append(parsedEmail.Attachments, Attachment{
				ContentDisposition: contentDisposition,
				ContentID:          contentID,
				ContentType:        contentType,
				Data:               attachmentData,
				FileName:           fileName,
				FileExtension:      fileExt,
				FileHash:           fileHash,
				ContentLength:      uint64(contentLength),
				TotalChunks:        uint64(totalChunks),
			})
		}
	}

	return parsedEmail, nil
}

// Helper function to extract boundary from Content-Type
func getBoundary(contentType string) string {
	parts := strings.Split(contentType, "boundary=")
	if len(parts) > 1 {
		return strings.Trim(parts[1], `"`)
	}
	return ""
}

func logAllFunctions(obj interface{}) {
	t := reflect.TypeOf(obj)

	// Kiểm tra nếu là con trỏ thì lấy giá trị thực tế của struct
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	// In ra thông tin của struct và các trường
	if t.Kind() == reflect.Struct {
		fmt.Printf("Struct Name: %s\n", t.Name())
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			fmt.Printf("Field %d: %s (%s)\n", i, field.Name, field.Type)
		}
	} else {
		fmt.Println("Không phải là struct.")
	}
}

// monitorFileAndRestart monitors the file for modification and restarts if modified
func monitorFileAndRestart() {

	for {
		// Get file information
		info, err := os.Stat(monitoredFile)
		if err != nil {
			fmt.Println("Failed to stat file %s: %v", monitoredFile, err)
		}

		// Check if the file modification time has changed
		if info.ModTime().After(lastModTime) {
			fmt.Println("File %s modified. Restarting application...", monitoredFile)
			lastModTime = info.ModTime()

			// Restart the application
			restartApp()
		}

		// Wait before checking again
		time.Sleep(pollingInterval)
	}
}

// restartApp stops the current process and starts a new instance of the executable
func restartApp() {
	// Start the new binary
	cmd := exec.Command("sudo", monitoredFile)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		fmt.Println("Failed to start application: %v", err)
		return
	}

	// Exit the current process
	os.Exit(0)
}

func isValidRecipientName(name string) bool {
	return len(name) > 0 && len(name) <= 42 && !strings.ContainsAny(name, " !@#$%^&*()")
}

func main() {

	// Initialize lastModTime to the current modification time of the file
	info, err := os.Stat(monitoredFile)
	if err != nil {
		fmt.Println("Failed to stat file %s: %v", monitoredFile, err)
		// return
	} else {
		lastModTime = info.ModTime()

		// Start monitoring the file
		go monitorFileAndRestart()
	}

	cconfig, err := config.LoadConfig("./config.yaml")
	if err != nil {
		log.Fatal("can not load config", err)
	}

	ChainClient, err = client.NewClient(
		&c_config.ClientConfig{
			Version_:                cconfig.MetaNodeVersion,
			PrivateKey_:             cconfig.PrivateKey_,
			ParentAddress:           cconfig.NodeAddress,
			ParentConnectionAddress: cconfig.NodeConnectionAddress,
			DnsLink_:                cconfig.DnsLink(),
		},
	)

	log.Println("Config ok")

	// create card abi
	reader, err := os.Open(cconfig.MailFactoryABIPath)
	if err != nil {
		log.Fatalf("Error occured while read baccarat abi")
	}
	defer reader.Close()

	mailFactoryAbi, err := abi.JSON(reader)
	if err != nil {
		log.Fatalf("Error occured while parse baccarat smart contract abi")
	}
	//
	readerMailStorage, err := os.Open(cconfig.MailStorageABIPath)
	if err != nil {
		log.Fatalf("Error occured while read baccarat abi")
	}
	defer readerMailStorage.Close()

	abiMailStorage, err := abi.JSON(readerMailStorage)
	if err != nil {
		log.Fatalf("Error occured while parse baccarat smart contract abi")
	}
	//
	readerFileStorage, err := os.Open(cconfig.FileABIPath)
	if err != nil {
		log.Fatalf("Error occured while read baccarat abi")
	}
	defer readerMailStorage.Close()

	abiFile, err := abi.JSON(readerFileStorage)
	if err != nil {
		log.Fatalf("Error occured while parse baccarat smart contract abi")
	}
	//
	servs := services.NewSendTransactionService(
		ChainClient,
		&mailFactoryAbi,
		common.HexToAddress(cconfig.MailFactoryAddress),
		&abiMailStorage,
		common.HexToAddress(cconfig.MailStorageAddress),
		common.HexToAddress(cconfig.NotiAddress),
		common.HexToAddress(cconfig.AdminAddress),
		&abiFile,
		common.HexToAddress(cconfig.FileAddress),
	)

	var emailStorageMap = sync.Map{}


	// Định nghĩa processor tùy chỉnh
	type myFooConfig struct {
		// SomeOption string `json:"some_option"` // Ví dụ về một cấu hình
	}

	MyFooProcessor := func() backends.Decorator {

		config := &myFooConfig{}

		// Hàm khởi tạo để thiết lập cấu hình cho MyFooProcessor
		initFunc := backends.InitializeWith(func(backendConfig backends.BackendConfig) error {
			// Trích xuất cấu hình từ backendConfig
			configType := backends.BaseConfig(&myFooConfig{})
			bcfg, err := backends.Svc.ExtractConfig(backendConfig, configType)
			if err != nil {
				return err
			}
			*config = *(bcfg.(*myFooConfig))
			return nil
		})

		// Đăng ký hàm khởi tạo cho MyFooProcessor
		backends.Svc.AddInitializer(initFunc)

		// log.Fatalf("Không thể khởi chạy server: 11")

		return func(p backends.Processor) backends.Processor {
			return backends.ProcessWith(func(e *guerrillaMail.Envelope, task backends.SelectTask) (backends.Result, error) {

				// log.Fatalf("Không thể khởi chạy server: hello")
				// log.Printf("---- Processing task: %v", task)

				if task == backends.TaskValidateRcpt {
					// Step 1: Extract the recipient email
					if len(e.RcptTo) == 0 {
						return backends.NewResult("550 No recipient provided"), nil
					}
					recipient := e.RcptTo[0].String()

					// Extract the portion before '@'
					recipientName := strings.Split(recipient, "@")[0]
					if recipientName == "" {
						log.Printf("Invalid recipient format: %s", recipient)
						return backends.NewResult("554 Invalid recipient email format"), nil
					}

					if !isValidRecipientName(recipientName) {
						log.Printf("Invalid recipient format: %s", recipient)
						return backends.NewResult("554 Invalid recipient email format"), nil
					}
					//get name from mns
					// Construct the URL with the domain name
					baseURL := cconfig.OwnerUrl + recipientName
					apiURL, err := url.Parse(baseURL)
					if err != nil {
						log.Printf("Can not parse url: %s", baseURL)
					}

					// Create the HTTP request
					resp, err := http.Get(apiURL.String())
					if err != nil {
						log.Printf("Can not get url: %s", baseURL)
					}
					defer resp.Body.Close()
					// Check the response status code
					if resp.StatusCode != http.StatusOK {
						log.Printf("Received non-OK HTTP status: %s", resp.Status)
						
					}

					// Read the response body
					body, err := ioutil.ReadAll(resp.Body)
					if err != nil {
						log.Printf("Fail in Read the response body")
						
					}
					// Convert the body to string
					bodyStr := string(body)
					// Trim any whitespace (including \n) from bodyStr
					bodyStr = strings.TrimSpace(bodyStr)
					fmt.Println("bodyStr:",bodyStr)
					//
					add := "0x"+bodyStr
					// Step 5: Fetch the smart contract address for the recipient
					// recipientName = "0xB50b908fFd42d2eDb12b325e75330c1AaAf35dc0"
					emailStorageAddress, err := servs.GetEmailStorage(add)
					if err != nil || emailStorageAddress == (common.Address{}) {
						log.Printf("Smart contract for recipient not found: %s", recipientName)
						return backends.NewResult("554 Recipient not found"), nil
					}

					// Step 6: Store the address for later use
					// emailStorageMap[recipientName] = emailStorageAddress
					emailStorageMap.Store(recipientName, emailStorageAddress)
					log.Printf("Smart contract for recipient %s: %s", recipientName, emailStorageAddress.(common.Address))

					return backends.NewResult("250 Recipient OK"), nil
				}

				if task == backends.TaskSaveMail {
					// Step 1: Extract the recipient email
					if len(e.RcptTo) == 0 {
						return backends.NewResult("550 No recipient provided"), nil
					}
					recipient := e.RcptTo[0].String()
					// Extract the portion before '@'
					recipientName := strings.Split(recipient, "@")[0]
					if recipientName == "" {
						log.Printf("Invalid recipient format: %s", recipient)
						return backends.NewResult("554 Invalid recipient email format"), nil
					}

					// Step 2: Retrieve the previously stored smart contract address
					emailStorageAddress, exists := emailStorageMap.Load(recipientName)
					if !exists {
						log.Printf("Smart contract for recipient not found during save: %s", recipientName)
						return backends.NewResult("554 Recipient validation not performed"), nil
					}

					ip := e.RemoteIP

					// return backends.NewResult("552 Error: Deny IP limit"), nil
					// return backends.NewResult(fmt.Sprintf("554 Deny IP limit: %s", ip)), nil

					sender := e.MailFrom.String()
					senderDomain := extractDomain(sender)

					// #debug chỗ này cần check thoả số lượng ký tự của ví và ENS
					// if !strings.EqualFold(recipientDomain, "example.com") {
					// 	// Trả về lỗi "Relay access denied"
					// 	errMessage := fmt.Sprintf("4.1.1 Error: Relay access denied hen : %s", recipientDomain)
					// 	log.Printf("DEBU: %s", errMessage)
					// 	return backends.NewResult(errMessage), nil
					// }

					// Kiểm tra kích thước email
					if len(e.Data.String()) > 1024*1024 {
						return backends.NewResult("552 Error: Message size exceeds 1MB limit"), nil
					}
					// Kiểm tra DKIM
					dkimResult, err := checkDKIM([]byte(e.Data.String()), senderDomain)
					if ip != "127.0.0.1" && !dkimResult {
						if err != nil {
					        log.Printf("DKIM error: %v", err)
					    }

				        log.Printf("DKIM failed, fallback to SPF and DMARC checks")

				        // Kiểm tra SPF
				        spfResult, spfErr := checkSPF(ip, senderDomain)
				        if spfErr != nil || !spfResult {
				            return backends.NewResult(fmt.Sprintf("554 SPF failed: %v", spfErr)), nil
				        }

				        // Kiểm tra DMARC
				        dmarcResult, dmarcErr := checkDMARC(senderDomain)
				        if dmarcErr != nil || !dmarcResult {
				            return backends.NewResult(fmt.Sprintf("554 DMARC failed: %v", dmarcErr)), nil
				        }
					}

					password, err := utils.GeneratePassword(recipient)

					// log.Printf("Ok content: %s", e.Data.String())
					log.Printf("Ok generating password: %s", password)

					if err != nil {
						log.Printf("Error generating password: %v", err)
						return backends.NewResult("554 Error generating password"), nil
					}

					encryptedEmail, err := encryptEmail(e.Data.String(), password)
					if err != nil {
						log.Printf("Error encrypting email: %v", err)
						return backends.NewResult("554 Error encrypting email"), nil
					}

					// Lưu email đã mã hoá xuống ổ cứng tại từng server
					err = utils.SaveEmailLocally(encryptedEmail)
					if err != nil {
						log.Printf("Error saving email locally: %v", err)
						return backends.NewResult("554 Error saving email locally"), nil
					}
					// Phân tích email
					parsedEmail, err := parseEmail(e.Data.String())
					if err != nil {
						log.Fatalf("Error parsing email: %v", err)
					}
					subject := parsedEmail.Subject
	                // html := sanitizeEmailHTML(parsedEmail.Html)
					body,err := sanitizeEmailHTML(parsedEmail.Body)
					if err != nil {
						log.Fatalf("Error sanitizeEmailHTML: %v", err)
					}
					encryptedBody, err := encryptEmail(body, password)
					if err != nil {
						log.Printf("Error encrypting email: %v", err)
						return backends.NewResult("554 Error encrypting email"), nil
					}
					// giai ma lai de dam bao noi dung
					decryptedBody, err := decryptEmail(encryptedBody, password)
					if err != nil {
						log.Printf("Error decrypting email: %v", err)
						return backends.NewResult("554 Error decrypting email"), nil
					}

					if (decryptedBody != body) {
						log.Printf("Error decrypting wrong email: %v", err)
						return backends.NewResult("554 Error decrypting wrong email"), nil
					}

					createdAt := uint64(time.Now().Unix())
					expireTime := createdAt + 30*24*60*60
					// Step 4: Convert attachments into the required format
					var infos []emailstorage.Info
					
					for _, attachment := range parsedEmail.Attachments {
						info := emailstorage.Info{
							Owner:common.HexToAddress(recipientName),
							Hash:attachment.FileHash,
							ContentLen:attachment.ContentLength,
							TotalChunks:attachment.TotalChunks,
							ExpireTime:expireTime ,
							Name:attachment.FileName,
							Ext:attachment.FileExtension,
							Status: 0,
							ContentDisposition: attachment.ContentDisposition,
							ContentID: attachment.ContentID,

						}
						infos = append(infos,info)
						
					}

					log.Printf("receiving email, subject: %s, body: %s, files: %d, createdAt %d, encryptedBody: %s", subject, body, createdAt, hex.EncodeToString(encryptedBody))
					// Step 8: Call the pushFileInfos method
					fileKeys,err := servs.PushFileInfos(infos)
					if err != nil {
						log.Printf("Failed to call pushFileInfos: %v", err)
						return backends.NewResult("554 Error calling pushFileInfos"), nil
					}
					fileKeysKq , ok := fileKeys.([][32]byte)
					if !ok {
						log.Printf("Failed to parse fileKeys: %v", err)
						return backends.NewResult("554 Error parse fileKeys"), nil
					}
					//Step 8: Call the uploadChunks method
					// Chunk the file data into smaller pieces
					 
					for index, attachment := range parsedEmail.Attachments {
						chunks, _ := utils.ChunkData(attachment.Data, chunkSize)
						// if err != nil {
						// 	log.Printf("Failed to parse UploadChunk: %v", err)
						// 	return backends.NewResult(fmt.Errorf("failed to chunk file data: %w", err)), nil
						// }

						// Initialize the lastChunkHash as an empty hash
						lastChunkHash := [32]byte{}
						var chunkHashes []string
						var chunkkq []string
						fmt.Println("len chunks la:",len(chunks))
						// Loop over each chunk and upload it to the contract
						for _, bchunk := range chunks {
							// bchunk ,err:= hex.DecodeString(chunk)
							// Upload each chunk with its hash
							lastChunkHash = utils.CalculateChunkHash(lastChunkHash, bchunk)
							kq,err := servs.UploadChunk(fileKeysKq[index], bchunk, lastChunkHash)
							if err != nil {
								log.Printf("Failed to call UploadChunk: %v", err)
								return backends.NewResult("554 Error calling UploadChunk"), nil
							}
							kqResult, ok := kq.(bool)
							if !ok || !kqResult {
								log.Printf("Failed to parse UploadChunk: %v", err)
								return backends.NewResult("554 Error calling UploadChunk"), nil
							}
							chunkHashes = append(chunkHashes, hex.EncodeToString(lastChunkHash[:]))
							chunkkq = append(chunkkq, hex.EncodeToString(bchunk))
						}
						// Write chunkHashes to a JSON file
						err = writeHashesToFile("chunk_hashes.json", chunkHashes)
						if err != nil {
							log.Printf("Failed to write chunk hashes to file: %v", err)
							return backends.NewResult("500 Internal Server Error"), nil
						}
						log.Println("File upload completed successfully.")
						// Write chunkHashes to a JSON file
						err = writeHashesToFile("chunkkq.json", chunkkq)
						if err != nil {
							log.Printf("Failed to write chunk hashes to file: %v", err)
							return backends.NewResult("500 Internal Server Error"), nil
						}
						log.Println("File upload completed successfully.")
					}
					// Step 8: Call the CreateEmail method
					discription := parsedEmail.Discription
					hash, err := servs.CreateEmail(emailStorageAddress.(common.Address), sender, subject, body, fileKeysKq, createdAt, discription)
					if err != nil {
						log.Printf("Failed to call CreateEmail: %v", err)
						return backends.NewResult("554 Error calling CreateEmail"), nil
					}

					// Step 9: Wait for the transaction to be mined
					log.Printf("Transaction submitted: %s", hash.(common.Hash))

					log.Printf("Email received and stored successfully: %+v", e)
					return backends.NewResult("250 OK: Email received and stored successfully"), nil
				}
				return p.Process(e, task)
			})
		}
	}

	// Cấu hình server SMTP (ví dụ)
	cfg := &guerrilla.AppConfig{
		LogFile:  "./go-guerrilla.log",
		LogLevel: "debug",
		// AllowedHosts: []string{}, // Không giới hạn domain, để kiểm tra logic tùy chỉnh
		AllowedHosts: []string{"m.pro", "payws.net", "payws.com"}, // Chỉ cho phép domain này
		Servers: []guerrilla.ServerConfig{
			{
				IsEnabled:       true,
				ListenInterface: "0.0.0.0:2025",
				MaxClients:      5,
				Timeout:         100,
			},
		},
		BackendConfig: backends.BackendConfig{
			"save_process":       "MyFooProcessor", // Sử dụng MyFooProcessor cho quá trình lưu email
			"validate_process":   "MyFooProcessor", // Sử dụng MyFooProcessor cho quá trình xác nhận người nhận
			"save_workers_size":  1,
			"log_received_mails": false,
		},
	}

	d := guerrilla.Daemon{Config: cfg}
	// Đặt backend cho server SMTP

	// Đăng ký MyFooProcessor vào backend
	d.AddProcessor("MyFooProcessor", MyFooProcessor)

	// Khởi chạy server SMTP
	if err := d.Start(); err != nil {
		log.Fatalf("Không thể khởi chạy SMTP server: %s", err)
	}

	select {}

}

// Hàm để lấy domain từ địa chỉ email
func extractDomain(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) == 2 {
		return parts[1]
	}
	return ""
}

// Hàm kiểm tra IP trong CIDR
func isIPInCIDR(ip, cidr string) bool {
	_, network, err := net.ParseCIDR(cidr)
	if err != nil {
		return false
	}
	return network.Contains(net.ParseIP(ip))
}

// Hàm kiểm tra SPF
func checkSPF(ip, domain string) (bool, error) {
	// Tạo yêu cầu DNS
	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(domain), dns.TypeTXT)

	// Gửi yêu cầu DNS đến Google Public DNS
	client := new(dns.Client)
	resp, _, err := client.Exchange(m, "8.8.8.8:53")
	if err != nil {
		return false, fmt.Errorf("DNS query failed: %v", err)
	}

	// Lấy bản ghi SPF từ câu trả lời DNS
	var spfRecord string
	for _, answer := range resp.Answer {
		if txt, ok := answer.(*dns.TXT); ok {
			for _, txtRecord := range txt.Txt {
				if strings.HasPrefix(txtRecord, "v=spf1") {
					spfRecord = txtRecord
					break
				}
			}
		}
	}

	if spfRecord == "" {
		return false, fmt.Errorf("no SPF record found for domain %s", domain)
	}

	// Phân tích và kiểm tra bản ghi SPF
	spfParts := strings.Split(spfRecord, " ")
	for _, part := range spfParts {
		if strings.HasPrefix(part, "ip4:") {
			allowedIP := strings.TrimPrefix(part, "ip4:")
			if strings.Contains(allowedIP, "/") {
				// Xử lý CIDR
				if isIPInCIDR(ip, allowedIP) {
					return true, nil
				}
			} else if ip == allowedIP {
				return true, nil
			}
		} else if strings.HasPrefix(part, "include:") {
			// Xử lý trường hợp "include"
			includedDomain := strings.TrimPrefix(part, "include:")
			result, err := checkSPF(ip, includedDomain) // Đệ quy kiểm tra include
			if err == nil && result {
				return true, nil
			}
		}
	}

	return false, fmt.Errorf("IP %s not authorized by SPF for domain %s", ip, domain)
}

func checkDMARC(domain string) (bool, error) {
	// client := &net.Resolver{}
	policy, err := dmarc.Lookup(domain)
	if err != nil {
		return false, fmt.Errorf("DMARC check failed: %v", err)
	}
	if policy.Policy == dmarc.PolicyReject {
		return false, fmt.Errorf("DMARC reject policy applied for domain: %s", domain)
	}
	return true, nil
}

func isUsingGmailMX(domain string) (bool, error) {
	// Query the MX records
	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(domain), dns.TypeMX)

	client := new(dns.Client)
	resp, _, err := client.Exchange(m, "8.8.8.8:53") // Use Google's public DNS
	if err != nil {
		return false, fmt.Errorf("DNS query failed: %v", err)
	}

	// Check for Google MX records
	for _, answer := range resp.Answer {
		if mx, ok := answer.(*dns.MX); ok {
			if strings.HasSuffix(mx.Mx, ".google.com.") {
				return true, nil
			}
		}
	}

	return false, nil
}

// Hàm kiểm tra DKIM
func checkDKIM(email []byte, senderDomain string) (bool, error) {
	result, err := dkim.Verify(&email)
	if err != nil {
		if strings.Contains(err.Error(), "signature has expired") {
			if senderDomain == "gmail.com" {
				return true, nil
			}

			usingMX, err := isUsingGmailMX(senderDomain)
			if err == nil && usingMX {
				return true, nil
			}

			return false, fmt.Errorf("DKIM signature expired")
		}
		return false, err
	}
	if result != 1 {
		return false, fmt.Errorf("DKIM verification failed")
	}
	return true, nil
}

// Hàm mã hoá email dạng gzip với password
// func encryptEmail(email, password string) ([]byte, error) {
// 	var buffer bytes.Buffer
// 	writer := gzip.NewWriter(&buffer)
// 	writer.Write([]byte(email))
// 	writer.Close()

// 	// Đơn giản hoá: kết hợp nội dung gzip với password (ở đây chỉ là minh hoạ)
// 	encrypted := append([]byte(password), buffer.Bytes()...)
// 	return encrypted, nil
// }

func encryptEmail(input, password string) ([]byte, error) {
	// Tạo buffer để lưu trữ dữ liệu nén gzip
	var buf bytes.Buffer
	gzipWriter := gzip.NewWriter(&buf)

	// Sao chép dữ liệu từ input vào gzip writer
	if _, err := io.WriteString(gzipWriter, input); err != nil {
		return nil, fmt.Errorf("failed to gzip data: %v", err)
	}
	gzipWriter.Close()

	// Mã hoá dữ liệu gzip bằng AES
	key := sha256.Sum256([]byte(password)) // Sử dụng SHA-256 để tạo khoá từ password
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher block: %v", err)
	}

	// Sử dụng Galois/Counter Mode (GCM) cho mã hoá
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %v", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	cipherText := gcm.Seal(nonce, nonce, buf.Bytes(), nil)

	return cipherText, nil
}
func decryptEmail(cipherText []byte, password string) (string, error) {
	// Tạo khóa từ password (SHA-256)
	key := sha256.Sum256([]byte(password))

	// Khởi tạo cipher block AES từ key
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return "", fmt.Errorf("failed to create cipher block: %v", err)
	}

	// Sử dụng GCM để giải mã
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %v", err)
	}

	// Lấy nonce từ cipherText (nonce được lưu trữ ở phần đầu của cipherText)
	nonceSize := gcm.NonceSize()
	if len(cipherText) < nonceSize {
		return "", fmt.Errorf("cipherText too short")
	}

	nonce, cipherText := cipherText[:nonceSize], cipherText[nonceSize:]

	// Giải mã dữ liệu
	plainText, err := gcm.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt data: %v", err)
	}

	// Giải nén gzip
	var buf bytes.Buffer
	buf.Write(plainText)
	gzipReader, err := gzip.NewReader(&buf)
	if err != nil {
		return "", fmt.Errorf("failed to create gzip reader: %v", err)
	}
	defer gzipReader.Close()

	// Đọc dữ liệu giải nén
	decompressedData, err := io.ReadAll(gzipReader)
	if err != nil {
		return "", fmt.Errorf("failed to read decompressed data: %v", err)
	}

	// Trả về chuỗi kết quả
	return string(decompressedData), nil
}

// Helper function to write hashes to a JSON file
func writeHashesToFile(filename string, hashes []string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // Pretty-print JSON
	return encoder.Encode(hashes)
}