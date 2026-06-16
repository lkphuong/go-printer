package constants

const (
	OK                         = "ok"
	PrintTypeKitchen PrintType = "kitchen"
	PrintTypeCashier PrintType = "cashier"
	API_KEY                    = ""
	MONGODB_URI                = ""
)

type PrintType string

const (
	PRINT_FAILED    = "print_failed"    // lỗi in
	PRINT_NOT_FOUND = "print_not_found" // không tìm thấy máy in
	CONNECT_TIMEOUT = "connect_timeout" // lỗi kết nối
	OFFLINE         = "offline"         // máy in offline
	OPEN            = "open"            // nắp máy in mở
	PAPER_OUT       = "paper_out"       // hết giấy
	PAPER_JAM       = "paper_jam"       // Nút cấp giấy được nhấn
	PAPER_NEAR_OUT  = "paper_near_out"  // Giấy gần hết
	CUT_ERROR       = "cut_error"       // lỗi cắt giấy
	UNKNOWN_ERROR   = "unknown_error"   // lỗi không xác định
	QUEUE_FULL      = "queue_full"      // hàng đợi in đầy
	API_KEY_FAIL    = "api_key_failed"  // thiếu API key

	GET_PRINTERS_FAILED = "get_printers_failed" // không lấy được danh sách máy in
	READ_CONFIG_FAILED  = "read_config_failed"  // lỗi đọc file config
	WRITE_CONFIG_FAILED = "write_config_failed" // lỗi ghi file config
	CONFIG_NOT_FOUND    = "config_not_found"    // không tìm thấy config máy in
	CLEAR_CACHE_FAILED  = "clear_cache_failed"  // lỗi xoá cache
	INVALID_REQUEST     = "invalid_request"     // dữ liệu request không hợp lệ
	FORM_PARSE_FAILED   = "form_parse_failed"   // lỗi phân tích form
	NO_FILES_UPLOADED   = "no_files_uploaded"   // không có file được tải lên
	NO_PRINTER          = "no_printer"          // không chỉ định máy in
	SAVE_FILE_FAILED    = "save_file_failed"    // lỗi lưu file tạm
	READ_FILE_FAILED    = "read_file_failed"    // lỗi mở file
	QUEUED              = "queued"              // job đã được đưa vào hàng đợi
)
