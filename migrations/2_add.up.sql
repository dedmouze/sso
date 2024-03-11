INSERT INTO admins (id, email, level)
VALUES (1, 'admin', 3)
ON CONFLICT DO NOTHING;

INSERT INTO apps (id, name, apiKey)
VALUES (1, 'adminApp', '12345')
ON CONFLICT DO NOTHING;