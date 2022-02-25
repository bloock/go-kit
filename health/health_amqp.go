package health

import (
	"github.com/bloock/go-kit/client"
)

type HealthAMQP struct {
	client      *client.AMQPClient
	description string
	version     string
}

func NewHealtAMQP(client *client.AMQPClient, description string) HealthAMQP {
	return HealthAMQP{
		client:      client,
		description: description,
		version:     DepVersion("streadway/amqp"),
	}
}

func (h HealthAMQP) HealthCheck() ExternalServiceDetails {
	s := "pass"
	var e string

	err := h.client.Ping()
	if err != nil {
		s = "error"
		e = err.Error()
	}

	return ExternalServiceDetails{
		Description: h.description,
		Version:     h.version,
		Status:      s,
		Error:       e,
	}
}
