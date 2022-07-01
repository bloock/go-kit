package httperror

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e AppError) Error() string {
	return fmt.Sprintf("error: %s. code: %d", e.Message, e.Code)
}

func NewAppError(code int, message string) *AppError {
	return &AppError{
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
			var parsedError *AppError

			switch err.(type) {
			case *AppError:
				parsedError = err.(*AppError)
			default:
				parsedError = &AppError{
					Code:    http.StatusInternalServerError,
					Message: "Internal Server Error",
				}
			}

			c.AbortWithStatusJSON(parsedError.Code, parsedError)
			return
		}

	}
}
