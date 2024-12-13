package presenters

import (
	"context"
	"github.com/go-chi/render"
	"net/http"
)

func RenderData(w http.ResponseWriter, r *http.Request, contentType string, code int, b []byte) {
	render.Status(r, code)
	*r = *r.WithContext(context.WithValue(r.Context(), render.ContentTypeCtxKey, contentType))
	render.Data(w, r, b)
}
