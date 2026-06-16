package middlewares

import (
	"go-printer/internal/constants"
	"go-printer/internal/logger"
	"go-printer/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ValidateAPIKey() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("API-Key")
		if apiKey == "" {
			logger.LogPrint(constants.API_KEY_FAIL, http.StatusUnauthorized, "missing API key: "+c.Request.URL.Path)
			c.Abort()
			utils.ResponseError(c, http.StatusUnauthorized, constants.API_KEY_FAIL, nil)
			return
		}

		if apiKey != constants.API_KEY {
			logger.LogPrint(constants.API_KEY_FAIL, http.StatusForbidden, "invalid API key: "+c.Request.URL.Path)
			c.Abort()
			utils.ResponseError(c, http.StatusForbidden, constants.API_KEY_FAIL, nil)
			return
		}

		c.Next()
	}
}
