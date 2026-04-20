-- migrations/000007_update_menu_schema.down.sql

ALTER TABLE menus
CHANGE COLUMN title name VARCHAR(100) NOT NULL,
CHANGE COLUMN url path VARCHAR(255) NOT NULL,
DROP COLUMN is_active,
DROP COLUMN created_at,
DROP COLUMN updated_at;
