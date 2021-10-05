mocks:
	mockgen -source=query/query.go -destination mocks/mock_query.go
	mockgen -source=command/command.go -destination mocks/mock_command.go
	mockgen -source=request/http_request.go -destination mocks/mock_http_request.go
	mockgen -source=cache/cache.go -destination mocks/mock_cache.go
	mockgen -source=event/event.go -destination mocks/mock_event.go
