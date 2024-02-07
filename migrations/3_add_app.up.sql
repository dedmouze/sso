INSERT INTO apps (id, name, secret)
VALUES (1, "test", "t-secret")
ON CONFLICT DO NOTHING;