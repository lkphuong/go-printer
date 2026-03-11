# Windows Service - Hướng dẫn chi tiết

## 📌 Tổng quan

Project này hỗ trợ **3 cách** tạo Windows Service:

1. **Go Program** - Dùng PowerShell `New-Service` command (từ Go code)
2. **PowerShell Script** - Direct script
3. **Batch File** - Wrapper để dễ chạy

Tất cả đều sử dụng **PowerShell `New-Service` command** thay vì third-party libraries.

---

## 🚀 Cách 1: Go Program (Khuyên dùng nhất)

### Build

```bash
make build-windows-service
```

Hoặc:

```bash
GOOS=windows GOARCH=amd64 go build -o build/printer-service.exe ./cmd/service/main.go
```

### Chạy

**Double-click quyền Admin:**

- Vào folder `build/`
- Right-click `printer-service.exe`
- Chọn **Run as administrator**

**Hoặc PowerShell:**

```powershell
Start-Process -FilePath "C:\path\to\build\printer-service.exe" -Verb RunAs
```

### Kết quả

```
╔════════════════════════════════════════╗
║   Go Printer - Windows Service         ║
╚════════════════════════════════════════╝

Service not found. Creating...
Executable: C:\path\to\build\printer-service.exe
Working directory: C:\path\to\build

Service created and started successfully

╔════════════════════════════════════════╗
║ ✅ Setup Complete!                    ║
╚════════════════════════════════════════╝

Access API at: http://localhost:9099
Logs at: C:\path\to\logs\app.log

The window will close in 5 seconds...
```

### Command line options

```powershell
# Interactive setup
.\build\printer-service.exe

# Create service
.\build\printer-service.exe -install

# Start service
.\build\printer-service.exe -start

# Stop service
.\build\printer-service.exe -stop

# Remove service
.\build\printer-service.exe -remove

# Check status
.\build\printer-service.exe -status

# Show help
.\build\printer-service.exe -help
```

---

## 🔧 Cách 2: PowerShell Script

### Chạy

```powershell
powershell -ExecutionPolicy Bypass -File .\scripts\install-service.ps1
```

**Hoặc double-click:**

- Vào folder `scripts/`
- Right-click `install-service.ps1`
- Chọn **Run with PowerShell**

### Menu tương tác

```
=================================
Go Printer - Windows Service Setup
=================================

✓ Tìm thấy executable: C:\path\to\build\printer-service.exe

Chọn hành động:
1. Cài đặt service (install)
2. Khởi động service (start)
3. Dừng service (stop)
4. Gỡ cài đặt service (uninstall)
5. Kiểm tra status

Nhập lựa chọn (1-5): _
```

---

## 📦 Cách 3: Batch File

### Chạy

Double-click: `scripts/install-service.bat`

**Hoặc:**

```cmd
scripts\install-service.bat
```

### Tự động:

- Request quyền Administrator (nếu chưa có)
- Gọi PowerShell script
- Đóng cửa sổ sau 3 giây

---

## 🎯 Tính năng

### Trong Go Program (`cmd/service/main.go`)

```go
// Sử dụng os/exec để chạy PowerShell
cmd := exec.Command("powershell", "-NoProfile", "-Command", script)
output, err := cmd.CombinedOutput()

// PowerShell script tạo service
script := `
    New-Service `
        -Name "PrinterService" `
        -BinaryPathName "C:\path\to\printer-service.exe" `
        -DisplayName "Go Printer Service" `
        -Description "Windows Service cho API Printer" `
        -StartupType Automatic
`
```

### Tính năng

✅ **Auto-check quyền Administrator**

- Nếu không có, sẽ request tự động
- Go program: hiển thị thông báo lỗi
- Batch file: tự động elevate quyền

✅ **Detect service tồn tại**

- Nếu đã tồn tại: cho menu chọn (start/stop/reinstall/remove)
- Nếu chưa: tạo mới

✅ **Tự động khởi động**

- Service sẽ start với `Automatic` startup type

✅ **Tự đóng cửa sổ**

- Go program: 5 giây
- Batch file: 3 giây
- PowerShell: click Enter

✅ **Logs chi tiết**

- Service logs: `logs/app.log`
- PowerShell output rõ ràng
- Color coded output

---

## 📊 So sánh 3 cách

| Tiêu chí      | Go Program | PowerShell | Batch            |
| ------------- | ---------- | ---------- | ---------------- |
| **Dễ dùng**   | ⭐⭐⭐⭐⭐ | ⭐⭐⭐     | ⭐⭐⭐⭐         |
| **Tính năng** | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐   | ⭐⭐             |
| **Yêu cầu**   | Chỉ .exe   | PowerShell | CMD + PowerShell |
| **Menu**      | ⭐⭐⭐⭐   | ⭐⭐⭐⭐   | ❌               |
| **Tự động**   | ⭐⭐⭐⭐⭐ | ⭐⭐⭐     | ⭐⭐⭐⭐⭐       |

---

## 🛑 Quản lý Service

### PowerShell

```powershell
# Xem tất cả service
Get-Service

