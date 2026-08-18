//go:build !windows

package utils

import "fmt"

// spoolerAvailable: Windows Print Spooler chỉ tồn tại trên Windows, nên các build
// khác (mac/linux) không thể in — printFile sẽ báo lỗi ngay từ services/print.go
// (SpoolerAvailable() == false) trước khi tới đây.
const spoolerAvailable = false

func printViaSpooler(printerName string, printCmd []byte, numCopies int) error {
	return fmt.Errorf("in qua Windows Spooler không được hỗ trợ trên hệ điều hành này")
}
