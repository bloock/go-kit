package http

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type appError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e appError) Error() string {
	return fmt.Sprintf("error: %s. code: %d", e.Message, e.Code)
}

func NewAppError(code int, message string) appError {
	return appError{
		Code:    code,
		Message: message,
	}
}

func ErrorMiddleware() gin.HandlerFunc {
	return errorMiddleware(gin.ErrorTypeAny)
}

func errorMiddleware(errType gin.ErrorType) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		detectedErrors := c.Errors.ByType(errType)

		if len(detectedErrors) > 0 {
			err := detectedErrors[0].Err
			var parsedError *appError

			switch err.(type) {
			case *appError:
				parsedError = err.(*appError)
			default:
				parsedError = &appError{
					Code:    http.StatusInternalServerError,
					Message: "Internal Server Error",
				}
			}

			c.AbortWithStatusJSON(parsedError.Code, parsedError)
			return
		}

	}
}
