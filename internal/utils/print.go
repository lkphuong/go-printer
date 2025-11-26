package utils

import (
	"bytes"
	"fmt"
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

func GetPrinters() ([]string, error) {
	switch runtime.GOOS {
	case "windows":
		cmd := exec.Command("powershell", "Get-Printer")
		var out bytes.Buffer
		cmd.Stdout = &out
		err := cmd.Run()
		if err != nil {
			return nil, err
		}
		// return array of printer names:portName
		lines := []string{}
		for _, line := range strings.Split(strings.TrimSpace(out.String()), "\n")[2:] {
			fields := strings.Fields(line)
			printerName := fields[0]
			portName := fields[3]
			if !strings.HasPrefix(printerName, "---") {
				lines = append(lines, fmt.Sprintf("%s|%s:9100", printerName, portName))
			}

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

func PrintFile(printer string, filePath string, copies string) error {

	numCopies := 1
	if copies != "" {
		if _, err := fmt.Sscan(copies, &numCopies); err != nil || numCopies < 1 {
			numCopies = 1
		}
	}

	// printer|192.168.1.100:9100
	parts := strings.Split(printer, "|")
	if len(parts) != 2 {
		return fmt.Errorf("invalid printer format")
	}
	address := parts[1]

	printerConn, err := net.Dial("tcp", address)
	if err != nil {
		log.Fatalf("Lỗi kết nối máy in: %v", err)
	}
	defer printerConn.Close()

	const MaxPrinterDots = 576 // Số điểm tối đa của máy in (tùy thuộc vào máy in)

	// Kiểm tra status máy in
	err = printStatus(printerConn)
	if err != nil {
		log.Println("Lỗi trạng thái máy in: ", err)
		return err
	}

	// Mở file ảnh
	imgFile, err := os.Open(filePath)
	if err != nil {
		log.Println("Mở file ảnh fail: ", err)
		return err
	}
	defer imgFile.Close()

	// Giải mã ảnh
	img, _, err := image.Decode(imgFile)
	if err != nil {
		log.Println("Giải mã ảnh fail: ", err)
		return err
	}

	// Tạo lệnh in ảnh
	printCmd, err := printImageCommand(img, MaxPrinterDots)
	if err != nil {
		log.Println("Tạo lệnh in ảnh fail: ", err)
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
			return err
		}
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
		errors = append(errors, "Offline")
	}
	if status&0x02 != 0 {
		errors = append(errors, "Nắp mở")
	}
	if status&0x04 != 0 {
		errors = append(errors, "Nút cấp giấy được nhấn")
	}
	if status&0x08 != 0 {
		errors = append(errors, "Hết giấy")
	}
	if status&0x10 != 0 {
		errors = append(errors, "Lỗi")
	}
	if status&0x20 != 0 {
		errors = append(errors, "Giấy gần hết")
	}
	if status&0x40 != 0 {
		errors = append(errors, "Lỗi cắt giấy")
	}
	if status&0x80 != 0 {
		errors = append(errors, "Lỗi không khắc phục")
	}

	if len(errors) > 0 {
		return fmt.Errorf("lỗi máy in: %v", errors)
	}

	return nil

}
