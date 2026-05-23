CREATE TABLE media (
    id INT PRIMARY KEY AUTO_INCREMENT,
    uuid VARCHAR(255) NOT NULL,
    media_id INT UNSIGNED NOT NULL,
    media_type VARCHAR(255) NOT NULL,
    collection_name VARCHAR(255) NOT NULL,
    filename VARCHAR(255) NOT NULL,
    disk VARCHAR(255) NOT NULL,
    path VARCHAR(255) NOT NULL,
    size INT UNSIGNED NOT NULL,
    mime_type VARCHAR(255) NOT NULL,
    status VARCHAR(255) NOT NULL DEFAULT 'pending', -- will be enum (pending, success, failed)

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    INDEX(media_id),
    INDEX(media_type)
) ENGINE=InnoDB;