# Xem chi tiết service
Get-Service -Name "PrinterService" | Format-List

# Start
Start-Service -Name "PrinterService"

# Stop
Stop-Service -Name "PrinterService" -Force

# Restart
Restart-Service -Name "PrinterService" -Force

# Gỡ cài đặt
# ⚠️ Phải stop trước!
Stop-Service -Name "PrinterService" -Force
Remove-Service -Name "PrinterService" -Force

# Thay đổi startup type
Set-Service -Name "PrinterService" -StartupType Automatic
Set-Service -Name "PrinterService" -StartupType Manual
Set-Service -Name "PrinterService" -StartupType Disabled
```

### Services GUI

1. Mở `services.msc` (Win + R → gõ `services.msc`)
2. Tìm **"PrinterService"**
3. Right-click → Properties
4. Quản lý startup type, description, etc.

### Command Prompt

```cmd
REM Start
net start PrinterService

REM Stop
net stop PrinterService

REM Check all services
net start
```

---

## 🐛 Troubleshooting

### Port 9099 bị chiếm

```powershell
# Xem ai dùng port 9099
netstat -ano | findstr :9099

# Kill process
taskkill /PID <PID> /F
```

### Service không start

1. Kiểm tra Event Viewer:
   - Mở `Event Viewer`
   - Windows Logs → Application
   - Tìm entries liên quan đến PrinterService

2. Xem service logs:

   ```bash
   type logs\app.log
   tail -f logs\app.log  # PowerShell
   ```

3. Test executable:
   ```powershell
   .\build\printer-service.exe
   ```

### Permission denied

- Chạy PowerShell/CMD dưới quyền Administrator
- Right-click → **Run as administrator**

### Service tồn tại nhưng không hoạt động

```powershell
# Remove cũ
Stop-Service -Name "PrinterService" -Force
Remove-Service -Name "PrinterService" -Force

# Reinstall
.\build\printer-service.exe -remove
.\build\printer-service.exe -install
```

---

## 📁 Cấu trúc File

```
go-printer/
├── cmd/
│   └── service/
│       └── main.go              ← Go program (sử dụng PowerShell)
├── scripts/
│   ├── install-service.ps1      ← PowerShell script (trực tiếp New-Service)
│   └── install-service.bat      ← Batch wrapper (gọi .ps1)
├── build/
│   └── printer-service.exe       ← Built executable
├── logs/
│   └── app.log                  ← Service logs
├── Makefile                     ← Có target build-windows-service
└── ServiceSetup.md              ← Quick start
```

---

## 💡 Development

### Cái gì được sử dụng?

**Go Program (`cmd/service/main.go`):**

```go
import (
    "os/exec"  // Chạy PowerShell commands
)

// Tạo service
cmd := exec.Command("powershell", "-NoProfile", "-Command", `
    New-Service ...
`)
```

**PowerShell Script (`scripts/install-service.ps1`):**

```powershell
# Direct New-Service command
New-Service `
    -Name "PrinterService" `
    -BinaryPathName $exePath `
    -DisplayName "Go Printer Service" `
    -StartupType Automatic
```

### Không sử dụng

❌ `kardianos/service` - Thay bằng `os/exec` + PowerShell
❌ NSSM - Không cần
❌ WinSW - Không cần

---

## ⚡ Quick Reference

```bash
# Build
make build-windows-service

# Install & run (double-click quyền Admin)
.\build\printer-service.exe

# Or manual command
.\build\printer-service.exe -install
Start-Service -Name PrinterService

# Check
Get-Service PrinterService

# Stop
Stop-Service -Name PrinterService -Force

# Remove
Remove-Service -Name PrinterService -Force
```

---

## 📝 Log Output

**Success:**

```
Go Printer - Windows Service
================================

✅ Setup Complete!
Access API at: http://localhost:9099
Logs at: C:\path\to\logs\app.log
```

**Logs file:**

```
[2026-03-11 10:30:45] Server starting on :9099
[2026-03-11 10:30:46] Initializing folders...
[2026-03-11 10:30:47] Folders and logger initialized.
```

---

## ✅ Checklist

- [ ] Build: `make build-windows-service`
- [ ] Run dưới Admin quyền
- [ ] Service tạo thành công
- [ ] Access http://localhost:9099
- [ ] Kiểm tra logs: `logs/app.log`
- [ ] Set auto-start (nếu cần)
