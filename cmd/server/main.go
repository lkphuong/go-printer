package main

import (
	"fmt"
	"go-printer/internal/app"
	"log"
	"net"
	"os"
	"path/filepath"

	"github.com/kardianos/service"
)

type program struct{}

func printIP() {
	interfaces, err := net.Interfaces()
	if err != nil {
		fmt.Println("Lỗi lấy interfaces:", err)
		return
	}

	for _, iface := range interfaces {
		if iface.Flags&net.FlagUp == 0 {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			ipNet, ok := addr.(*net.IPNet)
			if !ok {
				continue
			}

			ip := ipNet.IP

			if ip.To4() != nil && !ip.IsLoopback() {
				fmt.Printf("Interface: %s - IPv4: %s\n", iface.Name, ip.String())
			}
		}
	}
}

func (p *program) Start(s service.Service) error {
	go p.run()
	return nil
}

func (p *program) run() {
	exe, _ := os.Executable()
	dir := filepath.Dir(exe)
	os.Chdir(dir)

	logFile, _ := os.OpenFile("service.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	log.SetOutput(logFile)

	log.Println("Service started")

	a := app.NewApp()
	printIP()
	a.Run()
}

func (p *program) Stop(s service.Service) error {
	log.Println("Service stopped")
	return nil
}

func main() {
	svcConfig := &service.Config{
		Name:        "PrinterAMD64",
		DisplayName: "Printer AMD64 Service",
		Description: "Printer service",
	}

	prg := &program{}

	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)
	}

	// nếu chạy interactive (double click exe)
	if service.Interactive() {

		log.Println("Installing service...")

		err := s.Install()
		if err != nil {
			log.Println("Install failed:", err)
			return
		}

		log.Println("Starting service...")
		s.Start()

		return
	}

	err = s.Run()
	if err != nil {
		log.Fatal(err)
	}
}
