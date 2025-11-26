-- Create super admin user (password: admin123)
INSERT INTO users (email, password, name)
VALUES (
    'admin@gmail.com',
    '$2a$10$XATULCpd1nG8pgE2k38jWeP3pQ/fHPokEwDZEM.HTg2Jf22U.sM.W',  -- PASTE HASH HERE
    'Super Admin'
) ON CONFLICT (email) DO NOTHING;

-- Assign admin role to super admin
INSERT INTO user_roles (user_id, role_id)
SELECT u.id, r.id
FROM users u, roles r
WHERE u.email = 'admin@gmail.com'
AND r.role_name = 'admin'
ON CONFLICT DO NOTHING;