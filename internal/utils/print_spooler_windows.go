//go:build windows

package utils

import (
	"fmt"
	"log"
	"syscall"
	"time"
	"unsafe"
)

// spoolerAvailable báo cho printFile biết build này có thể in qua Windows Spooler.
const spoolerAvailable = true

var (
	modWinSpool          = syscall.NewLazyDLL("winspool.drv")
	procOpenPrinterW     = modWinSpool.NewProc("OpenPrinterW")
	procClosePrinter     = modWinSpool.NewProc("ClosePrinter")
	procStartDocPrinterW = modWinSpool.NewProc("StartDocPrinterW")
	procEndDocPrinter    = modWinSpool.NewProc("EndDocPrinter")
	procStartPagePrinter = modWinSpool.NewProc("StartPagePrinter")
	procEndPagePrinter   = modWinSpool.NewProc("EndPagePrinter")
	procWritePrinter     = modWinSpool.NewProc("WritePrinter")
)

// docInfo1W ánh xạ struct DOC_INFO_1 của WinAPI (winspool.h).
type docInfo1W struct {
	pDocName    *uint16
	pOutputFile *uint16
	pDatatype   *uint16
}

// printViaSpooler gửi printCmd (kèm feed + cut) qua Windows Print Spooler bằng datatype
// "RAW", để driver máy in đã cài trên Windows lo việc chờ máy in rảnh, đóng gói và
// truyền dữ liệu thật tới thiết bị, thay vì tự bắn thẳng qua socket TCP như trước.
//
// Lỗi ở đây KHÔNG được bọc thành retryablePrintError: khi đã trao lệnh in cho Spooler,
// service coi như xong việc — lỗi từ đây trở đi (máy Windows/driver treo, offline,
// hết giấy...) thuộc phạm vi trách nhiệm của Windows/driver, service chỉ là proxy
// trung gian nhận lệnh in, không tự lặp lại (retry) ở tầng ứng dụng cho các lỗi này.
func printViaSpooler(printerName string, printCmd []byte, numCopies int, sourceFilePath string) error {
	namePtr, err := syscall.UTF16PtrFromString(printerName)
	if err != nil {
		return fmt.Errorf("tên máy in không hợp lệ %q: %w", printerName, err)
	}

	var hPrinter syscall.Handle
	r1, _, errno := procOpenPrinterW.Call(uintptr(unsafe.Pointer(namePtr)), uintptr(unsafe.Pointer(&hPrinter)), 0)
	if r1 == 0 {
		return fmt.Errorf("OpenPrinter %q lỗi: %w", printerName, errno)
	}
	defer procClosePrinter.Call(uintptr(hPrinter))

	payload := make([]byte, 0, len(printCmd)+7)
	payload = append(payload, printCmd...)
	payload = append(payload, 0x0A, 0x0A, 0x0A, 0x0A) // feed giấy
	payload = append(payload, 0x1D, 0x56, 0x00)       // cắt giấy

	if err := dumpPrinterPayload(payload, sourceFilePath); err != nil {
		log.Println("Lưu file dump byte gửi máy in thất bại:", err)
	}

	for i := 0; i < numCopies; i++ {
		docName := fmt.Sprintf("go-printer copy %d/%d", i+1, numCopies)
		if err := writeSpoolerDocument(hPrinter, docName, payload); err != nil {
			return err
		}

		// await 0.5 second between copies, giống hành vi socket cũ
		time.Sleep(500 * time.Millisecond)
	}

	return nil
}

func writeSpoolerDocument(hPrinter syscall.Handle, docName string, data []byte) error {
	docNamePtr, err := syscall.UTF16PtrFromString(docName)
	if err != nil {
		return fmt.Errorf("tên document không hợp lệ: %w", err)
	}
	dataTypePtr, err := syscall.UTF16PtrFromString("RAW")
	if err != nil {
		return fmt.Errorf("datatype không hợp lệ: %w", err)
	}

	di := docInfo1W{pDocName: docNamePtr, pOutputFile: nil, pDatatype: dataTypePtr}

	r1, _, errno := procStartDocPrinterW.Call(uintptr(hPrinter), 1, uintptr(unsafe.Pointer(&di)))
	if r1 == 0 {
		return fmt.Errorf("StartDocPrinter lỗi: %w", errno)
	}
	defer procEndDocPrinter.Call(uintptr(hPrinter))

	r1, _, errno = procStartPagePrinter.Call(uintptr(hPrinter))
	if r1 == 0 {
		return fmt.Errorf("StartPagePrinter lỗi: %w", errno)
	}
	defer procEndPagePrinter.Call(uintptr(hPrinter))

	if len(data) == 0 {
		return nil
	}

	var written uint32
	r1, _, errno = procWritePrinter.Call(
		uintptr(hPrinter),
		uintptr(unsafe.Pointer(&data[0])),
		uintptr(len(data)),
		uintptr(unsafe.Pointer(&written)),
	)
	if r1 == 0 {
		return fmt.Errorf("WritePrinter lỗi: %w", errno)
	}
	if int(written) != len(data) {
		return fmt.Errorf("WritePrinter ghi thiếu: %d/%d byte", written, len(data))
	}

	return nil
}
