package utils

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

func CustomErrorMessage(err error) string {
	customMessages := map[string]string{
		"required": "This field is required.",
		"min":      "This field must be at least %s characters.",
		"max":      "This field must be at most %s characters.",
		"datetime": "Invalid date format. Use YYYY-MM-DD.",
	}

	errors := make(map[string]string)
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, fieldError := range validationErrors {
			field := strings.ToLower(fieldError.Field())
			tag := fieldError.Tag()
			if msg, exists := customMessages[tag]; exists {
				if tag == "min" || tag == "max" {
					msg = fmt.Sprintf(msg, fieldError.Param())
				}
				errors[field] = msg
			} else {
				errors[field] = fieldError.Error()
			}
		}
	}

	var msg string
	for k, v := range errors {
		msg = k + ": " + v
		break
	}
	return msg
}
