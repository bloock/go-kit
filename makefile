mocks:
	mockgen -source=query/query.go -destination mocks/query/mock_query.go
	mockgen -source=command/command.go -destination mocks/command/mock_command.go
	mockgen -source=request/http_request.go -destination mocks/request/mock_http_request.go
	mockgen -source=cache/cache.go -destination mocks/cache/mock_cache.go
	mockgen -source=event/event.go -destination mocks/event/mock_event.go
	mockgen -source=publisher/amqp_publisher.go -destination mocks/publisher/mock_amqp_publisher.go
