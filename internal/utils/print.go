package utils

import (
	"bytes"
	"errors"
	"fmt"
	"go-printer/internal/constants"
	"go-printer/internal/logger"
	"image"
	_ "image/jpeg"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type PrintRequest struct {
	Printer  string // tên máy in đã cài trên Windows, dùng để in qua Spooler/driver
	FilePath string
	Copies   string
}

var printQueue = make(chan PrintRequest, 1000)

type printFunc func(printer string, filePath string, copies string) error

type retryPolicy struct {
	interval time.Duration
	timeout  time.Duration
	now      func() time.Time
	sleep    func(time.Duration)
}

type retryablePrintError struct {
	cause error
}

func (e *retryablePrintError) Error() string {
	return e.cause.Error()
}

func (e *retryablePrintError) Unwrap() error {
	return e.cause
}

func newRetryablePrintError(err error) error {
	if err == nil {
		return nil
	}
	return &retryablePrintError{cause: err}
}

func isRetryablePrintError(err error) bool {
	var retryableErr *retryablePrintError
	return errors.As(err, &retryableErr)
}

func defaultPrintRetryPolicy() retryPolicy {
	return retryPolicy{
		interval: 30 * time.Second,
		timeout:  5 * time.Minute,
		now:      time.Now,
		sleep:    time.Sleep,
	}
}

func startPrintWorker() {
	log.Println("Starting print worker...")
	go runPrintWorker(printQueue, func(req PrintRequest) {
		_ = processPrintRequest(req, printFile, defaultPrintRetryPolicy())
	})
}

func runPrintWorker(queue <-chan PrintRequest, process func(PrintRequest)) {
	for req := range queue {
		func() {
			defer func() {
				if rec := recover(); rec != nil {
					log.Printf("print job panic recovered: %v", rec)
					logger.LogPrint(constants.UNKNOWN_ERROR, 500, fmt.Sprintf("panic: %v", rec))
				}
			}()
			process(req)
		}()
	}
}

func processPrintRequest(req PrintRequest, print printFunc, policy retryPolicy) error {
	startedAt := policy.now()
	retryCount := 0

	for {
		err := print(req.Printer, req.FilePath, req.Copies)
		if err == nil {
			logger.LogPrint(constants.OK, 200, "")
			return nil
		}

		if !isRetryablePrintError(err) {
			logger.LogPrint(err.Error(), 500, err.Error())
			return err
		}

		retryCount++
		if policy.now().Sub(startedAt) >= policy.timeout {
			log.Printf("Job expired after %s, retries: %d", policy.timeout, retryCount)
			logger.LogPrint(
				constants.PRINT_FAILED,
				500,
				fmt.Sprintf("job expired after %s, retries: %d", policy.timeout, retryCount),
			)
			if removeErr := os.Remove(req.FilePath); removeErr != nil && !errors.Is(removeErr, os.ErrNotExist) {
				log.Println("Xoá file tạm thời fail: ", removeErr)
				logger.LogPrint(
					constants.PRINT_FAILED,
					500,
					fmt.Sprintf("failed to remove temporary file: %v", removeErr),
				)
			}
			logger.LogPrint(err.Error(), 500, err.Error())
			return err
		}

		log.Printf("Retry %d after error: %v, sleeping %s", retryCount, err, policy.interval)
		logger.LogPrint(
			constants.PRINT_FAILED,
			500,
			fmt.Sprintf("retry %d after error: %v", retryCount, err),
		)
		policy.sleep(policy.interval)
	}
}

func init() {
	startPrintWorker()
}

func PrintFileQueued(printer string, filePath string, copies string) error {
	req := PrintRequest{
		Printer:  printer,
		FilePath: filePath,
		Copies:   copies,
	}
	return enqueuePrintRequest(printQueue, req)
}

func enqueuePrintRequest(queue chan<- PrintRequest, req PrintRequest) error {
	select {
	case queue <- req:
		return nil // Gửi thành công, không đợi
	default:
		return fmt.Errorf(constants.QUEUE_FULL) // Queue đầy, trả về lỗi
	}
}

func GetPrinters() ([]string, error) {
	switch runtime.GOOS {
	case "windows":
		cmd := exec.Command("powershell", " Get-Printer | Select Name")
		var out bytes.Buffer
		cmd.Stdout = &out
		err := cmd.Run()
		if err != nil {
			return nil, err
		}
		// return array of printer names:portName
		lines := []string{}

		if len(strings.Split(strings.TrimSpace(out.String()), "\n")) < 1 {
			return lines, nil
		}

		for _, line := range strings.Split(strings.TrimSpace(out.String()), "\n")[2:] {
			// tên máy in trước chữ local
			printerName := strings.TrimSpace(strings.Split(line, "Local")[0])

			lines = append(lines, printerName)
		}

		return lines, nil // iTP86|192.168.1.100:9100

	default:
		cmd := exec.Command("lpstat", "-p")
		var out bytes.Buffer
		cmd.Stdout = &out
		err := cmd.Run()
		if err != nil {
			return nil, err
		}
		lines := strings.Split(strings.TrimSpace(out.String()), "\n")
		var printers []string
		for _, line := range lines {
			if strings.HasPrefix(line, "printer ") {
				fields := strings.Fields(line)
				if len(fields) > 1 {
					printers = append(printers, fields[1])
				}
			}
		}
		return printers, nil
	}
}

func printFile(printer string, filePath string, copies string) error {

	numCopies := 1
	if copies != "" {
		if _, err := fmt.Sscan(copies, &numCopies); err != nil || numCopies < 1 {
			numCopies = 1
		}
	}

	const MaxPrinterDots = 576 // Số điểm tối đa của máy in (tùy thuộc vào máy in)

	if printer == "" {
		err := fmt.Errorf("thiếu tên máy in để gửi lệnh in qua Windows Spooler")
		log.Println(err)
		logger.LogPrint(constants.NO_PRINTER, 500, err.Error())
		return err
	}

	// Mở file ảnh
	imgFile, err := os.Open(filePath)
	if err != nil {
		log.Println("Mở file ảnh fail: ", err)
		logger.LogPrint(constants.READ_FILE_FAILED, 500, err.Error())
		return err
	}
	defer imgFile.Close()

	// Giải mã ảnh
	img, _, err := image.Decode(imgFile)
	if err != nil {
		log.Println("Giải mã ảnh fail: ", err)
		logger.LogPrint(constants.READ_FILE_FAILED, 500, "image decode failed: "+err.Error())
		return err
	}

	// Tạo lệnh in ảnh
	printCmd, err := printImageCommand(img, MaxPrinterDots)
	if err != nil {
		log.Println("Tạo lệnh in ảnh fail: ", err)
		logger.LogPrint(constants.PRINT_FAILED, 500, "build print command failed: "+err.Error())
		return err
	}

	// In 100% qua Windows Spooler + driver máy in — không còn đường socket TCP thô.
	// Spooler/driver tự lo việc chờ máy in rảnh, tránh tình trạng gửi bitmap thô làm
	// tràn buffer máy in (in nửa tờ, nửa trắng). Lỗi từ đây trở đi thuộc về
	// máy Windows/driver, không được retry ở tầng ứng dụng (xem printViaSpooler).
	if err := printViaSpooler(printer, printCmd, numCopies, filePath); err != nil {
		log.Println("In qua Windows Spooler fail: ", err)
		logger.LogPrint(constants.PRINT_FAILED, 500, "print via spooler failed: "+err.Error())
		return err
	}

	return nil
}

func resizeImage(img image.Image, newWidth int, newHeight int) image.Image {
	bounds := img.Bounds()                                         // lấy kích thước ban đầu của ảnh
	oldWidth := bounds.Dx()                                        // chiều rộng ban đầu
	oldHeight := bounds.Dy()                                       // chiều cao ban đầu
	newImg := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight)) // tạo ảnh mới với kích thước mới
	// thực hiện resize ảnh
	for y := 0; y < newHeight; y++ { // duyệt qua từng pixel theo chiều cao
		for x := 0; x < newWidth; x++ { // duyệt qua từng pixel theo chiều rộng
			oldX := x * oldWidth / newWidth      // tính toán vị trí pixel tương ứng trong ảnh gốc
			oldY := y * oldHeight / newHeight    // tính toán vị trí pixel tương ứng trong ảnh gốc
			newImg.Set(x, y, img.At(oldX, oldY)) // gán giá trị pixel từ ảnh gốc sang ảnh mới
		}
	}
	return newImg // trả về ảnh đã được resize
}

