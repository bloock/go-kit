package health

import (
	"encoding/json"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type HealthHttp struct {
	path        string
	description string
	version     string
}

func NewHealthHttp(path, description, version string) HealthHttp {
	return HealthHttp{
		path:        path,
		description: description,
		version:     version,
	}
}

func (h HealthHttp) HealthCheck() ExternalServiceDetails {
	client := &http.Client{}

	response, err := client.Get(h.path)

	if err != nil || response.StatusCode != 200 {
		var m string
		if err == nil {
			m = fmt.Sprintf("Unknown error with status %d path %s", response.StatusCode, h.path)
		} else {
			m = err.Error()
		}
		return ExternalServiceDetails{
			Description: h.description,
			Version:     h.version,
			Status:      "error",
			Error:       m,
		}
	}
	defer response.Body.Close()

	decoder := json.NewDecoder(response.Body)
	var data HealthResponse
	err = decoder.Decode(&data)

	if err != nil {
		return ExternalServiceDetails{
			Description: h.description,
			Version:     h.version,
			Status:      "error",
			Error:       err.Error(),
		}
	}

	description := data.Description
	if data.Description == "" {
		description = h.description
	}

	version := data.ServiceID
	if data.Version == "" {
		version = h.version
	}

	status := data.Status
	if data.Status == "" {
		status = "pass"
	}

	return ExternalServiceDetails{
		Description: description,
		Version:     version,
		Status:      status,
		Error:       data.Output,
	}
}
