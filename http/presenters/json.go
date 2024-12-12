package presenters

import (
	"github.com/go-chi/render"
	"net/http"
)

func RenderJSON(w http.ResponseWriter, r *http.Request, code int, v interface{}) {
	render.Status(r, code)
	render.JSON(w, r, v)
}
