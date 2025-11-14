package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Success  bool        `json:"success"`
	Message  string      `json:"message"`
	Data     interface{} `json:"data,omitempty"`
	MetaData interface{} `json:"metadata,omitempty"`
}

func ResponseSuccess(c *gin.Context, data interface{}, message string, metaData ...interface{}) {
	c.JSON(http.StatusOK, Response{
		Success:  true,
		Message:  message,
		Data:     data,
		MetaData: metaData,
	})
}

func ResponseError(c *gin.Context, statusCode int, message string, err interface{}) {
	c.JSON(statusCode, Response{
		Success: false,
		Message: message,
	})
}
