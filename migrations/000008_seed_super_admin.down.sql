-- Remove super admin
DELETE FROM user_roles 
WHERE user_id = (SELECT id FROM users WHERE email = 'admin@gmail.com');

DELETE FROM users WHERE email = 'admin@gmail.com';