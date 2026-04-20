-- migrations/000006_add_timestamps_to_role_permissions.down.sql

ALTER TABLE role_group_permissions
DROP COLUMN created_at,
DROP COLUMN updated_at;
