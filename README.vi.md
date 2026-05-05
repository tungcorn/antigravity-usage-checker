# Antigravity Usage Checker

[English](README.md)

![Version](https://img.shields.io/github/v/release/tungcorn/antigravity-usage-checker)
![Go](https://img.shields.io/badge/Go-1.25.5+-00ADD8?logo=go)
![Windows](https://img.shields.io/badge/Windows-0078D6?style=flat&logo=windows&logoColor=white)
![macOS](https://img.shields.io/badge/macOS-000000?style=flat&logo=apple&logoColor=white)
![Linux](https://img.shields.io/badge/Linux-FCC624?style=flat&logo=linux&logoColor=black)
[![CI](https://github.com/tungcorn/antigravity-usage-checker/actions/workflows/ci.yml/badge.svg)](https://github.com/tungcorn/antigravity-usage-checker/actions/workflows/ci.yml)
[![CodeQL](https://github.com/tungcorn/antigravity-usage-checker/actions/workflows/codeql.yml/badge.svg)](https://github.com/tungcorn/antigravity-usage-checker/actions/workflows/codeql.yml)

![Screenshot](assets/image.png)

Kiểm tra quota sử dụng Antigravity AI ngay trong terminal.

Antigravity Usage Checker là một công cụ dòng lệnh nhỏ viết bằng Go. Công cụ đọc thông tin quota từ API local của Antigravity/Windsurf và hiển thị kết quả theo định dạng dễ đọc trong terminal.

## Cảnh báo bảo mật quan trọng

Hiện có một fork độc hại của repository này: `Motasem-amer/antigravity-usage-checker`.

Fork đó từng phân phối malware giả dạng công cụ này bằng cách thay README bằng các liên kết "Download Now" giả, trỏ tới file ZIP chứa malware dropper cho Windows.

Chỉ tải công cụ này từ:

- **Repository chính**: [tungcorn/antigravity-usage-checker](https://github.com/tungcorn/antigravity-usage-checker)
- **Bản phát hành chính thức**: [GitHub Releases](https://github.com/tungcorn/antigravity-usage-checker/releases)

Không tải công cụ này từ fork hoặc nguồn không chính thức. Nếu bạn đã tải từ fork nêu trên, hãy quét hệ thống ngay lập tức.

## Mô hình bảo mật

Công cụ này được thiết kế để minh bạch và an toàn:

- **Chỉ dùng API local**: Công cụ kết nối tới server local của Antigravity tại `127.0.0.1`.
- **Không gọi mạng bên ngoài khi chạy**: Ứng dụng không gửi dữ liệu sử dụng tới server bên ngoài.
- **Không telemetry**: Không tracking, analytics hoặc thu thập dữ liệu.
- **Mã nguồn mở**: Code và script cài đặt đều có thể được kiểm tra.
- **Chỉ đọc dữ liệu**: Công cụ đọc thông tin tiến trình local và dữ liệu quota; không chỉnh sửa cấu hình Antigravity.

Các script cài đặt sẽ tải file release từ GitHub. Bạn có thể xem `install.ps1` hoặc `install.sh` trước khi chạy nếu muốn.

## Yêu cầu

- Antigravity cần đang chạy để lấy dữ liệu quota mới nhất.
- Chỉ cần kết nối internet khi dùng installer hoặc tải file release.

## Cài đặt nhanh

### Windows

Chạy trong PowerShell:

```powershell
iwr https://raw.githubusercontent.com/tungcorn/antigravity-usage-checker/main/install.ps1 -OutFile $env:TEMP\install.ps1; . $env:TEMP\install.ps1
```

Nếu PowerShell chặn việc chạy script, dùng lệnh sau:

```powershell
powershell -ExecutionPolicy Bypass -Command "iwr https://raw.githubusercontent.com/tungcorn/antigravity-usage-checker/main/install.ps1 -OutFile $env:TEMP\install.ps1; & $env:TEMP\install.ps1"
```

### macOS và Linux

Chạy trong shell:

```bash
curl -fsSL https://raw.githubusercontent.com/tungcorn/antigravity-usage-checker/main/install.sh | bash
```

Sau đó chạy:

```bash
agusage
```

## Cập nhật

Chạy lại lệnh cài đặt để cập nhật lên phiên bản mới nhất.

## Cài đặt phiên bản cụ thể

### Windows

```powershell
powershell -ExecutionPolicy Bypass -Command "iwr https://raw.githubusercontent.com/tungcorn/antigravity-usage-checker/main/install.ps1 -OutFile $env:TEMP\install.ps1; & $env:TEMP\install.ps1 -Version 0.5.0"
```

### macOS và Linux

```bash
curl -fsSL https://raw.githubusercontent.com/tungcorn/antigravity-usage-checker/main/install.sh | bash -s -- -v 0.5.0
```

## Cài đặt thủ công

1. Tải bản mới nhất từ [GitHub Releases](https://github.com/tungcorn/antigravity-usage-checker/releases/latest).

2. Chọn file phù hợp với nền tảng của bạn:

   - Windows: `antigravity-usage-checker-windows-amd64.zip`
   - macOS Intel: `antigravity-usage-checker-darwin-amd64.tar.gz`
   - macOS Apple Silicon: `antigravity-usage-checker-darwin-arm64.tar.gz`
   - Linux: `antigravity-usage-checker-linux-amd64.tar.gz`

3. Giải nén file.

4. Chạy chương trình:

   - Windows: chạy `agusage.exe`
   - macOS/Linux: chạy `chmod +x agusage`, sau đó `./agusage`

5. Tùy chọn: thêm file chạy vào `PATH`.

   Windows:

   1. Di chuyển `agusage.exe` vào một thư mục, ví dụ `C:\Tools`.
   2. Nhấn `Win + R`, gõ `sysdm.cpl`, rồi nhấn Enter.
   3. Mở **Advanced** > **Environment Variables**.
   4. Trong **User variables**, chọn `Path` > **Edit** > **New**.
   5. Thêm đường dẫn thư mục và khởi động lại terminal.

   macOS/Linux:

   ```bash
   sudo mv agusage /usr/local/bin/
   ```

   Hoặc thêm một thư mục tùy chọn vào `PATH` trong `~/.bashrc` hoặc `~/.zshrc`:

   ```bash
   export PATH="$PATH:/path/to/your/folder"
   ```

## Cách sử dụng

```bash
agusage            # Xem quota hiện tại
agu                # Alias ngắn được tạo bởi script cài đặt
agusage --json     # Xuất JSON
agusage -j         # Xuất JSON bằng tùy chọn ngắn
agusage --version  # Xem thông tin phiên bản
agusage --help     # Xem trợ giúp
```

## Tính năng

- **Hiển thị theo màu**: Dùng màu xanh, vàng và đỏ theo lượng quota còn lại.
- **Tính tổng thông minh**: Phát hiện và loại bỏ trùng lặp các quota pool dùng chung.
- **Thanh tiến độ Unicode**: Hiển thị quota bằng thanh tiến độ trong terminal.
- **Nhanh và nhẹ**: Viết bằng Go và phân phối dưới dạng một binary duy nhất.
- **Cache offline**: Có thể hiển thị dữ liệu quota đã cache khi Antigravity không chạy.

## Cách hoạt động

1. Tìm tiến trình Windsurf language server đang chạy.
2. Đọc thông tin kết nối local từ tham số của tiến trình.
3. Gọi API local của Antigravity tại `127.0.0.1`.
4. Phân tích dữ liệu quota và hiển thị trong terminal.

## Hạn chế đã biết

API local của Antigravity thường chỉ cập nhật thống kê sử dụng theo các mốc như 0%, 20%, 40%, 60% và 80%.

Công cụ này chủ động dùng API local an toàn, chỉ đọc dữ liệu. Một số công cụ khác có thể hiển thị phần trăm chi tiết hơn bằng cách hook vào event stream nội bộ của IDE hoặc bộ đếm token, đó là một cách tiếp cận khác.

## Hỗ trợ nền tảng

| Nền tảng | Trạng thái |
| --- | --- |
| Windows | Đã test đầy đủ |
| macOS | Đã test qua CI |
| Linux | Đã test qua CI |

## Đóng góp

Để xem hướng dẫn thiết lập môi trường phát triển, chạy test và build, đọc [hướng dẫn đóng góp](CONTRIBUTING.md).

## Giấy phép

[MIT](LICENSE) © 2024-present
