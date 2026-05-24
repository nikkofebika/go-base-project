CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE media (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    uuid UUID NOT NULL UNIQUE DEFAULT gen_random_uuid(),
    media_id BIGINT NOT NULL,
    media_type VARCHAR(255) NOT NULL,
    collection_name VARCHAR(255) NOT NULL,
    filename VARCHAR(255) NOT NULL,
    disk VARCHAR(255) NOT NULL,
    path VARCHAR(255) NOT NULL,
    size BIGINT NOT NULL,
    mime_type VARCHAR(255) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'pending',

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_media_media_id ON media(media_id);

CREATE INDEX idx_media_media_type ON media(media_type);