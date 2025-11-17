package handlers

import (
	"go-printer/internal/constants"
	"go-printer/internal/dto/request"
	"go-printer/internal/services"
	"go-printer/internal/utils"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PrintHandler struct {
	printService *services.PrintService
}

func NewPrintHandler(printService *services.PrintService) *PrintHandler {
	return &PrintHandler{
		printService: printService,
	}
}

func (ph *PrintHandler) GetPrinters(c *gin.Context) {
	printers, err := ph.printService.GetPrintersLocal()
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "Failed to get printers", err.Error())
		return

	}

	utils.ResponseSuccess(c, printers, constants.OK)
}

func (ph *PrintHandler) GetPrintConfig(c *gin.Context) {
	printer := c.Param("printer")
	config, err := ph.printService.GetPrintConfig(printer)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "Failed to get print config", err.Error())
		return
	}

	if config.PrinterName == "" && printer != "" {
		utils.ResponseError(c, http.StatusBadRequest, "Printer config not found", nil)
		return
	}

	utils.ResponseSuccess(c, config, constants.OK)
}

func (ph *PrintHandler) ConfigPrinter(c *gin.Context) {
	var request request.PrintConfigRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Println("err", err)
		msg := utils.CustomErrorMessage(err)
		utils.ResponseError(c, http.StatusBadRequest, msg, err.Error())
		return
	}

	if err := ph.printService.ConfigPrinter(request.PrinterName, request.Type); err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "Failed to config printer", err.Error())
		return
	}

	utils.ResponseSuccess(c, nil, constants.OK)
}

func (ph *PrintHandler) JobPrint(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "Failed to get form data", err.Error())
		return
	}

	files := form.File["file"]
	printType := form.Value["type"]
	copies := form.Value["copies"]
	if len(copies) > 0 {
		log.Printf("Copies: %s\n", copies[0])
	}

	if len(files) == 0 {
		utils.ResponseError(c, http.StatusBadRequest, "No files uploaded", nil)
		return
	}

	if err := ph.printService.JobPrint(c, printType[0], copies[0], files[0]); err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "Failed to print job", err.Error())
		return
	}

	utils.ResponseSuccess(c, nil, constants.OK)
}

func (ph *PrintHandler) ClearCache(c *gin.Context) {
	if err := ph.printService.ClearCache(); err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "Failed to clear cache", err.Error())
		return
	}

	utils.ResponseSuccess(c, nil, constants.OK)
}
