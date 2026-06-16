# Go Printer — Hướng dẫn cài đặt từ A đến Z

Tài liệu này hướng dẫn **từng bước, dễ hiểu cho người không rành kỹ thuật**, để cài đặt và vận hành dịch vụ Go Printer trên máy tính Windows.

> Đọc hết một lượt trước khi bắt đầu. Mỗi bước nên làm theo đúng thứ tự.

---

## 1. Go Printer là gì?

Go Printer là một **dịch vụ chạy ngầm trên máy tính Windows**. Nó hoạt động như sau:

1. Một phần mềm khác (ví dụ phần mềm bán hàng) **gửi lệnh in qua mạng** đến máy tính chạy Go Printer.
2. Go Printer nhận **ảnh hóa đơn (định dạng JPEG)** và gửi thẳng tới **máy in nhiệt** qua mạng nội bộ.
3. Máy in nhiệt in hóa đơn ra giấy.

**Cần nhớ 3 điều quan trọng:**

- Dịch vụ chỉ in được trên máy tính **Windows**.
- Máy in phải là **máy in nhiệt có kết nối mạng (LAN/WiFi)**, giao tiếp chuẩn ESC/POS qua cổng `9100`.
- File gửi đi để in phải là **ảnh JPEG** (không phải PDF, không phải Word).

---

## 2. Cần chuẩn bị gì?

Hãy kiểm tra đủ các mục sau trước khi bắt đầu (danh sách kiểm tra):

- [ ] Một máy tính chạy **Windows** (bắt buộc, nếu không thì không in được).
- [ ] Một **máy in nhiệt hỗ trợ mạng** (cắm dây LAN hoặc WiFi), giao tiếp ESC/POS qua cổng `9100`.
- [ ] Máy in và máy tính **nằm cùng một mạng nội bộ** (cùng router/WiFi).
- [ ] File chương trình **`printer-amd64.exe`** (nằm trong thư mục `build/` của dự án).
- [ ] Quyền **Administrator** trên máy tính Windows (để cài dịch vụ).

---

## 3. Bước 1 — Kết nối máy in vào mạng và lấy địa chỉ IP

Máy tính cần biết **địa chỉ IP** của máy in để gửi lệnh in.

1. Cắm dây mạng LAN vào máy in, hoặc kết nối máy in vào WiFi (xem hướng dẫn của hãng máy in).
2. Lấy địa chỉ IP của máy in. Có hai cách phổ biến:
   - **In trang cấu hình (self-test):** nhiều máy in nhiệt cho phép giữ nút **FEED** (cấp giấy) khi bật nguồn để in ra một trang thông tin có ghi địa chỉ IP.
   - **Xem trên màn hình máy in** (nếu máy có màn hình): vào menu mạng (Network) để xem IP.
3. **Ghi lại địa chỉ IP** này, ví dụ: `192.168.1.50`. Bạn sẽ cần nó ở bước sau.

> Mẹo: Địa chỉ IP thường có dạng `192.168.x.x`. Nên đặt máy in dùng **IP tĩnh** trên router để IP không bị đổi sau khi khởi động lại.

---

## 4. Bước 2 — Cài driver máy in trên Windows

1. Tải driver của máy in từ **trang web chính hãng** (theo đúng model máy in của bạn).
2. Chạy file cài đặt driver vừa tải và làm theo hướng dẫn trên màn hình (bấm **Next** đến khi xong).
3. Nếu không tìm thấy driver riêng, bạn có thể dùng driver chung **"Generic / Text Only"** có sẵn trong Windows ở bước tiếp theo.

---

## 5. Bước 3 — Thêm máy in bằng cổng Standard TCP/IP (rất quan trọng)

Đây là bước **bắt buộc** để Go Printer tra được địa chỉ IP của máy in. Máy in phải được thêm vào Windows bằng **cổng TCP/IP**.

1. Mở **Settings** (Cài đặt) → **Bluetooth & devices** → **Printers & scanners**.
2. Bấm **Add device** (Thêm thiết bị).
3. Chờ vài giây, bấm dòng **"Add manually"** (hoặc _"The printer that I want isn't listed"_).
4. Chọn **"Add a printer using a TCP/IP address or hostname"** → bấm **Next**.
5. Tại ô **Hostname or IP address**, nhập **địa chỉ IP máy in** đã ghi ở Bước 1 (ví dụ `192.168.1.50`) → bấm **Next**.
6. Chọn driver (driver đã cài ở Bước 2, hoặc **Generic / Text Only**) → bấm **Next**.
7. **Đặt tên cho máy in** (ô _Printer name_). **Ghi nhớ chính xác tên này** — bạn sẽ phải nhập đúng y hệt khi gửi lệnh in. Ví dụ: `MayInBep`.
8. Bấm **Finish** để hoàn tất.

