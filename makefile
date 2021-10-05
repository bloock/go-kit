mocks:
	mockgen -package=mocks -source=query/query.go -destination mocks/mock_query.go
	mockgen -package=mocks -source=command/command.go -destination mocks/mock_command.go
	mockgen -package=mocks -source=request/http_request.go -destination mocks/mock_http_request.go
	mockgen -package=mocks -source=cache/cache.go -destination mocks/mock_cache.go
	mockgen -package=mocks -source=event/event.go -destination mocks/mock_event.go
