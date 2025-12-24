-- ============================================
-- REPORT_VACATION queries
-- ============================================

-- name: GetUserVacations :many
SELECT id, user_id, start_date, end_date, year, COALESCE(description, '') as description, status, create_at
FROM report_vacation
WHERE user_id = ?
ORDER BY create_at DESC;

-- name: GetUserVacationsByYear :many
SELECT id, user_id, start_date, end_date, year, COALESCE(description, '') as description, status, create_at
FROM report_vacation
WHERE user_id = ? AND year = ?
ORDER BY create_at DESC;

-- name: GetVacationsByYear :many
SELECT id, user_id, start_date, end_date, year, COALESCE(description, '') as description, status, create_at
FROM report_vacation
WHERE year = ?
ORDER BY create_at DESC;

-- name: GetVacationById :one
SELECT id, user_id, start_date, end_date, year, COALESCE(description, '') as description, status, create_at
FROM report_vacation
WHERE id = ?;

-- name: GetVacationApproved :many
SELECT id, user_id, start_date, end_date, year, COALESCE(description, '') as description, status, create_at
FROM report_vacation
WHERE user_id = ? AND status = "approved";

-- name: GetYearsVacation :many
SELECT DISTINCT year
FROM report_vacation
WHERE user_id = ?
ORDER BY year DESC;

-- name: CreateVacation :exec
INSERT INTO report_vacation (id, user_id, start_date, end_date, year,  description, status)
VALUES (?, ?, ?, ?, ?, ?, ?);

-- name: UpdateVacationStatus :exec
UPDATE report_vacation
SET status = ?
WHERE id = ?;

-- name: UpdateVacation :exec
UPDATE report_vacation
SET start_date = ?, end_date = ?, year = ?, status = ?,  description = ?
WHERE id = ?;

-- name: DeleteVacation :exec
DELETE FROM report_vacation
WHERE id = ?;