> Lưu ý: Tên máy in phân biệt chính xác từng ký tự. Nên đặt tên **không dấu, không khoảng trắng** cho dễ dùng (ví dụ `MayInBep`, `MayInQuay`).

---

## 6. Bước 4 — Cài dịch vụ Go Printer

1. Tạo một thư mục cố định trên ổ đĩa, ví dụ: **`C:\GoPrinter`**.
2. Sao chép file **`printer-amd64.exe`** vào thư mục đó.
3. **Bấm chuột phải** vào `printer-amd64.exe` → chọn **"Run as administrator"** (Chạy với quyền quản trị).
4. Chương trình sẽ **tự cài đặt và khởi động** một dịch vụ Windows tên là **`PrinterAMD64`** ("Printer AMD64 Service"), sau đó tự thoát. Đây là hành vi bình thường.

> Cần quyền Administrator vì việc cài dịch vụ Windows đòi hỏi quyền quản trị. Nếu không chạy bằng quyền admin, dịch vụ sẽ không cài được.

Khi chạy lần đầu, chương trình tự tạo các thư mục **ngay cạnh file `.exe`**:

- `uploads/` — nơi lưu tạm file gửi đến để in.
- `config/` — chứa các file cấu hình (xem Bước 5).
- `logs/` — chứa nhật ký hoạt động.

---

## 7. Bước 5 — Kiểm tra cấu hình

Sau khi chạy lần đầu, mở thư mục `config/` (nằm cạnh file `.exe`).

### `config/device.json`

File này đặt **nhãn thiết bị** dùng trong nhật ký. Nội dung mặc định:

```json
{
  "location": "office"
}
```

Bạn có thể đổi `"office"` thành tên cửa hàng hoặc quầy của mình, ví dụ `"quay-1"`. Nên thay đổi để tổng hợp log để theo dõi thiết bị

### `config/config.json`

File này file này bỏ qua dùng flow.

---

## 8. Bước 6 — Kiểm tra dịch vụ đã chạy chưa

### Cách 1: Kiểm tra trong Services

1. Nhấn phím **Windows + R**, gõ `services.msc`, bấm **OK**.
2. Tìm dòng **"Printer AMD64 Service"**.
3. Cột **Status** phải hiển thị **Running** (Đang chạy).

### Cách 2: Kiểm tra bằng trình duyệt / lệnh

Dịch vụ chạy ở cổng **`9099`**. Mọi yêu cầu đều **bắt buộc** kèm header xác thực **`API-Key`**.

Mở **Command Prompt** và chạy (thay `<API_KEY_CUA_BAN>` bằng khóa thật được cấp):

```bash
curl -H "API-Key: <API_KEY_CUA_BAN>" http://localhost:9099/api/v1/printers
```

Nếu chạy đúng, kết quả trả về danh sách tên máy in dạng JSON. Tên máy in bạn vừa đặt ở Bước 3 phải xuất hiện trong danh sách này.

---

## 9. Bước 7 — Gán loại máy in và in thử

Tất cả địa chỉ đều bắt đầu bằng `http://localhost:9099/api/v1` và **bắt buộc** có header `API-Key`.

### 9.1. Gán loại cho máy in (tùy chọn)

Gán loại `kitchen` (bếp) hoặc `cashier` (thu ngân) cho máy in:

```bash
curl -X POST http://localhost:9099/api/v1/printers/config ^
  -H "API-Key: <API_KEY_CUA_BAN>" ^
  -H "Content-Type: application/json" ^
  -d "{\"printer_name\": \"MayInBep\", \"type\": [\"kitchen\"]}"
```

> Cả `printer_name` và `type` đều bắt buộc. `type` là một danh sách.

### 9.2. Gửi lệnh in thử

Lệnh in gửi theo dạng **multipart/form-data** (KHÔNG phải JSON), với các trường:

- `file` — **ảnh JPEG** cần in (bắt buộc).
- `printer` — **đúng tên máy in** đã đặt ở Bước 3 (bắt buộc).
- `copies` — số bản in (tùy chọn, mặc định là `1`).

Ví dụ in 1 bản ảnh `hoadon.jpg` ra máy `MayInBep`:

```bash
curl -X POST http://localhost:9099/api/v1/printers/jobs ^
  -H "API-Key: <API_KEY_CUA_BAN>" ^
  -F "printer=MayInBep" ^
  -F "copies=1" ^
  -F "file=@C:\duong-dan\hoadon.jpg"
```

