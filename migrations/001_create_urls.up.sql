-- Схема применяется приложением на старте (internal/repository/url/schema.go).
-- Файл держим синхронным с ней, чтобы схему можно было поднять вручную.
CREATE TABLE IF NOT EXISTS urls (
    id           BIGSERIAL PRIMARY KEY,
    original_url TEXT NOT NULL UNIQUE,
    uuid         VARCHAR(255) NOT NULL UNIQUE,
    user_id      TEXT NOT NULL DEFAULT '',
    is_deleted   BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE UNIQUE INDEX IF NOT EXISTS urls_original_url_uidx ON urls (original_url);

CREATE INDEX IF NOT EXISTS urls_user_id_idx ON urls (user_id) WHERE is_deleted = FALSE;
