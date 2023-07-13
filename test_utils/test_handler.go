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

		response := &Response{NewMsg3: "Hello"}

		c.JSON(http.StatusOK, response)
	}
}

func (r *TestHandler) Versions() []*versioning.Version {
	return []*versioning.Version{
		{
			Date: "2018-03-11",
			Changes: []*versioning.Change{
				{
					ResponseAction: msgFieldChange3,
				},
			},
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
		{
			Date: "2018-01-11",
			Changes: []*versioning.Change{
				{
					ResponseAction: msgFieldChange2,
				},
			},
		},
	}
}

type Response struct {
	NewMsg3 string `json:"new_msg_3"`
}

func msgFieldChange(mapping map[string]interface{}) map[string]interface{} {
	mapping["new_msg"] = mapping["new_msg_2"]
	delete(mapping, "new_msg_2")
	return mapping
}

func msgFieldChange2(mapping map[string]interface{}) map[string]interface{} {
	mapping["msg"] = mapping["new_msg"]
	delete(mapping, "new_msg")
	return mapping
}

func msgFieldChange3(mapping map[string]interface{}) map[string]interface{} {
	mapping["new_msg_2"] = mapping["new_msg_3"]
	delete(mapping, "new_msg_3")
	return mapping
}

func dateDeleteChange(mapping map[string]interface{}) map[string]interface{} {
	delete(mapping, "date")
	return mapping
}
