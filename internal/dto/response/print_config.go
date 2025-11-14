package response

type PrintConfigResponse struct {
	PrinterName string   `json:"printer_name"`
	Type        []string `json:"type"`
}
