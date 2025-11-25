package utils

import (
	"bytes"
	"embed"
	"fmt"
	"go-printer/internal/constants"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

//go:embed tools/SumatraPDF.exe
var sumatraPDF embed.FS
var sumatraPath string

func init() {
	log.Println("Initializing utils package")
	if runtime.GOOS == "windows" {
		tempDir := os.TempDir()
		tempPath := filepath.Join(tempDir, "SumatraPDF.exe")
		log.Printf("Temp dir: %s, Temp path: %s", tempDir, tempPath)
		if _, err := os.Stat(tempPath); os.IsNotExist(err) {
			// Chỉ extract nếu chưa tồn tại
			log.Println("Extracting SumatraPDF...")
			if err := extractEmbeddedFile("tools/SumatraPDF.exe", tempPath); err != nil {
				log.Printf("Failed to extract SumatraPDF: %v", err)
				return // Không gán sumatraPath nếu extract thất bại
			}
			log.Println("Extracted SumatraPDF successfully")
		} else {
			log.Println("SumatraPDF already exists in temp")
		}
		if _, err := os.Stat(tempPath); err == nil {
			sumatraPath = tempPath
			log.Printf("SumatraPath set to: %s", sumatraPath)
		} else {
			log.Printf("SumatraPDF file not found after extraction: %v", err)
		}
	}
}

func GetPrinters() ([]string, error) {
	switch runtime.GOOS {
	case "windows":
		cmd := exec.Command("powershell", "Get-Printer | Select-Object -ExpandProperty Name")
		var out bytes.Buffer
		cmd.Stdout = &out
		err := cmd.Run()
		if err != nil {
			return nil, err
		}
		line := strings.ReplaceAll(out.String(), "\r", "")
		lines := strings.Split(strings.TrimSpace(line), "\n")
		return lines, nil

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

func PrintFile(printer, filePath string, copies string) error {

	numCopies := 1
	if copies != "" {
		if _, err := fmt.Sscan(copies, &numCopies); err != nil || numCopies < 1 {
			numCopies = 1
		}
	}

	switch runtime.GOOS {
	case "windows":
		if sumatraPath == "" {
			// Fallback to mspaint if SumatraPDF not available
			log.Println("Using mspaint fallback")
			for i := 0; i < numCopies; i++ {
				psCmd := fmt.Sprintf("mspaint /pt %q %q", filePath, printer)
				cmd := exec.Command("powershell", "-NoProfile", "-Command", psCmd)
				var out bytes.Buffer
				cmd.Stdout = &out
				cmd.Stderr = &out
				if err := cmd.Run(); err != nil {
					fmt.Printf("print failed: %v: %s\n", err, out.String())
				}
			}
		} else {
			log.Printf("Using SumatraPDF at: %s", sumatraPath)
			for i := 0; i < numCopies; i++ {
				// use SumatraPDF for better performance
				cmd := exec.Command(sumatraPath, "-print-to", printer, "-silent", filePath)
				var out bytes.Buffer
				cmd.Stdout = &out
				cmd.Stderr = &out
				if err := cmd.Run(); err != nil {
					fmt.Printf("print failed: %v: %s\n", err, out.String())
				}
			}
		}

		return nil
	default:
		// Prefer lp, fall back to lpr
		for i := 0; i < numCopies; i++ {
			if _, err := exec.LookPath("lp"); err == nil {
				cmd := exec.Command("lp", "-d", printer, filePath)
				var out bytes.Buffer
				cmd.Stdout = &out
				cmd.Stderr = &out
				if err := cmd.Run(); err != nil {
					fmt.Printf("lp failed: %v: %s", err, out.String())
				}
			}
			if _, err := exec.LookPath("lpr"); err == nil {
				cmd := exec.Command("lpr", "-P", printer, filePath)
				var out bytes.Buffer
				cmd.Stdout = &out
				cmd.Stderr = &out
				if err := cmd.Run(); err != nil {
					fmt.Printf("lpr failed: %v: %s", err, out.String())
				}
			}
			fmt.Printf("no printing command found (lp or lpr)")
		}
		return nil
	}
}

func HealthCheckQueue(printer string) error {
	// interval 3 minutes, step 5s or queue empty return success
	for i := 0; i < 36; i++ {
		status, error := queuePrinter(printer)
		if status {
			break
		}

		if !status {
			if error != nil {
				return error
			}
			log.Printf("Health check printer %s passed\n", printer)
			time.Sleep(5 * time.Second)
		}
	}
	return nil
}

func queuePrinter(printer string) (bool, error) {
	// only windown
	psCmd := fmt.Sprintf("wmic printjob where \"name like '%%%s%%'\" list brief", printer)
	cmd := exec.Command("powershell", "-NoProfile", "-Command", psCmd)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		fmt.Printf("print failed: %v: %s\n", err, out.String())
	}

	jobStatus := out.String()

	// check queue empty
	if strings.TrimSpace(jobStatus) == "" {
		log.Printf("Print queue for printer %s is empty\n", printer)
		return true, nil
	}

	// check job status
	for _, line := range strings.Split(jobStatus, "\n") {
		if strings.TrimSpace(line) == "" || strings.HasPrefix(line, "JobId") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 3 {
			continue
		}
		jobID := fields[0]
		status := fields[2]
		log.Printf("Job ID: %s, Status: %s\n", jobID, status)

		// check status Available -> success
		if strings.ToLower(status) == "available" {
			return true, nil
		}

		if strings.ToLower(status) == "error" {

			// clean all job
			psCleanCmd := fmt.Sprintf("wmic printjob where \"name like '%%%s%%'\" delete", printer)
			cleanCmd := exec.Command("powershell", "-NoProfile", "-Command", psCleanCmd)
			var cleanOut bytes.Buffer
			cleanCmd.Stdout = &cleanOut
			cleanCmd.Stderr = &cleanOut
			if err := cleanCmd.Run(); err != nil {
				fmt.Printf("clean print job failed: %v: %s\n", err, cleanOut.String())
			} else {
				log.Printf("Cleaned all print jobs for printer %s\n", printer)
			}

			return false, fmt.Errorf(constants.PRINT_FAILED)
		}
	}
	return false, nil
}

func extractEmbeddedFile(embedPath, destPath string) error {
	data, err := sumatraPDF.ReadFile(embedPath)
	if err != nil {
		return err
	}
	file, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = io.Copy(file, bytes.NewReader(data))
	if err != nil {
		return err
	}
	return os.Chmod(destPath, 0755) // accept executable
}
