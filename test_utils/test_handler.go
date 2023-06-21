package test_utils

import (
	"github.com/bloock/go-kit/http/versioning"
	"github.com/bloock/go-kit/observability"
	"github.com/gin-gonic/gin"
	"net/http"
)

type TestHandler struct{}

var TestHandlerInstance = &TestHandler{}

func (r *TestHandler) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		s, ctx := observability.NewSpan(c, "service.repository.action")
		defer s.Finish()

		l := observability.InitLogger("local", "test_service", "1.0.0", true)
		l.Debug(ctx).Msg("a gin message")

		response := &Response{Msg: "Hello"}

		c.JSON(http.StatusOK, response)
	}
}

func (r *TestHandler) Versions() []*versioning.Version {
	return []*versioning.Version{
		{
			Date: "2018-02-10",
		},
		{
			Date: "2018-02-11",
			Changes: []*versioning.Change{
				{
					RequestAction:  msgFieldChange,
					ResponseAction: msgFieldChange,
				},
			},
		},
	}
}

type Response struct {
	Msg string
}

func msgFieldChange(mapping map[string]interface{}) map[string]interface{} {
	mapping["new_msg"] = mapping["Msg"]
	delete(mapping, "Msg")
	return mapping
}
func dateDeleteChange(mapping map[string]interface{}) map[string]interface{} {
	delete(mapping, "date")
	return mapping
}
