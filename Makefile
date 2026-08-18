run:
	go run ./cmd/server/main.go

build-windows:
	GOOS=windows GOARCH=amd64 go build -o build/printer-amd64.exe ./cmd/server/main.go

build-mac:
	GOOS=darwin GOARCH=arm64 go build -o build/printer-arm64 ./cmd/server/main.go

build-linux:
	GOOS=linux GOARCH=amd64 go build -o build/printer-amd64 ./cmd/server/main.go

docker-build:
	docker build -t go-printer .

bundle-p9099:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build \
		-trimpath \
		-ldflags="-s -w -buildid=" \
		-o build/printer-p9099-amd64.exe \
		./cmd/server/main.go

bundle-p9098:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build \
		-trimpath \
		-ldflags="-s -w -buildid=" \
		-o build/printer-p9098-amd64.exe \
		./cmd/server/main.go

# Windows Service Setup
# 1. Build: make build-windows-service
# 2. Run: Double-click build/printer-service.exe (as Administrator)
# 3. That's it! Service will auto-create via PowerShell New-Service