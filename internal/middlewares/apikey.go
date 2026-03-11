package middlewares

import (
	"go-printer/internal/constants"
	"go-printer/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ValidateAPIKey() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("api-key")
		if apiKey == "" {
			c.Abort()
			utils.ResponseError(c, http.StatusUnauthorized, constants.API_KEY_FAIL, nil)
			return
		}

		if apiKey != constants.API_KEY {
			c.Abort()
			utils.ResponseError(c, http.StatusForbidden, constants.API_KEY_FAIL, nil)
			return
		}

		c.Next()
	}
}
