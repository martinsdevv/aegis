CREATE TABLE IF NOT EXISTS api_keys (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    key TEXT UNIQUE NOT NULL,         -- aqui vai armazenar SHA256
    upstream_host TEXT,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    monthly_quota INTEGER NOT NULL DEFAULT 10000,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_api_keys_key ON api_keys(key);
