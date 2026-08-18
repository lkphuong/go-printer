// Lệnh binpreview giải mã ngược file .bin (payload đã dump bởi dumpPrinterPayload
// trong internal/utils/print_spooler_windows.go — chính là byte thật gửi cho máy in
// qua WritePrinter) thành ảnh PNG, để so bằng mắt với ảnh gốc xem bitmap có khớp không.
//
// Cách dùng:
//
//	go run ./cmd/binpreview uploads/spool-dumps/20260818153000.000_iTP86.bin
//	go run ./cmd/binpreview uploads/spool-dumps/xxx.bin out.png
package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "cách dùng: binpreview <file.bin> [output.png]")
		os.Exit(1)
	}

	binPath := os.Args[1]
	outPath := ""
	if len(os.Args) >= 3 {
		outPath = os.Args[2]
	} else {
		ext := filepath.Ext(binPath)
		outPath = strings.TrimSuffix(binPath, ext) + ".png"
	}

	data, err := os.ReadFile(binPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "không đọc được file:", err)
		os.Exit(1)
	}

	img, err := decodeEscPosRasterImage(data)
	if err != nil {
		fmt.Fprintln(os.Stderr, "không giải mã được lệnh in:", err)
		os.Exit(1)
	}

	f, err := os.Create(outPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "không tạo được file ảnh:", err)
		os.Exit(1)
	}
	defer f.Close()

	if err := png.Encode(f, img); err != nil {
		fmt.Fprintln(os.Stderr, "không ghi được ảnh PNG:", err)
		os.Exit(1)
	}

	fmt.Printf("đã lưu ảnh preview tại %s (%dx%d)\n", outPath, img.Bounds().Dx(), img.Bounds().Dy())
}

// decodeEscPosRasterImage giải ngược lệnh ESC/POS "GS v 0" (raster bit image):
// 4 byte header 0x1D 0x76 0x30 0x00, kế tiếp 4 byte width/height (xL xH yL yH),
// rồi dữ liệu bitmap 1-bit. Payload dump có thể có thêm feed giấy + lệnh cắt phía
// sau — phần đó bị bỏ qua vì đã biết trước đúng số byte bitmap cần đọc.
func decodeEscPosRasterImage(data []byte) (image.Image, error) {
	const headerLen = 8
	if len(data) < headerLen {
		return nil, fmt.Errorf("dữ liệu quá ngắn, thiếu header GS v 0")
	}
	if data[0] != 0x1D || data[1] != 0x76 || data[2] != 0x30 {
		return nil, fmt.Errorf("không tìm thấy header lệnh in ảnh (GS v 0) ở đầu file")
	}

	xL, xH := int(data[4]), int(data[5])
	yL, yH := int(data[6]), int(data[7])

	byteWidth := xL | (xH << 8)
	height := yL | (yH << 8)
	width := byteWidth * 8

	bitmap := data[headerLen:]
	if len(bitmap) < byteWidth*height {
		return nil, fmt.Errorf("dữ liệu bitmap thiếu: cần %d byte, có %d byte", byteWidth*height, len(bitmap))
	}

	img := image.NewGray(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			byteIndex := y*byteWidth + x/8
			bitPosition := 7 - (x % 8)
			bit := (bitmap[byteIndex] >> bitPosition) & 1

			c := color.Gray{Y: 255}
			if bit == 1 {
				c = color.Gray{Y: 0}
			}
			img.SetGray(x, y, c)
		}
	}

	return img, nil
}
