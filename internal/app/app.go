package app

import (
	"encoding/json"
	"fmt"
	"go-printer/internal/handlers"
	"go-printer/internal/logger"
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
	logger.LogPrint("Initializing folders...", 200, "Initializing folders...")
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

	// init MongoDB logger; on failure the service continues with file-only logging.
	if err := logger.Init(); err != nil {
		log.Println("MongoDB logger init failed, continuing with file logging only:", err)
		logger.LogPrint("MongoDB logger init failed, continuing with file logging only:", 500, fmt.Sprintf("MongoDB logger init failed: %v", err))
	}

	log.Println("Folders and logger initialized.")
	logger.LogPrint("Folders and logger initialized.", 200, "Folders and logger initialized.")
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
	logger.LogPrint(fmt.Sprintf("hour: %d", currentTime.Hour()), 200, fmt.Sprintf("hour: %d", currentTime.Hour()))
	if currentTime.Hour() <= 4 || currentTime.Hour() > 6 {
		logger.LogPrint("Not time for cleanup yet.", 200, "Not time for cleanup yet.")
		return
	}
	logger.LogPrint("Running cleanup of uploads folder...", 200, "Running cleanup of uploads folder...")
	uploadsDir := filepath.Join(".", "uploads")
	files, err := os.ReadDir(uploadsDir)
	if err != nil {
		logger.LogPrint(fmt.Sprintf("Error reading uploads directory: %v", err), 500, fmt.Sprintf("Error reading uploads directory: %v", err))
	} else {
		for _, file := range files {
			_ = os.Remove(filepath.Join(uploadsDir, file.Name()))
		}
		logger.LogPrint("Uploads cleanup completed.", 200, "Uploads cleanup completed.")
	}

	// cleanup logs older than 14 days
	logsDir := filepath.Join(".", "logs")
	files, err = os.ReadDir(logsDir)
	if err != nil {
		logger.LogPrint(fmt.Sprintf("Error reading logs directory: %v", err), 500, fmt.Sprintf("Error reading logs directory: %v", err))
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
			logger.LogPrint(fmt.Sprintf("Deleted old log file: %s", file.Name()), 200, fmt.Sprintf("Deleted old log file: %s", file.Name()))
		}
	}
	logger.LogPrint("Logs cleanup completed.", 200, "Logs cleanup completed.")

	// drop MongoDB log collections older than the retention window.
	logger.CleanupOldCollections()
}

func (a *App) Run() {
	logger.LogPrint("Server starting on :9099", 200, "Server starting on :9099")
	if a.router == nil {
		logger.LogPrint("router is nil, check NewApp() initialization", 500, "router is nil, check NewApp() initialization")
	}

	go func() {
		// interval to run cleanupFolder every day at 5:00 AM
		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()
		for {
			func() {
				defer func() {
					if rec := recover(); rec != nil {
						logger.LogPrint(fmt.Sprintf("cleanup goroutine panic recovered: %v", rec), 500, fmt.Sprintf("cleanup goroutine panic recovered: %v", rec))
					}
				}()
				cleanupFolder()
			}()
			<-ticker.C
		}
	}()

	a.router.Run(":9099")

}
