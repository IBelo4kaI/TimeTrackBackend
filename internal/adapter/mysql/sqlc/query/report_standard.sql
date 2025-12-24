-- ============================================
-- REPORT_STANDART queries
-- ============================================

-- name: GetStandardHours :one
SELECT id, month, year, hours, gender_id
FROM report_standard
WHERE month = ? AND year = ? AND gender_id = ?;

-- name: GetStandardHoursByMonth :many
SELECT id, month, year, hours, gender_id
FROM report_standard
WHERE month = ? AND year = ?;

-- name: CreateStandardHours :exec
INSERT INTO report_standard (id, month, year, hours, gender_id)
VALUES (?, ?, ?, ?, ?);

-- name: UpdateStandardHours :exec
UPDATE report_standard
SET hours = ?
WHERE id = ?;

-- name: DeleteStandardHours :exec
DELETE FROM report_standard
WHERE id = ?;

-- name: CheckStandardExists :one
SELECT COUNT(*) as exists_count
FROM report_standard
WHERE month = ? AND year = ? AND gender_id = ?;