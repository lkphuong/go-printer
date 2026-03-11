package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const serviceName = "PrinterService"
const serviceDisplayName = "Go Printer Service"
const serviceDescription = "Windows Service cho API Printer"

func runPowerShellCommand(script string) (string, error) {
	if runtime.GOOS != "windows" {
		return "", fmt.Errorf("Service only works on Windows")
	}

	cmd := exec.Command("powershell", "-NoProfile", "-Command", script)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

func isServiceExists() (bool, error) {
	script := fmt.Sprintf("$service = Get-Service -Name \"%s\" -ErrorAction SilentlyContinue\nif ($service) { \"true\" } else { \"false\" }", serviceName)

	output, err := runPowerShellCommand(script)
	if err != nil {
		return false, err
	}

	return strings.TrimSpace(output) == "true", nil
}

func createAndStartService(exePath string) error {
	// Escape path for PowerShell
	exePath = strings.ReplaceAll(exePath, `\`, `\\`)
	exePath = strings.ReplaceAll(exePath, `"`, `\"`)

	script := fmt.Sprintf("$ErrorActionPreference = 'Stop'\n\n"+
		"New-Service `\n"+
		"    -Name \"%s\" `\n"+
		"    -BinaryPathName \"%s\" `\n"+
		"    -DisplayName \"%s\" `\n"+
		"    -Description \"%s\" `\n"+
		"    -StartupType Automatic | Out-Null\n\n"+
		"Start-Service -Name \"%s\"\n\n"+
		"Write-Host \"Service created and started successfully\"",
		serviceName, exePath, serviceDisplayName, serviceDescription, serviceName)

	output, err := runPowerShellCommand(script)
	if err != nil {
		fmt.Printf("Error creating service: %s\n", output)
		return err
	}

	fmt.Println(output)
	return nil
}

func startService() error {
	script := fmt.Sprintf("Start-Service -Name \"%s\"\nWrite-Host \"Service started successfully\"", serviceName)

	output, err := runPowerShellCommand(script)
	if err != nil {
		fmt.Printf("Error: %s\n", output)
		return err
	}

	fmt.Println(output)
	return nil
}

func stopService() error {
	script := fmt.Sprintf("Stop-Service -Name \"%s\" -Force -ErrorAction SilentlyContinue\nWrite-Host \"Service stopped successfully\"", serviceName)

	output, err := runPowerShellCommand(script)
	if err == nil {
		fmt.Println(output)
	}
	return err
}

func removeService() error {
	script := fmt.Sprintf("Stop-Service -Name \"%s\" -Force -ErrorAction SilentlyContinue\nRemove-Service -Name \"%s\" -Force\nWrite-Host \"Service removed successfully\"", serviceName, serviceName)

	output, err := runPowerShellCommand(script)
	if err != nil {
		fmt.Printf("Error: %s\n", output)
		return err
	}

	fmt.Println(output)
	return nil
}

func getServiceStatus() string {
	script := fmt.Sprintf("$service = Get-Service -Name \"%s\" -ErrorAction SilentlyContinue\nif ($service) { $service.Status.ToString() } else { \"Not installed\" }", serviceName)

	output, err := runPowerShellCommand(script)
	if err != nil {
		return "Error getting status"
	}

	return strings.TrimSpace(output)
}

func printBanner() {
	fmt.Println("")
	fmt.Println("╔═══════════════════════════════════════════╗")
	fmt.Println("║   Go Printer - Windows Service Setup      ║")
	fmt.Println("╚═══════════════════════════════════════════╝")
	fmt.Println("")
}

func main() {
	printBanner()

	// Check if running on Windows
	if runtime.GOOS != "windows" {
		fmt.Println("Error: This tool only works on Windows")
		fmt.Print("Press Enter to exit...")
		fmt.Scanln()
		os.Exit(1)
	}

	// Check if running as Administrator
	// cmd := exec.Command("powershell", "-Command", `
	// 	if ([Security.Principal.WindowsPrincipal] [Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole] 'Administrator')) {
	// 		"true"
	// 	} else {
	// 		"false"
	// 	}
	// `)
	//output, err := cmd.Output()
	//isAdmin := err == nil && strings.TrimSpace(string(output)) == "true"

	// if !isAdmin {
	// 	fmt.Println("Requires Administrator privileges")
	// 	fmt.Println("")
	// 	fmt.Println("Please run this program as Administrator:")
	// 	fmt.Println("  Right-click on printer-service.exe")
	// 	fmt.Println("  Select 'Run as administrator'")
	// 	fmt.Println("")
	// 	fmt.Print("Press Enter to exit...")
	// 	fmt.Scanln()
	// 	os.Exit(1)
	// }

	exePath, err := os.Executable()
	if err != nil {
		fmt.Println("Error getting executable path:", err)
		fmt.Print("Press Enter to exit...")
		fmt.Scanln()
		os.Exit(1)
	}

	exeDir := filepath.Dir(exePath)

	// Parse command line arguments (optional)
	if len(os.Args) > 1 {
		cmd := strings.ToLower(os.Args[1])
		switch cmd {
		case "-start":
			if err := startService(); err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}
			fmt.Println("✓ Service started")
			return

		case "-stop":
			if err := stopService(); err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}
			fmt.Println("✓ Service stopped")
			return

		case "-status":
			status := getServiceStatus()
			fmt.Printf("Service status: %s\n", status)
			return

		case "-help", "-h", "-?":
			fmt.Println("Usage: printer-service.exe [command]")
			fmt.Println("Commands: -start, -stop, -status, -help")
			return
		}
	}

	// ===== DEFAULT: Auto-setup/replace mode (double-click) =====
	fmt.Println("Checking for existing service...")
	exists, err := isServiceExists()
	if err != nil {
		fmt.Printf("⚠ Warning: Could not check existing service: %v\n", err)
		exists = false
	}

	if exists {
		fmt.Println("✓ Service already exists")
		fmt.Println("")
		fmt.Println("Removing old service...")
		if err := removeService(); err != nil {
			fmt.Printf("⚠ Warning: %v\n", err)
		}
		time.Sleep(1 * time.Second)
	}

	fmt.Println("Installing service...")
	fmt.Printf("Path: %s\n", exePath)
	fmt.Println("")

	if err := createAndStartService(exePath); err != nil {
		fmt.Printf("Error: %v\n", err)
		fmt.Print("Press Enter to exit...")
		fmt.Scanln()
		os.Exit(1)
	}

	fmt.Println("")
	fmt.Println("Setup Complete!")
	fmt.Println("")
	fmt.Printf("✓ Service installed and started\n")
	fmt.Printf("✓ Access API: http://localhost:9099\n")
	fmt.Printf("✓ Logs: %s\\logs\\app.log\n", exeDir)
	fmt.Println("")
	fmt.Println("Closing in 4 seconds...")
	time.Sleep(4 * time.Second)
}
