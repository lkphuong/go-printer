package utils

import (
	"errors"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"testing"
)

// TestPrintImageCommand_ConvertAndSaveBitmap convert file anh nguon sang bitmap
// bang chinh ham printImageCommand dang duoc dung trong printFile (print.go),
// roi luu lai ket qua de kiem tra bang mat.
func TestPrintImageCommand_ConvertAndSaveBitmap(t *testing.T) {
	srcPath := filepath.Join("..", "..", "20260815214347_2de0e644b357d412_print_3bfc6da5-90bb-4889-845c-1d8eaff9b0a8_1204485130663952667.jpg")

	imgFile, err := os.Open(srcPath)
	if err != nil {
		t.Fatalf("khong mo duoc file anh nguon %s: %v", srcPath, err)
	}
	defer imgFile.Close()

	img, _, err := image.Decode(imgFile)
	if err != nil {
		t.Fatalf("giai ma anh that bai: %v", err)
	}

	const maxPrinterDots = 576 // giong MaxPrinterDots trong printFile (print.go)

	printCmd, err := printImageCommand(img, maxPrinterDots)
	if err != nil {
		t.Fatalf("printImageCommand loi: %v", err)
	}

	outDir := "testdata"
	if err := os.MkdirAll(outDir, 0755); err != nil {
		t.Fatalf("khong tao duoc thu muc %s: %v", outDir, err)
	}

	rawPath := filepath.Join(outDir, "print_command.bin")
	if err := os.WriteFile(rawPath, printCmd, 0644); err != nil {
		t.Fatalf("khong ghi duoc file lenh in: %v", err)
	}

	previewPath := filepath.Join(outDir, "print_preview.png")
	if err := saveBitmapPreview(printCmd, previewPath); err != nil {
		t.Fatalf("khong tao duoc anh preview: %v", err)
	}

	t.Logf("da luu lenh in (%d bytes) tai %s va anh preview tai %s", len(printCmd), rawPath, previewPath)
}

// saveBitmapPreview giai nguoc du lieu bitmap 1-bit trong lenh ESC/POS (GS v 0)
// thanh anh PNG den trang de xem lai bang mat.
func saveBitmapPreview(printCmd []byte, path string) error {
	const headerLen = 8 // 0x1D 0x76 0x30 0x00 + xL xH yL yH
	if len(printCmd) < headerLen {
		return errors.New("print command qua ngan, thieu header")
	}

	xL := int(printCmd[4])
	xH := int(printCmd[5])
	yL := int(printCmd[6])
	yH := int(printCmd[7])

	byteWidth := xL | (xH << 8)
	height := yL | (yH << 8)
	width := byteWidth * 8

	data := printCmd[headerLen:]
	if len(data) < byteWidth*height {
		return errors.New("du lieu bitmap khong du do dai da khai bao")
	}

	preview := image.NewGray(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			byteIndex := y*byteWidth + x/8
			bitPosition := 7 - (x % 8)
			bit := (data[byteIndex] >> bitPosition) & 1

			c := color.Gray{Y: 255}
			if bit == 1 {
				c = color.Gray{Y: 0}
			}
			preview.SetGray(x, y, c)
		}
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	return png.Encode(f, preview)
}