> **Quan trọng:** Hệ thống chỉ in được **ảnh JPEG**. File PDF, Word, PNG sẽ không in được.

---

## 10. Quản lý dịch vụ

Mở `services.msc` (Windows + R → gõ `services.msc`), tìm **"Printer AMD64 Service"**, bấm chuột phải để:

- **Start** — khởi động dịch vụ.
- **Stop** — dừng dịch vụ.
- **Restart** — khởi động lại (dùng sau khi đổi cấu hình).

**Xem nhật ký (log)** khi cần kiểm tra sự cố — các file nằm cạnh `printer-amd64.exe`:

- `logs/app.log` — nhật ký chi tiết hoạt động (tự xoay vòng khi đạt 50 MB, giữ tối đa 14 ngày).
- `service.log` — nhật ký mức dịch vụ.

> Hệ thống **tự dọn dẹp vào khoảng 5 giờ sáng**: xóa file tạm trong `uploads/` và xóa nhật ký cũ hơn 14 ngày.

---

## 11. Xử lý sự cố

### 11.1. Quy trình chẩn đoán theo thứ tự (làm lần lượt từ trên xuống)

Khi máy in không in được, hãy kiểm tra **lần lượt từng bước** sau. Dừng lại ngay tại bước phát hiện ra vấn đề:

1. **Kiểm tra dịch vụ có đang chạy không.**
   - Nhấn **Windows + R**, gõ `services.msc`, bấm **OK**.
   - Tìm dòng **"Printer AMD64 Service"** và xem cột **Status**.
   - Nếu **không phải Running**: bấm chuột phải → **Start**. Nếu không có trong danh sách, chạy lại **Bước 4** bằng quyền Administrator.

2. **Kiểm tra API có phản hồi không.** (Chỉ làm khi dịch vụ đã Running.)
   - Mở **Command Prompt**, chạy (thay `<API_KEY_CUA_BAN>` bằng khóa thật):

     ```bash
     curl -H "API-Key: <API_KEY_CUA_BAN>" http://localhost:9099/api/v1/printers
     ```

   - Nếu trả về **danh sách máy in** (chấp nhận đúng `API-Key`) → dịch vụ đang hoạt động bình thường, chuyển sang bước 4.
   - Nếu **không phản hồi / báo lỗi kết nối** → dịch vụ chưa thật sự chạy, quay lại bước 1 hoặc xem nhật ký ở bước 3.

3. **Đọc nhật ký để tìm lỗi.** (Khi dịch vụ Running nhưng vẫn không in được.)
   - Mở **thư mục cài đặt** (nơi chứa `printer-amd64.exe`, ví dụ `C:\GoPrinter`).
   - Mở thư mục **`logs/`** → mở file **`app.log`**.
   - Xem **dòng dưới cùng** (mới nhất) để biết hệ thống đang báo lỗi gì, rồi đối chiếu với **bảng mã lỗi ở mục 11.2**.

4. **Kiểm tra mạng giữa máy tính và máy in.** (Khi nhật ký không có lỗi rõ ràng, hoặc gặp `connect_timeout` / `print_not_found`.)
   - Đảm bảo máy tính và máy in **nằm cùng một mạng nội bộ (LAN/WiFi)**.
   - Thử ping địa chỉ IP máy in (thay bằng IP thật ở Bước 1):

     ```bash
     ping 192.168.1.50
     ```

   - Nếu ping **không thông** → kiểm tra dây mạng/WiFi, máy in đã bật chưa, IP có bị đổi không (nên đặt IP tĩnh — xem mẹo ở Bước 1).
   - Nếu ping **thông** nhưng vẫn không in → kiểm tra tên máy in nhập có **đúng y hệt** tên đã đặt ở Bước 3, và máy in được thêm bằng **cổng TCP/IP**.

### 11.2. Bảng mã lỗi

Khi gửi lệnh in, hệ thống trả về **mã lỗi** giúp xác định nguyên nhân. Bảng dưới đây liệt kê các mã thường gặp:

