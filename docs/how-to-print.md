# Hướng Dẫn Cài Đặt Dịch Vụ Máy In

## Tổng Quan

Dịch vụ máy in này là một ứng dụng Go cung cấp API REST để quản lý máy in, cấu hình in ấn và gửi công việc in. Ứng dụng hỗ trợ đa nền tảng (Windows, macOS, Linux) và sử dụng kết nối TCP/IP để giao tiếp với máy in.

## Yêu Cầu Hệ Thống

- **Hệ điều hành**: Windows, macOS hoặc Linux.
- **Máy in**: Máy in hỗ trợ giao thức ESC/POS qua TCP/IP (thường là port 9100).
- **Driver máy in**: Cần cài đặt driver máy in trước khi sử dụng.

## Cài Đặt

### Bước 1: Tải Xuống Ứng Dụng

- Tải file thực thi.

### Bước 2: Chuẩn Bị Thư Mục

- Tạo một thư mục riêng để chứa chương trình.
- Sao chép file thực thi (ví dụ: `printer-amd64.exe`) vào thư mục này.

### Bước 3: Khởi Động Ứng Dụng

- Nhấp đúp vào file thực thi để chạy (trên Windows) hoặc chạy từ terminal:
- Lần đầu tiên chạy, ứng dụng sẽ tự động khởi tạo các thư mục và file cấu hình cần thiết:
  - `config/`: Chứa file cấu hình.
  - `uploads/`: Chứa file tải lên tạm thời.
  - `logs/`: Chứa file log.

### Bước 4: Cấu Hình Máy In

- **Cài đặt driver máy in**: Đảm bảo máy in được cài đặt với chế độ kết nối TCP/IP. Trên Windows, vào **Settings > Devices > Printers & scanners** để kiểm tra.
- **Cấu hình port thủ công (nếu cần)**:
  - Vào **Printer Properties > Ports**.
  - Chọn **Add Port > New Port Type**.
  - Nhập địa chỉ IP và port của máy in (thường là port 9100).
  - Áp dụng và xác nhận.
- **File cấu hình**:
  - `config/device.json`: Lưu thông tin về máy đang chạy dịch vụ (ví dụ: vị trí, tên máy). Thay đổi để dễ dàng kiểm tra lỗi và log.
  - `config/config.json`: Lưu cấu hình máy in (IP, port, loại in). File này được quản lý tự động bởi hệ thống; không khuyến nghị chỉnh sửa thủ công để tránh lỗi kết nối.

## Chạy Ứng Dụng

- Ứng dụng sẽ chạy trên port 9099.

## Lưu Ý Vận Hành

- Máy cài đặt phải cùng mạng LAN với máy in và có kết nối internet để gửi tín hiệu.
- Đảm bảo firewall không chặn port 9100 hoặc 9099.
- Log được lưu trong `logs/app.log` với rotation tự động (tối đa 50MB, giữ 14 ngày).
- Ứng dụng tự động dọn dẹp file tạm trong `uploads/` lúc 5:00 AM hàng ngày.
- Nếu gặp lỗi, kiểm tra log hoặc khởi động lại dịch vụ.

## Xử Lý Sự Cố

- **Không tìm thấy máy in**: Kiểm tra `config/config.json` và đảm bảo máy in được cấu hình đúng.
- **Lỗi kết nối**: Kiểm tra IP/port và firewall.
- **File không in được**: Đảm bảo file là hình ảnh (JPEG/PNG) và kích thước hợp lý (tối đa 576 điểm ngang).
