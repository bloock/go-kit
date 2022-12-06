package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrorMiddleware() gin.HandlerFunc {
	return errorMiddleware(gin.ErrorTypeAny)
}

func errorMiddleware(errType gin.ErrorType) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		detectedErrors := c.Errors.ByType(errType)

		if len(detectedErrors) > 0 {
			err := detectedErrors[0].Err
			var parsedError *HttpAppError

			switch err.(type) {
			case HttpAppError:
				parsedError = err.(*HttpAppError)
			default:

				parsedError = NewHttpAppError(http.StatusInternalServerError, "Internal Server Error")
			}

			c.AbortWithStatusJSON(parsedError.Code, parsedError)
			return
		}

	}
}
