-- migrations/000004_update_role_group_schema.up.sql

ALTER TABLE role_groups
ADD COLUMN code VARCHAR(50) UNIQUE AFTER id,
ADD COLUMN is_active BOOLEAN DEFAULT TRUE AFTER description;
