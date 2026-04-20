-- migrations/000008_create_provinces_table.up.sql

CREATE TABLE IF NOT EXISTS provinces (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    code VARCHAR(50) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL UNIQUE,
    is_active BOOLEAN DEFAULT TRUE,
    created_at DATETIME(3) NULL,
    updated_at DATETIME(3) NULL
);
