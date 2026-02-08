CREATE TABLE IF NOT EXISTS key_value_cache (
    key TEXT PRIMARY KEY,
    value TEXT NOT NULL,
    expires_at DATETIME NOT NULL
);