| Mã lỗi              | Nghĩa                             | Cách xử lý                                                                                      |
| ------------------- | --------------------------------- | ----------------------------------------------------------------------------------------------- |
| `api_key_failed`    | Thiếu hoặc sai `API-Key`          | Kiểm tra đã gửi đúng header `API-Key` với giá trị được cấp                                      |
| `print_not_found`   | Không tra được IP máy in          | Kiểm tra tên máy in nhập có **đúng y hệt** tên ở Bước 3; máy in phải được thêm bằng cổng TCP/IP |
| `connect_timeout`   | Không kết nối được máy in         | Kiểm tra máy in đã bật, cùng mạng, và mở cổng `9100`; thử ping IP máy in                        |
| `offline`           | Máy in đang offline               | Bật lại máy in, kiểm tra dây mạng/WiFi                                                          |
| `open`              | Nắp máy in đang mở                | Đóng nắp máy in lại                                                                             |
| `paper_out`         | Hết giấy                          | Lắp cuộn giấy mới                                                                               |
| `paper_near_out`    | Giấy gần hết                      | Chuẩn bị thay cuộn giấy mới                                                                     |
| `paper_jam`         | Kẹt giấy / nút cấp giấy được nhấn | Mở nắp, gỡ giấy kẹt, lắp lại giấy                                                               |
| `cut_error`         | Lỗi cắt giấy                      | Kiểm tra bộ phận cắt, gỡ giấy kẹt ở dao cắt                                                     |
| `queue_full`        | Hàng đợi in đã đầy                | Chờ các lệnh in trước hoàn tất rồi gửi lại                                                      |
| `no_printer`        | Chưa chỉ định máy in              | Bổ sung trường `printer` trong lệnh in                                                          |
| `no_files_uploaded` | Chưa đính kèm file                | Bổ sung trường `file` (ảnh JPEG) trong lệnh in                                                  |
| `unknown_error`     | Lỗi không xác định                | Xem `logs/app.log` để biết chi tiết                                                             |

**Dịch vụ không tự chạy?** Mở `services.msc`, tìm "Printer AMD64 Service" và bấm **Start**. Nếu không có trong danh sách, chạy lại Bước 4 bằng quyền Administrator.

---

## 12. Cảnh báo bảo mật (đọc kỹ)

> **QUAN TRỌNG:** Khóa API (`API-Key`) và chuỗi kết nối cơ sở dữ liệu MongoDB hiện đang được **nhúng cứng bên trong file `printer-amd64.exe`**. Bất kỳ ai có được file `.exe` đều có thể trích xuất các thông tin này và có **toàn quyền** truy cập API cũng như cơ sở dữ liệu.

Khuyến nghị để giảm rủi ro:

- **Chỉ chạy trong mạng nội bộ.** Tuyệt đối **không mở cổng `9099` ra Internet** (không port-forward trên router).
- **Không chia sẻ file `.exe`** ra ngoài. Coi nó như một tệp chứa mật khẩu.
- Nếu cần phân phối cho nhiều nơi, hãy **đổi (rotate) khóa API và mật khẩu MongoDB**, và cân nhắc đưa các bí mật này ra file cấu hình/biến môi trường thay vì nhúng cứng.
- Trong mọi tài liệu hay ví dụ, **chỉ dùng placeholder** như `<API_KEY_CUA_BAN>`, không bao giờ ghi giá trị thật.

---

## Phụ lục — Thông tin kỹ thuật nhanh

| Hạng mục                      | Giá trị                                  |
| ----------------------------- | ---------------------------------------- |
| Cổng API của dịch vụ          | `9099` (mọi địa chỉ, CORS mở)            |
| Cổng giao tiếp máy in         | `9100` (raw ESC/POS, timeout 5 giây)     |
| Header xác thực               | `API-Key`                                |
| Đường dẫn gốc API             | `/api/v1`                                |
| Tên dịch vụ Windows           | `PrinterAMD64` ("Printer AMD64 Service") |
| Định dạng file in             | Ảnh **JPEG**                             |
| Thư mục runtime (cạnh `.exe`) | `uploads/`, `config/`, `logs/`           |

### Danh sách API

| Phương thức | Đường dẫn                          | Chức năng                                            |
| ----------- | ---------------------------------- | ---------------------------------------------------- |
| GET         | `/api/v1/printers`                 | Lấy danh sách máy in                                 |
| GET         | `/api/v1/printers/:printer/config` | Xem cấu hình loại của một máy in                     |
| POST        | `/api/v1/printers/config`          | Gán loại cho máy in (JSON: `printer_name`, `type[]`) |
| POST        | `/api/v1/printers/jobs`            | Gửi lệnh in (multipart: `file`, `printer`, `copies`) |
| DELETE      | `/api/v1/printers/cache`           | Xóa toàn bộ cấu hình (reset về rỗng)                 |

### Build lại từ mã nguồn (dành cho lập trình viên)

Yêu cầu Go 1.23 trở lên.

```bash
# Build cho Windows (AMD64)
make build-windows   # tạo build/printer-amd64.exe

# Chạy trực tiếp để phát triển
make run             # chạy go run ./cmd/server/main.go
```