func printImageCommand(img image.Image, maxWidth int) ([]byte, error) {
	bounds := img.Bounds() // Lấy kích thước ảnh
	w := bounds.Dx()       // Chiều rộng ảnh
	h := bounds.Dy()       // Chiều cao ảnh

	if w > maxWidth {
		ratio := float64(maxWidth) / float64(w) // Tính tỉ lệ resize
		newH := int(float64(h) * ratio)         // Tính chiều cao mới theo tỉ lệ
		img = resizeImage(img, maxWidth, newH)  // Resize ảnh
		w = maxWidth                            // Cập nhật lại chiều rộng
		h = newH                                // Cập nhật lại chiều cao
	}

	byteWidth := (w + 7) / 8 // Chiều rộng tính theo byte (mỗi byte chứa 8 pixel)

	xL := byte(byteWidth & 0xFF) // LSB của chiều rộng
	xH := byte(byteWidth >> 8)   // MSB của chiều rộng
	yL := byte(h & 0xFF)         // LSB của chiều cao
	yH := byte(h >> 8)           // MSB của chiều cao

	command := []byte{0x1D, 0x76, 0x30, 0x00} // Khởi đầu lệnh in ảnh

	command = append(command, xL, xH, yL, yH) // Thêm thông tin chiều rộng và chiều cao

	// Chuyển đổi ảnh sang định dạng đen trắng 1 bit
	data := make([]byte, byteWidth*h) // Dữ liệu ảnh in
	for y := 0; y < h; y++ {          // Duyệt qua từng hàng
		for x := 0; x < w; x++ { // Duyệt qua từng cột
			r, g, b, _ := img.At(x, y).RGBA() // Lấy giá trị màu của pixel
			// Chuyển đổi sang thang độ xám
			grayValue := 0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b)

			if grayValue < 32768 { // Ngưỡng để xác định đen trắng
				byteIndex := y*byteWidth + x/8        // Vị trí byte trong dữ liệu
				bitPosition := 7 - (x % 8)            // Vị trí bit trong byte
				data[byteIndex] |= (1 << bitPosition) // Đặt bit thành 1 (đen)
			}
		}
	}

	command = append(command, data...) // Thêm dữ liệu ảnh vào lệnh in

	return command, nil
}

