# Go Printer — Hướng dẫn vận hành và xử lý sự cố

Tài liệu này hướng dẫn cách **vận hành** dịch vụ Go Printer và **xử lý sự cố** khi máy in không hoạt động. Viết cho người không rành kỹ thuật.

> Đọc hết một lượt trước khi bắt đầu. Làm theo đúng thứ tự.

---

## 1. Tổng quan nhanh

Go Printer là một **dịch vụ chạy ngầm trên Windows**. Phần mềm bán hàng gửi **ảnh hóa đơn (JPEG)** qua mạng đến máy tính chạy Go Printer; dịch vụ chuyển ảnh đó tới **máy in nhiệt** để in ra giấy.

Ba điều cần nhớ:

- Dịch vụ chỉ chạy trên máy tính **Windows**.
- Máy in phải là **máy in nhiệt có kết nối mạng (LAN/WiFi)**.
- File gửi đi để in phải là **ảnh JPEG** (không phải PDF, Word, PNG).

---

## 2. Cài đặt và khởi động

### 2.1. Kết nối máy in vào mạng và lấy địa chỉ IP

1. Cắm dây LAN hoặc kết nối máy in vào WiFi.
2. Lấy địa chỉ IP máy in (in trang self-test bằng cách giữ nút **FEED** khi bật nguồn, hoặc xem trong menu Network trên máy in).
3. **Ghi lại địa chỉ IP**, ví dụ `192.168.1.50`.

> Mẹo: Nên đặt máy in dùng **IP tĩnh** trên router để IP không bị đổi sau khi khởi động lại.

### 2.2. Thêm máy in bằng cổng Standard TCP/IP

Đây là bước **bắt buộc** để dịch vụ tra được địa chỉ máy in.

1. Mở **Settings** → **Bluetooth & devices** → **Printers & scanners** → **Add device**.
2. Bấm **"Add manually"** → chọn **"Add a printer using a TCP/IP address or hostname"** → **Next**.
3. Nhập **địa chỉ IP máy in** đã ghi ở bước 2.1 → **Next**.
4. Chọn driver (driver của hãng, hoặc **Generic / Text Only**) → **Next**.
5. **Đặt tên cho máy in** và **ghi nhớ chính xác tên này** — phải nhập đúng y hệt khi gửi lệnh in. Nên đặt tên **không dấu, không khoảng trắng**, ví dụ `MayInBep`.
6. Bấm **Finish**.

### 2.3. Cài dịch vụ

1. Tạo thư mục cố định, ví dụ **`C:\GoPrinter`**, và sao chép file **`printer-amd64.exe`** vào đó.
2. **Bấm chuột phải** vào file → **"Run as administrator"**.
3. Chương trình tự cài và khởi động dịch vụ Windows tên **"Printer AMD64 Service"**, sau đó tự thoát. Đây là hành vi bình thường.

Lần chạy đầu, chương trình tự tạo các thư mục ngay cạnh file `.exe`: `uploads/` (file tạm), `config/` (cấu hình), `logs/` (nhật ký).

### 2.4. Kiểm tra cấu hình

Mở `config/device.json` (cạnh file `.exe`) và đổi nhãn thiết bị cho dễ theo dõi log:

```json
{
  "location": "quay-1"
}
```

---

## 3. Kiểm tra dịch vụ đã chạy chưa

**Cách 1 — Trong Services:** Nhấn **Windows + R**, gõ `services.msc`, bấm **OK**. Tìm **"Printer AMD64 Service"**, cột **Status** phải là **Running**.

**Cách 2 — Bằng lệnh:** Mở **Command Prompt** và chạy (thay `<API_KEY_CUA_BAN>` bằng khóa thật được cấp):

```bash
curl -H "API-Key: <API_KEY_CUA_BAN>" http://localhost:9099/api/v1/printers
```

Nếu chạy đúng, kết quả trả về danh sách tên máy in. Tên máy in đã đặt ở bước 2.2 phải xuất hiện trong danh sách.

---

## 4. Quản lý dịch vụ

Mở `services.msc`, tìm **"Printer AMD64 Service"**, bấm chuột phải để:

- **Start** — khởi động dịch vụ.
- **Stop** — dừng dịch vụ.
- **Restart** — khởi động lại (dùng sau khi đổi cấu hình).

**Xem nhật ký** khi cần kiểm tra sự cố (các file cạnh `printer-amd64.exe`):

- `logs/app.log` — nhật ký chi tiết hoạt động.
- `service.log` — nhật ký mức dịch vụ.

