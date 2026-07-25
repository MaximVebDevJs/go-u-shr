CREATE TABLE IF NOT EXISTS urls (
    id           BIGSERIAL PRIMARY KEY,
    original_url TEXT NOT NULL UNIQUE,
    uuid         VARCHAR(255) NOT NULL UNIQUE,
    user_id      TEXT NOT NULL,
    is_deleted   BOOLEAN NOT NULL DEFAULT FALSE
);
