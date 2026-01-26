# Contributing

🌐 [🇻🇳 Tiếng Việt](#vietnamese) | [🇬🇧 English](#english)

<a id="english"></a>
## 🇬🇧 English

Thank you for your interest in contributing to Antigravity Usage Checker!

### Prerequisites

- [Go](https://go.dev/dl/) 1.25.5 or higher

### Setup

1. Clone the repository:
```bash
git clone https://github.com/tungcorn/antigravity-usage-checker.git
cd antigravity-usage-checker
```

2. Download dependencies:
```bash
go mod download
```

### Run without Build

You can run the tool directly from source:
```bash
go run ./cmd/agusage/
```
Or with specific arguments:
```bash
go run ./cmd/agusage/ --json
```

### Testing

Run all unit tests:
```bash
go test ./... -v
```

### Build

#### Current Platform
```bash
go build -o agusage ./cmd/agusage/
```

#### Cross-platform Build
Generate binaries for different operating systems:

**Windows (64-bit):**
```bash
GOOS=windows GOARCH=amd64 go build -o agusage.exe ./cmd/agusage/
```

**macOS (Intel):**
```bash
GOOS=darwin GOARCH=amd64 go build -o agusage ./cmd/agusage/
```

**macOS (Apple Silicon):**
```bash
GOOS=darwin GOARCH=arm64 go build -o agusage ./cmd/agusage/
```

**Linux (64-bit):**
```bash
GOOS=linux GOARCH=amd64 go build -o agusage ./cmd/agusage/
```

### Project Structure

- `cmd/agusage/`: Entry point of the application.
- `internal/`: Private library code.
  - `auth/`: Authentication and CSRF token handling.
  - `client/`: Local API client for Antigravity.
  - `discovery/`: Windsurf process discovery logic.
  - `display/`: Terminal output formatting and progress bars.
  - `quota/`: Data structures and quota logic.

---
<a id="vietnamese"></a>
## 🇻🇳 Tiếng Việt

Cảm ơn bạn đã quan tâm đóng góp cho dự án Antigravity Usage Checker!

### Yêu cầu tiên quyết (Prerequisites)

- [Go](https://go.dev/dl/) 1.25.5 hoặc mới hơn

### Cài đặt (Setup)

1. Clone repository về máy:
```bash
git clone https://github.com/tungcorn/antigravity-usage-checker.git
cd antigravity-usage-checker
```

2. Tải các thư viện phụ thuộc (dependencies):
```bash
go mod download
```

### Chạy trực tiếp (Run without Build)

Bạn có thể chạy công cụ trực tiếp từ mã nguồn mà không cần build:
```bash
go run ./cmd/agusage/
```
Hoặc chạy với các tham số cụ thể (ví dụ xuất JSON):
```bash
go run ./cmd/agusage/ --json
```

### Kiểm thử (Testing)

Chạy tất cả các unit test:
```bash
go test ./... -v
```

### Đóng gói (Build)

#### Nền tảng hiện tại (Current Platform)
```bash
go build -o agusage ./cmd/agusage/
```

#### Đóng gói đa nền tảng (Cross-platform Build)
Tạo file chạy cho các hệ điều hành khác nhau:

**Windows (64-bit):**
```bash
GOOS=windows GOARCH=amd64 go build -o agusage.exe ./cmd/agusage/
```

**macOS (Intel):**
```bash
GOOS=darwin GOARCH=amd64 go build -o agusage ./cmd/agusage/
```

**macOS (Apple Silicon):**
```bash
GOOS=darwin GOARCH=arm64 go build -o agusage ./cmd/agusage/
```

**Linux (64-bit):**
```bash
GOOS=linux GOARCH=amd64 go build -o agusage ./cmd/agusage/
```

### Cấu trúc dự án (Project Structure)

- `cmd/agusage/`: Điểm bắt đầu của ứng dụng (hàm main).
- `internal/`: Mã nguồn thư viện nội bộ.
  - `auth/`: Xử lý xác thực và CSRF token.
  - `client/`: Client gọi Local API của Antigravity.
  - `discovery/`: Logic tìm kiếm tiến trình Windsurf.
  - `display/`: Định dạng hiển thị terminal và thanh tiến độ.
  - `quota/`: Cấu trúc dữ liệu và logic xử lý quota.
