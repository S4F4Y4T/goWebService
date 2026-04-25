CREATE TABLE IF NOT EXISTS credentials (
    id            SERIAL PRIMARY KEY,
    user_id       INTEGER UNIQUE NOT NULL,
    email         VARCHAR(255) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    created_at    TIMESTAMPTZ DEFAULT NOW()
);
