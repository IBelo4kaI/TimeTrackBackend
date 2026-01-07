-- ============================================
-- REPORT_TYPE queries
-- ============================================

-- name: GetTypeById :one
SELECT id, name, system_name
FROM report_type
WHERE id = ?;

-- name: GetTypeBySystemName :one
SELECT id, name, system_name
FROM report_type
WHERE system_name = ?;

-- name: GetTypeAll :many
SELECT id, name, system_name
FROM report_type
ORDER BY name ASC;

-- name: CreateType :exec
INSERT INTO report_type (id, name, system_name)
VALUES (?, ?, ?);

-- name: UpdateType :exec
UPDATE report_type
SET name = ?, system_name = ?
WHERE id = ?;

-- name: DeleteType :exec
DELETE FROM report_type
WHERE id = ?;