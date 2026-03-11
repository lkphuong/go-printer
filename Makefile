run:
	go run ./cmd/server/main.go

build-windows:
	GOOS=windows GOARCH=amd64 go build -o build/printer-amd64.exe ./cmd/server/main.go

build-windows-service:
	GOOS=windows GOARCH=amd64 go build -o build/printer-service-amd64.exe ./cmd/service/main.go

build-mac:
	GOOS=darwin GOARCH=arm64 go build -o build/printer-arm64 ./cmd/server/main.go

# Windows Service Setup
# 1. Build: make build-windows-service
# 2. Run: Double-click build/printer-service.exe (as Administrator)
# 3. That's it! Service will auto-create via PowerShell New-Service