run:
	go run ./cmd/server/main.go

build-windows:
	GOOS=windows GOARCH=amd64 go build -o build/printer.exe ./cmd/server/main.go

build-mac:
	GOOS=darwin GOARCH=arm64 go build -o build/printer ./cmd/server/main.go