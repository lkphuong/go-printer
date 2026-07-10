package services

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"go-printer/internal/constants"
	"go-printer/internal/dto/response"
	"go-printer/internal/logger"
	"go-printer/internal/utils"
	"log"
	"mime/multipart"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type PrintService struct{}

func newTempUploadPath(now time.Time, filename string) (string, error) {
	randomSuffix := make([]byte, 8)
	if _, err := rand.Read(randomSuffix); err != nil {
		return "", err
	}

	baseName := filepath.Base(filename)
	if baseName == "." || baseName == string(filepath.Separator) {
		baseName = "upload"
	}

	return filepath.Join(
		"uploads",
		now.Format("20060102150405")+"_"+hex.EncodeToString(randomSuffix)+"_"+baseName,
	), nil
}

func (ps *PrintService) GetPrintersLocal() ([]string, error) {

	printers, err := utils.GetPrinters()
	if err != nil {
		logger.LogPrint(constants.GET_PRINTERS_FAILED, 400, err.Error())
		return nil, err
	}

	return printers, nil
}

func (ps *PrintService) GetPrintConfig(printer string) (response.PrintConfigResponse, error) {
	// Đọc file json
	data, err := os.ReadFile("./config/config.json")
	if err != nil {
		log.Println("Error reading config file:", err)
		logger.LogPrint(constants.READ_CONFIG_FAILED, 400, err.Error())
		return response.PrintConfigResponse{}, err
	}

	var config []response.PrintConfigResponse
	if err := json.Unmarshal(data, &config); err != nil {
		log.Println("Error unmarshalling config:", err)
		logger.LogPrint(constants.READ_CONFIG_FAILED, 400, err.Error())
		return response.PrintConfigResponse{}, err
	}

	// Lọc theo printer nếu có
	if printer != "" {
		filteredConfig := response.PrintConfigResponse{}
		for _, c := range config {
			if c.PrinterName == printer {
				filteredConfig = c
				break
			}
		}
		return filteredConfig, nil
	}

	return response.PrintConfigResponse{}, nil
}

func (ps *PrintService) ConfigPrinter(printerName string, types []string) error {

	// Đọc file json
	data, err := os.ReadFile("./config/config.json")
	if err != nil {
		log.Println("Error reading config file:", err)
		logger.LogPrint(constants.READ_CONFIG_FAILED, 400, err.Error())
		return err
	}

	var config []response.PrintConfigResponse
	if err := json.Unmarshal(data, &config); err != nil {
		log.Println("Error unmarshalling config:", err)
		logger.LogPrint(constants.READ_CONFIG_FAILED, 400, err.Error())
		return err
	}

	// Xoá config với type
	newConfig := []response.PrintConfigResponse{}
	for _, c := range config {
		keep := true
		for _, t := range types {
			if c.PrinterName == printerName {
				for _, existingType := range c.Type {
					if existingType == t {
						keep = false
						break
					}
				}
			}
		}
		if keep {
			newConfig = append(newConfig, c)
		}
	}

	// Thêm config mới
	newConfig = append(newConfig, response.PrintConfigResponse{
		PrinterName: printerName,
		Type:        types,
	})

	// Ghi lại file json
	newData, err := json.MarshalIndent(newConfig, "", "  ")
	if err != nil {
		log.Println("Error marshalling config:", err)
		logger.LogPrint(constants.WRITE_CONFIG_FAILED, 400, err.Error())
		return err
	}

	if err := os.WriteFile("./config/config.json", newData, 0644); err != nil {
		log.Println("Error writing config file:", err)
		logger.LogPrint(constants.WRITE_CONFIG_FAILED, 400, err.Error())
		return err
	}

	return nil
}

func (ps *PrintService) JobPrint(c *gin.Context, printer string, copies string, file *multipart.FileHeader) error {
	// lấy những printer theo type
	cmd := exec.Command(
		"powershell",
		fmt.Sprintf(
			`(Get-PrinterPort -Name (Get-Printer -Name "%s").PortName).PrinterHostAddress`,
			printer,
		),
	)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		logger.LogPrint(constants.PRINT_NOT_FOUND, 400, err.Error())
		return fmt.Errorf(constants.PRINT_NOT_FOUND)
	}

	ip := strings.TrimSpace(out.String())

	// save file tạm thời
	tempFilePath, err := newTempUploadPath(time.Now(), file.Filename)
	if err != nil {
		logger.LogPrint(constants.SAVE_FILE_FAILED, 400, err.Error())
		return err
	}
	if err := c.SaveUploadedFile(file, tempFilePath); err != nil {
		logger.LogPrint(constants.SAVE_FILE_FAILED, 400, err.Error())
		return err
	}

	f, err := file.Open()
	if err != nil {
		log.Println("read file error: ", err)
		logger.LogPrint(constants.READ_FILE_FAILED, 400, err.Error())
		return err
	}
	defer f.Close()

	// ping printer
	// gửi in cho từng printer
	if err := utils.PrintFileQueued(ip, tempFilePath, copies); err != nil {
		log.Println("print file error: ", err)
		logger.LogPrint(constants.QUEUE_FULL, 400, err.Error())
		return err
	}

	// job đã được đưa vào hàng đợi; kết quả in cuối cùng được worker ghi log sau.
	logger.LogPrint(constants.QUEUED, 200, "")
	log.Println("print queued")
	return nil
}

func (ps *PrintService) ClearCache() error {
	// xoá hết file json dữ liệu trong file json
	if err := os.WriteFile("./config/config.json", []byte("[]"), 0644); err != nil {
		log.Println("Error clearing cache:", err)
		logger.LogPrint(constants.CLEAR_CACHE_FAILED, 400, err.Error())
		return err
	}
	return nil
}
