package main

import (
	"go-printer/internal/app"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/kardianos/service"
)

type program struct{}

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

	// Supervisor loop: if NewApp() or Run() panics, recover, pause briefly, and
	// respawn the app so the printer service stays available without OS-level restart.
	for {
		superviseApp()
		log.Println("App exited, restarting in 2s...")
		time.Sleep(2 * time.Second)
	}
}

func superviseApp() {
	defer func() {
		if rec := recover(); rec != nil {
			log.Printf("app panic recovered by supervisor: %v", rec)
		}
	}()

	a := app.NewApp()
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
