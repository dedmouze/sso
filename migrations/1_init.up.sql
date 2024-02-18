CREATE TABLE IF NOT EXISTS users
(
    id          INTEGER PRIMARY KEY,
    email       TEXT NOT NULL UNIQUE,
    pass_hash   BLOB NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_email ON users(email);

CREATE TABLE IF NOT EXISTS apps
(
    id      INTEGER PRIMARY KEY,
    name    TEXT NOT NULL UNIQUE,
    secret  TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS admins
(
    id INTEGER PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    level INTEGER NOT NULL CHECK(level IN(1, 2, 3))
);
CREATE INDEX IF NOT EXISTS idx_email ON admins(email);