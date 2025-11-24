package utils

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Status    bool        `json:"status"`
	ErrorCode string      `json:"errorCode,omitempty"`
	Data      interface{} `json:"data,omitempty"`
}

func ResponseSuccess(c *gin.Context, data interface{}, errorCode string, metaData ...interface{}) {
	c.JSON(http.StatusOK, Response{
		Status:    true,
		ErrorCode: "",
		Data:      data,
	})
}

func ResponseError(c *gin.Context, statusCode int, errorCode string, err interface{}) {
	log.Println("error message: ", errorCode)
	c.JSON(statusCode, Response{
		Status:    false,
		ErrorCode: errorCode,
	})
}
