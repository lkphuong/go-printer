package utils

import (
	"bytes"
	"fmt"
	"go-printer/internal/constants"
	"image"
	_ "image/jpeg"
	"log"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

type PrintRequest struct {
	IP         string
	FilePath   string
	Copies     string
	Result     chan error
	StartTime  time.Time
	RetryCount int
}

var printQueue = make(chan PrintRequest, 1000)

func startPrintWorker() {
	log.Println("Starting print worker...")
	fmt.Printf("Starting print worker...\n")
	go func() {
		for req := range printQueue {
			go func(r PrintRequest) {
				var err error
				for {
					err = printFile(r.IP, r.FilePath, r.Copies)
					if err == nil {
						break
					}
					r.RetryCount++
					if time.Since(r.StartTime) > 1*time.Hour {
						log.Printf("Job expired after 1 hour, retries: %d", r.RetryCount)
						fmt.Printf("Job expired after 1 hour, retries: %d\n", r.RetryCount)
						// xoá file tạm thời
						if err := os.Remove(r.FilePath); err != nil {
							log.Println("Xoá file tạm thời fail: ", err)
						}
						break
					}
					log.Printf("Retry %d after error: %v, sleeping 30s", r.RetryCount, err)
					fmt.Printf("Retry %d after error: %v, sleeping 30s\n", r.RetryCount, err)
					time.Sleep(30 * time.Second)
				}
				r.Result <- err
			}(req)
		}
	}()
}

func init() {
	startPrintWorker()
}

func PrintFileQueued(ip string, filePath string, copies string) error {
	req := PrintRequest{
		IP:         ip,
		FilePath:   filePath,
		Copies:     copies,
		Result:     make(chan error, 1),
		StartTime:  time.Now(),
		RetryCount: 0,
	}
	select {
	case printQueue <- req:
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

			fmt.Println("name: ", printerName)

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

func printFile(ip string, filePath string, copies string) error {

	numCopies := 1
	if copies != "" {
		if _, err := fmt.Sscan(copies, &numCopies); err != nil || numCopies < 1 {
			numCopies = 1
		}
	}

	// printer|192.168.1.100:9100
	printerConn, err := ConnectPrinter(ip)
	if err != nil {
		return err
	}
	defer printerConn.Close()

	const MaxPrinterDots = 576 // Số điểm tối đa của máy in (tùy thuộc vào máy in)

	// Kiểm tra status máy in
	err = printStatus(printerConn)
	if err != nil {
		log.Println("Lỗi trạng thái máy in: ", err)
		fmt.Printf("Lỗi trạng thái máy in: %v\n", err)
		return err
	}

	// Mở file ảnh
	imgFile, err := os.Open(filePath)
	if err != nil {
		log.Println("Mở file ảnh fail: ", err)
		fmt.Printf("Mở file ảnh fail: %v\n", err)
		return err
	}
	defer imgFile.Close()

	// Giải mã ảnh
	img, _, err := image.Decode(imgFile)
	if err != nil {
		log.Println("Giải mã ảnh fail: ", err)
		fmt.Printf("Giải mã ảnh fail: %v\n", err)
		return err
	}

	// Tạo lệnh in ảnh
	printCmd, err := printImageCommand(img, MaxPrinterDots)
	if err != nil {
		log.Println("Tạo lệnh in ảnh fail: ", err)
		fmt.Printf("Tạo lệnh in ảnh fail: %v\n", err)
		return err
	}

	// Gửi lệnh in nhiều bản
	for i := 0; i < numCopies; i++ {
		_, err = printerConn.Write(printCmd)
		// Thêm vài dòng trắng dòng trắng sau khi in ảnh
		printerConn.Write([]byte{0x0A, 0x0A, 0x0A, 0x0A})
		// conn.Write([]byte{0x0A, 0x0A})
		printerConn.Write([]byte{0x1D, 0x56, 0x00}) // Cắt giấy (GS V 0)

		if err != nil {
			log.Println("Gửi lệnh in fail: ", err)
			fmt.Printf("Gửi lệnh in fail: %v\n", err)
			return err
		}

		// await 0.5 second between copies
		time.Sleep(500 * time.Millisecond)
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

func ConnectPrinter(ip string) (net.Conn, error) {

	address := fmt.Sprintf("%s:%d", ip, 9100)

	printerConn, err := net.DialTimeout("tcp", address, 5*time.Second)
	if err != nil {
		log.Println("lỗi kết nối máy in: ", err)
		fmt.Printf("Lỗi kết nối máy in: %v\n", err)
		return nil, fmt.Errorf(constants.CONNECT_TIMEOUT)
	}
	return printerConn, nil
}

func printStatus(printerConn net.Conn) error {
	// Gửi lệnh truy vấn trạng thái: GS r 1 (0x1D 0x72 0x01)
	statusCmd := []byte{0x1D, 0x72, 0x01}
	_, err := printerConn.Write(statusCmd)
	if err != nil {
		return err
	}

	// Đọc phản hồi trạng thái từ máy in
	buffer := make([]byte, 1)
	printerConn.SetReadDeadline(time.Now().Add(1 * time.Second))
	n, err := printerConn.Read(buffer)
	if err != nil || n != 1 {
		return err
	}

	status := buffer[0]
	log.Printf("Trạng thái máy in: 0x%02X", status)

	// Log các lỗi dựa trên bit
	errors := []string{}
	if status&0x01 != 0 {
		errors = append(errors, constants.OFFLINE)
	}
	if status&0x02 != 0 {
		errors = append(errors, constants.OPEN)
	}
	if status&0x04 != 0 {
		errors = append(errors, constants.PAPER_JAM)
	}
	if status&0x08 != 0 {
		errors = append(errors, constants.PAPER_OUT)
	}
	if status&0x10 != 0 {
		errors = append(errors, constants.CONNECT_TIMEOUT)
	}
	if status&0x20 != 0 {
		errors = append(errors, constants.PAPER_NEAR_OUT)
	}
	if status&0x40 != 0 {
		errors = append(errors, constants.CUT_ERROR)
	}
	if status&0x80 != 0 {
		errors = append(errors, constants.UNKNOWN_ERROR)
	}

	if len(errors) > 0 {
		log.Println("Lỗi máy in: ", strings.Join(errors, ", "))
		return fmt.Errorf(constants.UNKNOWN_ERROR)
	}

	return nil

}
