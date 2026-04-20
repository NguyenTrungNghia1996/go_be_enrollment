-- migrations/000004_update_role_group_schema.down.sql

ALTER TABLE role_groups
DROP COLUMN code,
DROP COLUMN is_active;
