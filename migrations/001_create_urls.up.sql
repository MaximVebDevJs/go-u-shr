CREATE TABLE IF NOT EXISTS urls (
    id           BIGSERIAL PRIMARY KEY,
    original_url TEXT NOT NULL,
    uuid         VARCHAR(255) NOT NULL UNIQUE
);
