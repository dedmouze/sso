INSERT INTO admins (id, email, level)
VALUES (1, 'admin', 3)
ON CONFLICT DO NOTHING;

INSERT INTO apps (id, name, secret) VALUES 
(1, 'url-shortener', 'test-secret-1'),
(2, 'permission', 'test-secret-2'),
(3, 'userInfo', 'test-secret-3')
ON CONFLICT DO NOTHING;

