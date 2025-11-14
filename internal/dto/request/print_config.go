package request

type PrintConfigRequest struct {
	PrinterName string   `json:"printer_name" binding:"required"`
	Type        []string `json:"type" binding:"required"`
}
