package handlers

import (
	"go-printer/internal/constants"
	"go-printer/internal/dto/request"
	"go-printer/internal/logger"
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
		utils.ResponseError(c, http.StatusBadRequest, constants.GET_PRINTERS_FAILED, err.Error())
		return

	}

	utils.ResponseSuccess(c, printers, constants.OK)
}

func (ph *PrintHandler) GetPrintConfig(c *gin.Context) {
	printer := c.Param("printer")
	config, err := ph.printService.GetPrintConfig(printer)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, constants.READ_CONFIG_FAILED, err.Error())
		return
	}

	if config.PrinterName == "" && printer != "" {
		logger.LogPrint(constants.CONFIG_NOT_FOUND, http.StatusBadRequest, "printer config not found: "+printer)
		utils.ResponseError(c, http.StatusBadRequest, constants.CONFIG_NOT_FOUND, nil)
		return
	}

	utils.ResponseSuccess(c, config, constants.OK)
}

func (ph *PrintHandler) ConfigPrinter(c *gin.Context) {
	var request request.PrintConfigRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Println("err", err)
		msg := utils.CustomErrorMessage(err)
		logger.LogPrint(constants.INVALID_REQUEST, http.StatusBadRequest, err.Error())
		utils.ResponseError(c, http.StatusBadRequest, msg, err.Error())
		return
	}

	if err := ph.printService.ConfigPrinter(request.PrinterName, request.Type); err != nil {
		utils.ResponseError(c, http.StatusBadRequest, constants.WRITE_CONFIG_FAILED, err.Error())
		return
	}

	utils.ResponseSuccess(c, nil, constants.OK)
}

func (ph *PrintHandler) JobPrint(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		logger.LogPrint(constants.FORM_PARSE_FAILED, http.StatusBadRequest, err.Error())
		utils.ResponseError(c, http.StatusBadRequest, constants.FORM_PARSE_FAILED, err.Error())
		return
	}

	files := form.File["file"]
	printer := form.Value["printer"]
	copies := form.Value["copies"]
	log.Println("copies: ", len(copies))
	if len(copies) == 0 {
		copies = []string{"1"} // default 1 copy
	}

	if len(files) == 0 {
		logger.LogPrint(constants.NO_FILES_UPLOADED, http.StatusBadRequest, "no files in form")
		utils.ResponseError(c, http.StatusBadRequest, constants.NO_FILES_UPLOADED, nil)
		return
	}

	if len(printer) == 0 {
		logger.LogPrint(constants.NO_PRINTER, http.StatusBadRequest, "printer field missing")
		utils.ResponseError(c, http.StatusBadRequest, constants.NO_PRINTER, nil)
		return
	}

	if err := ph.printService.JobPrint(c, printer[0], copies[0], files[0]); err != nil {
		utils.ResponseError(c, http.StatusBadRequest, err.Error(), err.Error())
		return
	}

	utils.ResponseSuccess(c, nil, constants.OK)
}

func (ph *PrintHandler) ClearCache(c *gin.Context) {
	if err := ph.printService.ClearCache(); err != nil {
		utils.ResponseError(c, http.StatusBadRequest, constants.CLEAR_CACHE_FAILED, err.Error())
		return
	}

	utils.ResponseSuccess(c, nil, constants.OK)
}
