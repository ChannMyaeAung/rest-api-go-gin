CREATE TABLE IF NOT EXISTS users_tmp (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    email TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL,
    password TEXT NOT NULL
);

INSERT INTO users_tmp (id, email, name, password)
SELECT id, email, name, password FROM users;

DROP TABLE users;

ALTER TABLE users_tmp RENAME TO users;
