package health

import (
	"fmt"
	"github.com/bloock/go-kit/http/presenters"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

type ExternalService interface {
	HealthCheck() ExternalServiceDetails
}

type Health struct {
	services    []ExternalService
	output      string
	id          string
	notes       []string
	links       []string
	description string
}

func NewHealth(o, i, d string, n, l []string, e []ExternalService) Health {
	return Health{
		services:    e,
		output:      o,
		id:          i,
		notes:       n,
		links:       l,
		description: d,
	}
}

func (h *Health) AddService(s ExternalService) {
	h.services = append(h.services, s)
}

// swagger:response Health
type BodyProductListResponse struct {
	// in: body
	Body HealthResponse
}

type ExternalServiceDetails struct {
	// Required: true
	// description: human-friendly description of the service.
	Description string `json:"description,omitempty"`
	// description: public version of the service.
	Version string `json:"version,omitempty"`
	// Required: true
	// description: indicates whether the service status is acceptable or not. API publishers
	Status string `json:"status,omitempty"`
	// description: error msg
	Error string `json:"error,omitempty"`
}

type HealthResponse struct {
	// Required: true
	// description: indicates whether the service status is acceptable or not. API publishers
	Status string `json:"status"`
	// description: public version of the service.
	Version string `json:"version,omitempty"`
	// description:  in well-designed APIs, backwards-compatible changes in the service should not update a version number.
	RelaseID string `json:"relaseID,omitempty"`
	// description: array of notes relevant to current state of health
	Notes []string `json:"notes,omitempty"`
	// description: raw error output, in case of “fail” or “warn” states. This field SHOULD be omitted for “pass” state.
	Output string `json:"output,omitempty"`
	// description: an object representing status of sub-components of the service in question
	Details []ExternalServiceDetails `json:"details,omitempty"`
	// description: an array of objects containing link relations and URIs
	Links []string `json:"links,omitempty"`
	// description: unique identifier of the service, in the application scope
	ServiceID string `json:"serviceID,omitempty"`
	// Required: true
	// description: human-friendly description of the service.
	Description string `json:"description"`
}

// swagger:route GET /health infrastructure Health
//
//	Produces:
//	- application/health+json
//
//	Schemes: http, https
//
//	Responses:
//		201: Health
//		400:
//		500:
func (h Health) CheckGinHandler() gin.HandlerFunc {

	return func(ctx *gin.Context) {
		status, resp := h.getHealth()
		ctx.JSON(status, resp)
	}
}

// swagger:route GET /health infrastructure Health
//
//	Produces:
//	- application/health+json
//
//	Schemes: http, https
//
//	Responses:
//		201: Health
//		400:
//		500:
func (h Health) CheckChiHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		status, resp := h.getHealth()
		presenters.RenderJSON(w, r, status, resp)
	}
}

func (h Health) getHealth() (int, interface{}) {
	var details []ExternalServiceDetails
	var err []string

	serviceDown := false
	for _, service := range h.services {
		detail := service.HealthCheck()
		details = append(details, detail)
		if detail.Error != "" {
			serviceDown = true
			err = append(err, detail.Error)
		}
	}

	var output string
	if len(err) > 0 {
		output = fmt.Sprintf("%s %s", h.output, strings.Join(err[:], ","))
		return http.StatusServiceUnavailable, HealthResponse{}
	}

	version := os.Getenv("IMAGE")
	release := os.Getenv("IMAGE")

	status := "pass"
	if serviceDown {
		status = "unestable"
	}

	health := HealthResponse{
		Status:      status,
		Version:     version,
		RelaseID:    release,
		Notes:       h.notes,
		Output:      output,
		Details:     details,
		Links:       h.links,
		ServiceID:   h.id,
		Description: h.description,
	}
	return http.StatusOK, health
}
