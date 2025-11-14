package app

import (
	"go-printer/internal/handlers"
	"go-printer/internal/routers"
	"go-printer/internal/services"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type App struct {
	services *Services
	handlers *Handlers
	// mongoDB  *database.MongoDB
	router *gin.Engine
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
