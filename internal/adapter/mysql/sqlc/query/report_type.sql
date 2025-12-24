-- ============================================
-- REPORT_TYPE queries
-- ============================================

-- name: GetReportTypeById :one
SELECT id, name, system_name
FROM report_type
WHERE id = ?;

-- name: GetReportTypeBySystemName :one
SELECT id, name, system_name
FROM report_type
WHERE system_name = ?;

-- name: GetAllReportTypes :many
SELECT id, name, system_name
FROM report_type
ORDER BY name ASC;

-- name: CreateReportType :exec
INSERT INTO report_type (id, name, system_name)
VALUES (?, ?, ?);

-- name: UpdateReportType :exec
UPDATE report_type
SET name = ?, system_name = ?
WHERE id = ?;

-- name: DeleteReportType :exec
DELETE FROM report_type
WHERE id = ?;