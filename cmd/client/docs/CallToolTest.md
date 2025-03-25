# Hướng dẫn sử dụng công cụ `call_tool_test`

## Action `read`

- Thêm action `read` tương tự action `call`, gọi contract nhưng chỉ đọc dữ liệu.  Không thực hiện thay đổi trạng thái trên blockchain.


## Cập nhật loại tài khoản

Để cập nhật loại tài khoản, hãy chạy công cụ `cmd/client/call_tool_test` với lệnh sau:

```
go run main.go -data=UpdateAccountType.json
```

**File ví dụ: `UpdateAccountType.json`**

File `UpdateAccountType.json` cần có các trường sau:

* **`action`**:  Luôn có giá trị `"action": "update_account"`.
* **`input`**:  Chuỗi 8 ký tự hex (4 byte data) để chỉ định loại tài khoản.

| Giá trị `input` | Mô tả                                      |
|-----------------|----------------------------------------------|
| `"00000000"`     | Loại mặc định. Chỉ giao dịch bằng chữ ký BLS. |
| `"00000001"`     | Yêu cầu xác thực từ cả BLS và SECP (cho giao dịch từ Metamask hoặc Web3). |

