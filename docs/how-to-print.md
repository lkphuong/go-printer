## In Trên macOS/Linux

Trên macOS và Linux, hệ thống in thường dựa trên CUPS (Common Unix Printing System). Bạn có thể sử dụng lệnh `lpstat` để kiểm tra trạng thái máy in và các công việc in đang chờ.

### Bước 1: Tìm Máy In

Mở Terminal và chạy lệnh sau để xem danh sách máy in có sẵn:

```
lpstat -p
```

Lệnh này sẽ hiển thị các máy in đã cài đặt và trạng thái của chúng (ví dụ: ipos, canon).

### Bước 2: In

Để in một file (ví dụ: PDF hoặc văn bản), sử dụng lệnh `lp`:

```
lp filename.pdf
```

- Thay `filename.pdf` bằng tên file bạn muốn in.
- Bạn có thể chỉ định máy in cụ thể bằng tùy chọn `-d`:
  ```
  lp -d printer_name filename.pdf
  ```

### Bước 3: Kiểm Tra Queue

Sau khi gửi in, kiểm tra trạng thái:

```
lpstat -o
```

Lệnh này liệt kê các công việc in đang chờ hoặc đang xử lý.

### Note

- Đảm bảo máy in đã được cài đặt và kết nối đúng cách qua System Preferences (macOS) hoặc CUPS web interface (Linux).
- Nếu gặp vấn đề, kiểm tra log của CUPS hoặc khởi động lại dịch vụ in.

## In Trên Windows

Trên Windows, việc in thường yêu cầu qua một ứng dụng trung gian vì hệ thống không hỗ trợ gửi trực tiếp file đến máy in từ dòng lệnh như trên Unix-based systems.

### Bước 1: Sử Dụng Paint (MSPaint)

Windows có sẵn ứng dụng Paint (mspaint.exe), có thể dùng để in file hình ảnh hoặc văn bản.

- Mở Command Prompt hoặc PowerShell.
- Chạy lệnh:
  ```
  mspaint /pt "C:\path\to\filename.jpg"
  ```
  - Thay `C:\path\to\filename.jpg` bằng đường dẫn đầy đủ đến file bạn muốn in.
  - Lệnh này sẽ mở file trong Paint và tự động in nó ra máy in mặc định.

### Tại Sao Dùng Paint?

- Paint có sẵn trên mọi phiên bản Windows, không cần cài đặt thêm.
- Bạn cũng có thể dùng Microsoft Edge hoặc các ứng dụng khác như Adobe Reader cho PDF, nhưng Paint là lựa chọn đơn giản nhất cho hình ảnh.

### Bước 2: Kiểm Tra Máy In

- Mở Settings > Devices > Printers & scanners để xem danh sách máy in.
- Nếu cần in qua máy in cụ thể, mở file trong ứng dụng tương ứng và chọn máy in từ dialog in.

### Note

- Đảm bảo máy in đã được cài đặt driver và kết nối.
- Nếu file không phải hình ảnh, chuyển đổi sang định dạng hỗ trợ (ví dụ: PDF sang JPG) trước khi in.
