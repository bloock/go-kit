package health

import (
	"encoding/json"
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

	if err != nil || response.StatusCode == 200 {
		return ExternalServiceDetails{
			Description: h.description,
			Version:     h.version,
			Status:      "error",
			Error:       err.Error(),
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

	return ExternalServiceDetails{
		Description: description,
		Version:     version,
		Status:      data.Status,
		Error:       data.Output,
	}
}
