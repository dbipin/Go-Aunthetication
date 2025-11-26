-- Insert default roles
INSERT INTO roles (role_name, description) VALUES
    ('admin', 'Full system access'),
    ('user', 'Standard user access'),
    ('moderator', 'Content moderation access')
ON CONFLICT (role_name) DO NOTHING;

-- Insert default permissions
INSERT INTO permissions (permission_name, resource, action, description) VALUES
    ('users.create', 'users', 'create', 'Create users'),
    ('users.read', 'users', 'read', 'View users'),
    ('users.update', 'users', 'update', 'Update users'),
    ('users.delete', 'users', 'delete', 'Delete users'),
    ('users.list', 'users', 'list', 'List users')
ON CONFLICT (resource, action) DO NOTHING;

-- Assign permissions to admin role
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.role_name = 'admin'
ON CONFLICT DO NOTHING;

-- Assign basic permissions to user role
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r, permissions p
WHERE r.role_name = 'user' 
AND p.permission_name IN ('users.read', 'users.update')
ON CONFLICT DO NOTHING;