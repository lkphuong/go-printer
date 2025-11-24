package constants

const (
	OK                         = "ok"
	PrintTypeKitchen PrintType = "kitchen"
	PrintTypeCashier PrintType = "cashier"
)

type PrintType string

const (
	PRINT_FAILED    = "print_failed"
	PRINT_NOT_FOUND = "print_not_found"
)
