CREATE TABLE IF NOT EXISTS cache_usage
(
    _key     VARCHAR(255) PRIMARY KEY,
    value   INT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
) ENGINE = InnoDB
    CHARSET = utf8mb4
    COLLATE = utf8mb4_general_ci;