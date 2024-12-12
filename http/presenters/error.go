package presenters

import (
	"github.com/bloock/go-kit/errors"
	"github.com/go-chi/render"
	"net/http"
)

func RenderError(w http.ResponseWriter, r *http.Request, err error) {
	var parsedError errors.HttpAppError

	switch err.(type) {
	case errors.HttpAppError:
		parsedError = err.(errors.HttpAppError)
	default:
		parsedError = errors.NewHttpAppError(http.StatusInternalServerError, "Internal Server Error")
	}

	render.Status(r, parsedError.Code)
	render.JSON(w, r, parsedError.Error())
}
