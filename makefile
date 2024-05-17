mocks:
	mockgen -source=cqrs/query/query.go -destination test_utils/mocks/query/mock_query.go
	mockgen -source=cqrs/command/command.go -destination test_utils/mocks/command/mock_command.go
	mockgen -source=http/http_request.go -destination test_utils/mocks/request/mock_http_request.go
	mockgen -source=cache/cache.go -destination test_utils/mocks/cache/mock_cache.go
	mockgen -source=domain/event.go -destination test_utils/mocks/event/mock_event.go
	mockgen -source=repository/publisher/publisher.go -destination test_utils/mocks/publisher/mock_publisher.go
	mockgen -source=cache/cache_usage_repository.go -destination test_utils/mocks/cache_usage/mock_cache_usage.go