> Hệ thống **tự dọn dẹp vào khoảng 5 giờ sáng**: xóa file tạm trong `uploads/` và nhật ký cũ.

---

## 5. Xử lý sự cố

### 5.1. Quy trình chẩn đoán (làm lần lượt từ trên xuống)

Khi máy in không in được, kiểm tra **lần lượt từng bước**. Dừng lại ngay tại bước phát hiện ra vấn đề:

1. **Dịch vụ có đang chạy không?**
   - Mở `services.msc`, tìm **"Printer AMD64 Service"**, xem cột **Status**.
   - Nếu **không phải Running**: bấm chuột phải → **Start**. Nếu không có trong danh sách, chạy lại **mục 2.3** bằng quyền Administrator.

2. **API có phản hồi không?** (Chỉ làm khi dịch vụ đã Running.)
   - Mở **Command Prompt**, chạy:

     ```bash
     curl -H "API-Key: <API_KEY_CUA_BAN>" http://localhost:9099/api/v1/printers
     ```

   - Trả về **danh sách máy in** → dịch vụ hoạt động bình thường, chuyển sang bước 4.
   - **Không phản hồi / báo lỗi kết nối** → quay lại bước 1 hoặc xem nhật ký ở bước 3.

3. **Đọc nhật ký để tìm lỗi.** (Khi dịch vụ Running nhưng vẫn không in được.)
   - Mở thư mục cài đặt (ví dụ `C:\GoPrinter`) → thư mục **`logs/`** → file **`app.log`**.
   - Xem **dòng dưới cùng** (mới nhất), đối chiếu với **bảng mã lỗi ở mục 6.2**.

4. **Kiểm tra mạng giữa máy tính và máy in.** (Khi gặp `connect_timeout` / `print_not_found`.)
   - Đảm bảo máy tính và máy in **cùng một mạng nội bộ**.
   - Thử ping IP máy in (thay bằng IP thật):

     ```bash
     ping 192.168.1.50
     ```

   - Ping **không thông** → kiểm tra dây mạng/WiFi, máy in đã bật chưa, IP có bị đổi không (nên đặt IP tĩnh).
   - Ping **thông** nhưng vẫn không in → kiểm tra tên máy in nhập có **đúng y hệt** tên đã đặt, và máy in được thêm bằng **cổng TCP/IP**.

### 6.2. Bảng mã lỗi

Khi gửi lệnh in, hệ thống trả về **mã lỗi** giúp xác định nguyên nhân:

| Mã lỗi              | Nghĩa                             | Cách xử lý                                                                  |
| ------------------- | --------------------------------- | --------------------------------------------------------------------------- |
| `api_key_failed`    | Thiếu hoặc sai `API-Key`          | Kiểm tra đã gửi đúng header `API-Key` với giá trị được cấp                  |
| `print_not_found`   | Không tra được IP máy in          | Kiểm tra tên máy in nhập đúng y hệt; máy in phải được thêm bằng cổng TCP/IP |
| `connect_timeout`   | Không kết nối được máy in         | Kiểm tra máy in đã bật, cùng mạng; thử ping IP máy in                       |
| `offline`           | Máy in đang offline               | Bật lại máy in, kiểm tra dây mạng/WiFi                                      |
| `open`              | Nắp máy in đang mở                | Đóng nắp máy in lại                                                         |
| `paper_out`         | Hết giấy                          | Lắp cuộn giấy mới                                                           |
| `paper_near_out`    | Giấy gần hết                      | Chuẩn bị thay cuộn giấy mới                                                 |
| `paper_jam`         | Kẹt giấy / nút cấp giấy được nhấn | Mở nắp, gỡ giấy kẹt, lắp lại giấy                                           |
| `cut_error`         | Lỗi cắt giấy                      | Kiểm tra bộ phận cắt, gỡ giấy kẹt ở dao cắt                                 |
| `queue_full`        | Hàng đợi in đã đầy                | Chờ các lệnh in trước hoàn tất rồi gửi lại                                  |
| `no_printer`        | Chưa chỉ định máy in              | Bổ sung trường `printer` trong lệnh in                                      |
| `no_files_uploaded` | Chưa đính kèm file                | Bổ sung trường `file` (ảnh JPEG) trong lệnh in                              |
| `unknown_error`     | Lỗi không xác định                | Xem `logs/app.log` để biết chi tiết                                         |

---

## 7. Lưu ý bảo mật

- **Chỉ chạy trong mạng nội bộ.** Tuyệt đối **không mở cổng `9099` ra Internet** (không port-forward trên router).
