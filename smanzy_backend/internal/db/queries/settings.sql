-- name: GetSetting :one
SELECT value FROM settings WHERE key = $1;

-- name: UpsertSetting :one
INSERT INTO settings (key, value, updated_at)
VALUES ($1, $2, (EXTRACT(EPOCH FROM NOW()) * 1000)::BIGINT)
ON CONFLICT (key) DO UPDATE SET
    value = EXCLUDED.value,
    updated_at = EXCLUDED.updated_at
RETURNING *;

-- name: ListSettings :many
SELECT key, value FROM settings;
