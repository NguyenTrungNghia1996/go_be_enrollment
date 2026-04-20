-- migrations/000006_add_timestamps_to_role_permissions.up.sql

ALTER TABLE role_group_permissions
ADD COLUMN created_at DATETIME(3) NULL,
ADD COLUMN updated_at DATETIME(3) NULL;
