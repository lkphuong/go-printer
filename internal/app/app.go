package app

import (
	"encoding/json"
	"go-printer/internal/handlers"
	"go-printer/internal/middlewares"
	"go-printer/internal/routers"
	"go-printer/internal/services"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gopkg.in/natefinch/lumberjack.v2"
)

type App struct {
	services *Services
	handlers *Handlers
	router   *gin.Engine
}

type Services struct {
	PrintService *services.PrintService
}

type Handlers struct {
	PrintHandler *handlers.PrintHandler
}

func (a *App) setupServices() {
	a.services = &Services{
		PrintService: &services.PrintService{},
	}
}

func (a *App) setupHandlers() {
	a.handlers = &Handlers{
		PrintHandler: handlers.NewPrintHandler(a.services.PrintService),
	}
}

func (a *App) setupRouter() {
	config := cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:    []string{"Origin", "Content-Type", "Authorization", "X-API-Key"},
	}

	a.router.Use(cors.New(config))

	v1 := a.router.Group("/api/v1")
	v1.Use(middlewares.ValidateAPIKey())
	{
		routers.SetupPrintRoutes(v1, a.handlers.PrintHandler)
	}
}

func initializing() {
	// init folder database, uploads, logs
	log.Println("Initializing folders...")
	uploadsDir := filepath.Join(".", "uploads")
	os.MkdirAll(uploadsDir, 0755)

	configDir := filepath.Join(".", "config")
	os.MkdirAll(configDir, 0755)

	logsDir := filepath.Join(".", "logs")
	os.MkdirAll(logsDir, 0755)

	// init config.json
	configsFile := filepath.Join(configDir, "config.json")
	if _, err := os.Stat(configsFile); os.IsNotExist(err) {
		emptyConfigs := []interface{}{}
		data, _ := json.MarshalIndent(emptyConfigs, "", "  ")
		os.WriteFile(configsFile, data, 0644)
	}

	// init device.json
	deviceFile := filepath.Join(configDir, "device.json")
	if _, err := os.Stat(deviceFile); os.IsNotExist(err) {
		defaultDeviceConfig := map[string]string{
			"location": "office",
		}
		data, _ := json.MarshalIndent(defaultDeviceConfig, "", "  ")
		os.WriteFile(deviceFile, data, 0644)
	}

	// setup log file with lumberjack
	logFile := filepath.Join(logsDir, "app.log")
	log.SetOutput(&lumberjack.Logger{
		Filename:   logFile,
		MaxSize:    50, // megabytes before rotating
		MaxBackups: 0,  // unlimited backups
		MaxAge:     14, // days
		Compress:   false,
	})
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	log.Println("Folders and logger initialized.")
}

func NewApp() *App {

	initializing()

	app := &App{
		router: gin.Default(),
	}

	app.setupServices()
	app.setupHandlers()
	app.setupRouter()

	return app
}

func cleanupFolder() {
	// cleanup at 5:00 AM
	currentTime := time.Now()
	log.Println("hour: ", currentTime.Hour())
	if currentTime.Hour() <= 4 || currentTime.Hour() > 6 {
		log.Println("Not time for cleanup yet.")
		return
	}
	log.Println("Running cleanup of uploads folder...")
	uploadsDir := filepath.Join(".", "uploads")
	files, err := os.ReadDir(uploadsDir)
	if err != nil {
		log.Println("Error reading uploads directory:", err)
	} else {
		for _, file := range files {
			_ = os.Remove(filepath.Join(uploadsDir, file.Name()))
		}
		log.Println("Uploads cleanup completed.")
	}

	// cleanup logs older than 14 days
	logsDir := filepath.Join(".", "logs")
	files, err = os.ReadDir(logsDir)
	if err != nil {
		log.Println("Error reading logs directory:", err)
		return
	}
	now := time.Now()
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		info, err := file.Info()
		if err != nil {
			continue
		}
		if now.Sub(info.ModTime()) > 14*24*time.Hour {
			os.Remove(filepath.Join(logsDir, file.Name()))
			log.Printf("Deleted old log file: %s", file.Name())
		}
	}
	log.Println("Logs cleanup completed.")
}

func (a *App) Run() {
	log.Println("Server starting on :9099")
	if a.router == nil {
		log.Fatal("router is nil, check NewApp() initialization")
	}

	go func() {
		// interval to run cleanupFolder every day at 5:00 AM
		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()
		for {
			cleanupFolder()
			<-ticker.C
		}
	}()

	a.router.Run(":9099")

}
