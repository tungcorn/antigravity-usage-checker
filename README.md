# Antigravity Usage Checker

🚀 Check your Antigravity AI usage quota from terminal

![Version](https://img.shields.io/github/v/release/tungcorn/antigravity-usage-checker)
![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)
![Windows](https://img.shields.io/badge/Windows-0078D6?style=flat&logo=windows&logoColor=white)
![macOS](https://img.shields.io/badge/macOS-000000?style=flat&logo=apple&logoColor=white)
![Linux](https://img.shields.io/badge/Linux-FCC624?style=flat&logo=linux&logoColor=black)
[![CI](https://github.com/tungcorn/antigravity-usage-checker/actions/workflows/ci.yml/badge.svg)](https://github.com/tungcorn/antigravity-usage-checker/actions/workflows/ci.yml)
[![CodeQL](https://github.com/tungcorn/antigravity-usage-checker/actions/workflows/codeql.yml/badge.svg)](https://github.com/tungcorn/antigravity-usage-checker/actions/workflows/codeql.yml)

![Screenshot](assets/image.png)

## 🔒 Security

This tool is designed to be transparent and safe:

- ✅ **Local only** - Only connects to `localhost` (Antigravity local server)
- ✅ **No network calls** - Does NOT send any data to external servers
- ✅ **No telemetry** - No tracking, analytics, or data collection
- ✅ **Open source** - All code is public and auditable
- ✅ **CodeQL scanned** - Automatically scanned for security vulnerabilities

> 💡 You can review the [install scripts](install.ps1) before running them.

## 🔧 How It Works

1. Detects the running Windsurf language server process
2. Extracts connection parameters (port, CSRF token) from process arguments
3. Calls the local API at `127.0.0.1` to fetch quota data
4. Parses and displays the information in terminal

> **Note**: This tool only reads publicly available process information and communicates with localhost. No external network requests are made.

---

## 🇬🇧 English

### Quick Install

**Windows (PowerShell):**
```powershell
iwr https://raw.githubusercontent.com/tungcorn/antigravity-usage-checker/main/install.ps1 -OutFile $env:TEMP\install.ps1; . $env:TEMP\install.ps1
```

**macOS / Linux (Bash):**
```bash
curl -fsSL https://raw.githubusercontent.com/tungcorn/antigravity-usage-checker/main/install.sh | bash
```

Then run:
```bash
agusage
```

### Update

Run the install command again to update to the latest version.

### Install Specific Version

**Windows:**
```powershell
iwr https://raw.githubusercontent.com/tungcorn/antigravity-usage-checker/main/install.ps1 -OutFile $env:TEMP\install.ps1; . $env:TEMP\install.ps1 -Version 0.5.0
```

**macOS / Linux:**
```bash
curl -fsSL https://raw.githubusercontent.com/tungcorn/antigravity-usage-checker/main/install.sh | bash -s -- -v 0.5.0
```

### Manual Install

1. **Download** from [Releases](https://github.com/tungcorn/antigravity-usage-checker/releases/latest)
   - Windows: `antigravity-usage-checker-windows-amd64.zip`
   - macOS Intel: `antigravity-usage-checker-darwin-amd64.tar.gz`
   - macOS Apple Silicon: `antigravity-usage-checker-darwin-arm64.tar.gz`
   - Linux: `antigravity-usage-checker-linux-amd64.tar.gz`

2. **Extract** the archive to a folder of your choice

3. **Run** the executable:
   - Windows: Double-click `agusage.exe` or run from terminal
   - macOS/Linux: Run `chmod +x agusage` first, then `./agusage`

4. **(Optional) Add to PATH** for global access:
   - Windows: Move `agusage.exe` to a folder in your PATH
   - macOS/Linux: Move to `/usr/local/bin/` or add the folder to `$PATH`

> ⚠️ Antigravity must be running

### Features ✨

- 🎨 **Color-coded display** - Green when quota is healthy (>50%), yellow when moderate (>20%), red when low
- 📊 **Smart total calculation** - Automatically detects and deduplicates shared quota pools
- 🔄 **Unicode progress bars** - Beautiful █ and ░ characters for visual progress
- ⚡ **Fast and lightweight** - Written in Go, single binary, no dependencies
- 💾 **Offline cache** - Works even when Antigravity is not running

### Commands

```bash
agusage          # Show quota
agusage --json   # JSON output
agusage --help   # Help
```

### Platform Support

| Platform | Status |
|----------|--------|
| Windows | ✅ Fully tested |
| macOS | ✅ CI tested |
| Linux | ✅ CI tested |

---

## 🇻🇳 Tiếng Việt

### Cài đặt nhanh

**Windows (PowerShell):**
```powershell
iwr https://raw.githubusercontent.com/tungcorn/antigravity-usage-checker/main/install.ps1 -OutFile $env:TEMP\install.ps1; . $env:TEMP\install.ps1
```

**macOS / Linux (Bash):**
```bash
curl -fsSL https://raw.githubusercontent.com/tungcorn/antigravity-usage-checker/main/install.sh | bash
```

Sau đó chạy:
```bash
agusage
```

### Cập nhật

Chạy lại lệnh cài đặt để cập nhật lên phiên bản mới nhất.

### Cài đặt phiên bản cụ thể

**Windows:**
```powershell
iwr https://raw.githubusercontent.com/tungcorn/antigravity-usage-checker/main/install.ps1 -OutFile $env:TEMP\install.ps1; . $env:TEMP\install.ps1 -Version 0.5.0
```

**macOS / Linux:**
```bash
curl -fsSL https://raw.githubusercontent.com/tungcorn/antigravity-usage-checker/main/install.sh | bash -s -- -v 0.5.0
```

### Cài thủ công

1. **Tải** từ [Releases](https://github.com/tungcorn/antigravity-usage-checker/releases/latest)
   - Windows: `antigravity-usage-checker-windows-amd64.zip`
   - macOS Intel: `antigravity-usage-checker-darwin-amd64.tar.gz`
   - macOS Apple Silicon: `antigravity-usage-checker-darwin-arm64.tar.gz`
   - Linux: `antigravity-usage-checker-linux-amd64.tar.gz`

2. **Giải nén** file vào thư mục bạn chọn

3. **Chạy** chương trình:
   - Windows: Double-click `agusage.exe` hoặc chạy từ terminal
   - macOS/Linux: Chạy `chmod +x agusage` trước, sau đó `./agusage`

4. **(Tùy chọn) Thêm vào PATH** để chạy từ mọi nơi:
   - Windows: Di chuyển `agusage.exe` vào thư mục trong PATH
   - macOS/Linux: Di chuyển vào `/usr/local/bin/` hoặc thêm thư mục vào `$PATH`

> ⚠️ Antigravity phải đang chạy

### Tính năng ✨

- 🎨 **Màu sắc thông minh** - Xanh lá khi quota còn nhiều (>50%), vàng khi trung bình (>20%), đỏ khi sắp hết
- 📊 **Tính tổng thông minh** - Tự động phát hiện và loại bỏ trùng lặp các quota pools dùng chung
- 🔄 **Progress bar Unicode** - Ký tự █ và ░ đẹp mắt cho thanh tiến độ
- ⚡ **Nhanh và nhẹ** - Viết bằng Go, binary đơn giản, không cần dependencies
- 💾 **Cache offline** - Hoạt động ngay cả khi Antigravity không chạy

### Các lệnh

```bash
agusage          # Xem quota
agusage --json   # Xuất JSON
agusage --help   # Trợ giúp
```

### Hỗ trợ nền tảng

| Nền tảng | Trạng thái |
|----------|------------|
| Windows | ✅ Đã test đầy đủ |
| macOS | ✅ Đã test CI |
| Linux | ✅ Đã test CI |

---

## Development

### Run tests
```bash
go test ./... -v
```

### Build
```bash
go build -o agusage ./cmd/agusage/
```

---

## License

MIT © 2024-present

---

<p align="center">
  <b>If you find this useful, give it a ⭐!</b>
</p>
