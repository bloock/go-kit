CREATE TABLE IF NOT EXISTS cache_usage
(
    _key     VARCHAR(255) PRIMARY KEY,
    value   BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
)