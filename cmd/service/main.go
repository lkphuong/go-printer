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

func printHelp() {
	fmt.Println("Usage:")
	fmt.Println("  printer-service.exe [command]")
	fmt.Println("")
	fmt.Println("Commands:")
	fmt.Println("  (no args)     - Create and start service (interactive)")
	fmt.Println("  -install      - Create service")
	fmt.Println("  -start        - Start service")
	fmt.Println("  -stop         - Stop service")
	fmt.Println("  -remove       - Remove service")
	fmt.Println("  -status       - Check service status")
	fmt.Println("  -help         - Show this help")
	fmt.Println("")
}

func main() {
	// Check if running on Windows
	if runtime.GOOS != "windows" {
		fmt.Println("Error: This tool only works on Windows")
		os.Exit(1)
	}

	// Check if running as Administrator
	cmd := exec.Command("powershell", "-Command", `
		if ([Security.Principal.WindowsPrincipal] [Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole] 'Administrator')) {
			"true"
		} else {
			"false"
		}
	`)
	output, err := cmd.Output()
	isAdmin := err == nil && strings.TrimSpace(string(output)) == "true"

	if !isAdmin {
		fmt.Println("Requires Administrator privileges")
		fmt.Println("Please run this program as Administrator:")
		fmt.Println("  1. Right-click on printer-service.exe")
		fmt.Println("  2. Select 'Run as administrator'")
		fmt.Println("")
		fmt.Print("Press Enter to exit...")
		fmt.Scanln()
		os.Exit(1)
	}

	exePath, err := os.Executable()
	if err != nil {
		fmt.Println("Error getting executable path:", err)
		os.Exit(1)
	}

	exeDir := filepath.Dir(exePath)

	// Parse command line arguments
	if len(os.Args) > 1 {
		cmd := strings.ToLower(os.Args[1])

		switch cmd {
		case "-install":
			fmt.Println("Installing service...")
			if err := createAndStartService(exePath); err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}
			fmt.Println("✓ Service installed and started")

		case "-start":
			if err := startService(); err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}

		case "-stop":
			if err := stopService(); err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}

		case "-remove":
			if err := removeService(); err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}

		case "-status":
			status := getServiceStatus()
			fmt.Printf("Service status: %s\n", status)

		case "-help", "-h", "-?":
			printHelp()

		default:
			fmt.Printf("Unknown command: %s\n", cmd)
			printHelp()
			os.Exit(1)
		}
		return
	}

	// No arguments - interactive mode

	// Check if service exists
	exists, err := isServiceExists()
	if err != nil {
		fmt.Printf("Error checking service: %v\n", err)
		os.Exit(1)
	}

	if exists {
		status := getServiceStatus()
		fmt.Printf("Service already exists (Status: %s)\n", status)
		fmt.Println("")
		fmt.Println("Choose action:")
		fmt.Println("  1. Start service")
		fmt.Println("  2. Stop service")
		fmt.Println("  3. Reinstall service")
		fmt.Println("  4. Remove service")
		fmt.Println("  5. Exit")
		fmt.Print("\nEnter choice (1-5): ")

		var choice string
		fmt.Scanln(&choice)

		switch choice {
		case "1":
			fmt.Println("\nStarting service...")
			if err := startService(); err != nil {
				fmt.Printf("Error: %v\n", err)
			} else {
				fmt.Println("✓ Service started")
			}

		case "2":
			fmt.Println("\nStopping service...")
			if err := stopService(); err != nil {
				fmt.Printf("Error: %v\n", err)
			} else {
				fmt.Println("✓ Service stopped")
			}

		case "3":
			fmt.Println("\nReinstalling service...")
			if err := removeService(); err != nil {
				fmt.Printf("Error: %v\n", err)
				time.Sleep(2 * time.Second)
			}
			fmt.Println("\nCreating new service...")
			if err := createAndStartService(exePath); err != nil {
				fmt.Printf("Error: %v\n", err)
			} else {
				fmt.Println("✓ Service reinstalled and started")
			}

		case "4":
			fmt.Println("\nRemoving service...")
			if err := removeService(); err != nil {
				fmt.Printf("Error: %v\n", err)
			} else {
				fmt.Println("✓ Service removed")
			}

		case "5":
			fmt.Println("Exiting...")
			return

		default:
			fmt.Println("Invalid choice")
		}
	} else {
		// Service doesn't exist, create and start it
		fmt.Println("Service not found. Creating...")
		fmt.Printf("Executable: %s\n", exePath)
		fmt.Printf("Working directory: %s\n", exeDir)
		fmt.Println("")

		if err := createAndStartService(exePath); err != nil {
			fmt.Printf("Error: %v\n", err)
			fmt.Print("\nPress Enter to exit...")
			fmt.Scanln()
			os.Exit(1)
		}

		fmt.Println("")
		fmt.Println("Setup Complete!")
		fmt.Println("")
		fmt.Printf("Access API at: http://localhost:9099\n")
		fmt.Printf("Logs at: %s\\logs\\app.log\n", exeDir)
		fmt.Println("")
		fmt.Println("The window will close in 5 seconds...")
		time.Sleep(5 * time.Second)
	}
}
