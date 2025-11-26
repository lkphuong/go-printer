package constants

const (
	OK                         = "ok"
	PrintTypeKitchen PrintType = "kitchen"
	PrintTypeCashier PrintType = "cashier"
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
)
