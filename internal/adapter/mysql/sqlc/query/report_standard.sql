-- ============================================
-- REPORT_STANDART queries
-- ============================================

-- name: GetStandard :one
SELECT id, month, year, hours, gender_id
FROM report_standard
WHERE month = ? AND year = ? AND gender_id = ?;

-- name: GetStandardByMonth :many
SELECT id, month, year, hours, gender_id
FROM report_standard
WHERE month = ? AND year = ?;

-- name: GetStandardByYear :many
SELECT id, month, year, hours, gender_id
FROM report_standard
WHERE year = ?
ORDER BY month ASC;

-- name: CreateStandard :exec
INSERT INTO report_standard (id, month, year, hours, gender_id)
VALUES (?, ?, ?, ?, ?);

-- name: UpdateStandard :exec
UPDATE report_standard
SET hours = ?
WHERE id = ?;

-- name: DeleteStandard :exec
DELETE FROM report_standard
WHERE id = ?;

-- name: CheckStandard :one
SELECT COUNT(*) as exists_count
FROM report_standard
WHERE month = ? AND year = ? AND gender_id = ?;