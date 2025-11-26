DROP TRIGGER IF EXISTS update_permissions_updated_at ON permissions;
DROP INDEX IF EXISTS idx_permissions_resource_action;
DROP INDEX IF EXISTS idx_permissions_permission_name;
DROP TABLE IF EXISTS permissions;