// SpoolerAvailable báo cho tầng service biết build hiện tại có thể in qua Windows Spooler.
func SpoolerAvailable() bool {
	return spoolerAvailable
}

// dumpPrinterPayload lưu nguyên văn (không sửa 1 byte nào) dữ liệu vừa gửi cho máy in
// qua WritePrinter — bao gồm cả feed giấy + lệnh cắt được nối thêm trong printViaSpooler.
// Đây chính là nội dung mà Windows Spooler ghi vào file .SPL (datatype RAW), nên đối
// chiếu 2 file này (fc /b hoặc hex diff) sẽ chứng minh được service có gửi đúng dữ liệu
// hay không, tách bạch lỗi service khỏi lỗi máy in/driver.
//
// File dump được đặt tên trùng với file gốc trong uploads/ (chỉ đổi đuôi thành .bin)
// để dễ đối soát 1-1 giữa ảnh gốc và dữ liệu thực tế đã gửi cho máy in — tempFilePath
// đã có timestamp + random suffix riêng cho từng job (newTempUploadPath) nên không lo
// trùng tên giữa các lần in.
func dumpPrinterPayload(payload []byte, sourceFilePath string) error {
	dumpDir := filepath.Join("uploads", "spool-dumps")
	if err := os.MkdirAll(dumpDir, 0755); err != nil {
		return err
	}

	base := filepath.Base(sourceFilePath)
	dumpName := strings.TrimSuffix(base, filepath.Ext(base)) + ".bin"

	return os.WriteFile(filepath.Join(dumpDir, dumpName), payload, 0644)
}
