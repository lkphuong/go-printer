Tạo một thư mục riêng để chứa chương trình.
Sao chép file thực thi printer-amd64.exe vào thư mục này.
Nhấp đúp chuột vào file printer-amd64.exe để khởi động.
Ở lần chạy đầu tiên, hệ thống sẽ tự động khởi tạo thư mục config cùng các file cấu hình cần thiết.
Trong folder configs

- device.json: Lưu trữ thông tin của máy đang cài đặt dịch vụ ví dụ máy in Ung Văn Khiêm, nên thay đổi khi cài đặt để dễ vận hành kiểm tra lỗi. Hỗ trợ ghi log và truy vết lỗi trong quá trình vận hành thực tế
- config.json: Lưu trữ cấu hình máy in trong cùng mạng LAN (IP, port, driver, cấu hình kết nối…). File này là cấu hình dùng chung, được hệ thống quản lý. Không khuyến nghị chỉnh sửa thủ công để tránh gây lỗi kết nối hoặc sai lệch cấu hình.

Lưu ý vận hành:

- Máy cài đặt phải chung mạng LAN với máy in, luôn đảm bảo thiết bị có internet để gửi tín hiệu tới máy in
