package routers

import (
	"go-printer/internal/handlers"

	"github.com/gin-gonic/gin"
)

func SetupPrintRoutes(router *gin.RouterGroup, handler *handlers.PrintHandler) {
	printers := router.Group("/printers")
	{
		printers.GET("", handler.GetPrinters)
		printers.GET("/:printer/config", handler.GetPrintConfig)
		printers.POST("/config", handler.ConfigPrinter)
		printers.POST("/jobs", handler.JobPrint)
		printers.DELETE("/cache", handler.ClearCache)
	}
}
