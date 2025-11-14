package app

import (
	"encoding/json"
	"go-printer/internal/handlers"
	"go-printer/internal/routers"
	"go-printer/internal/services"
	"log"
	"os"
	"path/filepath"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
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
		AllowHeaders:    []string{"Origin", "Content-Type", "Authorization"},
	}

	a.router.Use(cors.New(config))

	v1 := a.router.Group("/api/v1")
	{
		routers.SetupPrintRoutes(v1, a.handlers.PrintHandler)
	}
}

func NewApp() *App {

	//init folder database, uploads
	log.Println("Initializing folders...")
	uploadsDir := filepath.Join(".", "uploads")
	os.MkdirAll(uploadsDir, 0755)

	configDir := filepath.Join(".", "config")
	os.MkdirAll(configDir, 0755)

	//init config.json

	configsFile := filepath.Join(configDir, "config.json")

	if _, err := os.Stat(configsFile); os.IsNotExist(err) {
		emptyConfigs := []interface{}{}
		data, _ := json.MarshalIndent(emptyConfigs, "", "  ")
		os.WriteFile(configsFile, data, 0644)
	}
	log.Println("Folders initialized.")

	app := &App{
		router: gin.Default(),
	}

	app.setupServices()
	app.setupHandlers()
	app.setupRouter()

	return app
}

func (a *App) Run() {
	log.Println("Server starting on :9099")
	if a.router == nil {
		log.Fatal("router is nil, check NewApp() initialization")
	}
	a.router.Run(":9099")
}
