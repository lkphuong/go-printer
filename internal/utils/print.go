package utils

import (
	"bytes"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

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
		// Use mspaint with /pt option to print

		go func() {
			{
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
			}
		}()
		return nil
	default:
		// Prefer lp, fall back to lpr
		go func() {
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
		}()
		return nil
	}
}
