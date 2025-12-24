-- ============================================
-- REPORT_USER queries
-- ============================================

-- name: GetUserMonthReport :many
SELECT
    ru.id,
    ru.user_id,
    ru.day,
    ru.month,
    ru.year,
    ru.hours,
    ru.type_id,
    rt.name as type_name,
    rt.system_name as type_system_name
FROM report_user ru
INNER JOIN report_type rt ON ru.type_id = rt.id
WHERE ru.user_id = ? AND ru.month = ? AND ru.year = ?
ORDER BY ru.day ASC;

-- name: GetUserDayReport :one
SELECT
    ru.id,
    ru.user_id,
    ru.day,
    ru.month,
    ru.year,
    ru.hours,
    ru.type_id,
    rt.name as type_name,
    rt.system_name as type_system_name
FROM report_user ru
INNER JOIN report_type rt ON ru.type_id = rt.id
WHERE ru.id = ?;

-- name: GetMonthTotalHours :one
SELECT CAST(COALESCE(SUM(hours), 0.0) AS FLOAT) AS total_hours
FROM report_user
WHERE user_id = ? AND month = ? AND year = ?;

-- name: CountDaysByType :one
SELECT COUNT(DISTINCT day) as days_count
FROM report_user
WHERE user_id = ? AND month = ? AND year = ? AND type_id = ?;

-- name: CountDaysWork :one
SELECT COUNT(DISTINCT day) as days_count
FROM report_user ru
INNER JOIN report_type rt ON ru.type_id = rt.id
WHERE ru.user_id = ? AND ru.month = ? AND ru.year = ? AND (rt.system_name = 'work' OR rt.system_name = 'weekend');

-- name: CreateUserReport :exec
INSERT INTO report_user (id, user_id, day, month, year, hours, type_id)
VALUES (?, ?, ?, ?, ?, ?, ?);

-- name: UpdateUserReport :exec
UPDATE report_user
SET hours = ?, type_id = ?
WHERE id = ?;

-- name: CheckUserReportExists :one
SELECT COUNT(*) as exists_count
FROM report_user
WHERE user_id = ? AND day = ? AND month = ? AND year = ?;

-- name: DeleteUserReport :exec
DELETE FROM report_user
WHERE user_id = ? AND day = ? AND month = ? AND year = ?;
