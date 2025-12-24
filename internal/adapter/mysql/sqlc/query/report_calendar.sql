-- ============================================
-- REPORT_CALENDAR queries
-- ============================================

-- name: GetDays :many
SELECT
    rc.id,
    rc.day,
    rc.month,
    rc.year,
    COALESCE(rc.description, '') as description,  -- Возвращаем пустую строку вместо NULL
    rc.is_paid_vacation,
    rc.type_id,
    rt.name as type_name,
    rt.system_name as type_system_name
FROM report_calendar rc
INNER JOIN report_type rt ON rc.type_id = rt.id
WHERE rc.month = ? AND rc.year = ?
ORDER BY rc.day ASC;

-- name: GetDaysByType :many
SELECT
    rc.id,
    rc.day,
    rc.month,
    rc.year,
    COALESCE(rc.description, '') as description,  -- Возвращаем пустую строку вместо NULL
    rc.is_paid_vacation,
    rc.type_id,
    rt.name as type_name,
    rt.system_name as type_system_name
FROM report_calendar rc
INNER JOIN report_type rt ON rc.type_id = rt.id
WHERE rc.month = ? AND rc.year = ? AND rt.system_name = ?
ORDER BY rc.day ASC;

-- name: GetAllDays :many
SELECT
    rc.id,
    rc.day,
    rc.month,
    rc.year,
    COALESCE(rc.description, '') as description,  -- Возвращаем пустую строку вместо NULL
    rc.is_paid_vacation,
    rc.type_id,
    rt.name as type_name,
    rt.system_name as type_system_name
FROM report_calendar rc
INNER JOIN report_type rt ON rc.type_id = rt.id
WHERE rc.year = ?
ORDER BY rc.day ASC;

-- name: GetAllDaysByType :many
SELECT
    rc.id,
    rc.day,
    rc.month,
    rc.year,
    COALESCE(rc.description, '') as description,  -- Возвращаем пустую строку вместо NULL
    rc.is_paid_vacation,
    rc.type_id,
    rt.name as type_name,
    rt.system_name as type_system_name
FROM report_calendar rc
INNER JOIN report_type rt ON rc.type_id = rt.id
WHERE rc.year = ? AND rt.system_name = ?
ORDER BY rc.day ASC;

-- name: GetDay :one
SELECT
    rc.id,
    rc.day,
    rc.month,
    rc.year,
    rc.description,
    rc.is_paid_vacation,
    rc.type_id,
    rt.name as type_name,
    rt.system_name as type_system_name
FROM report_calendar rc
INNER JOIN report_type rt ON rc.type_id = rt.id
WHERE rc.day = ? AND rc.month = ? AND rc.year = ?;

-- name: CreateDay :exec
INSERT INTO report_calendar (id, day, month, year, description, is_paid_vacation, type_id)
VALUES (?, ?, ?, ?, ?, ?, ?);

-- name: UpdateDay :exec
UPDATE report_calendar
SET description = ?, type_id = ?
WHERE id = ?;

-- name: DeleteCalendarDay :exec
DELETE FROM report_calendar
WHERE id = ?;

-- name: CheckDayExists :one
SELECT COUNT(*) as exists_count
FROM report_calendar
WHERE day = ? AND month = ? AND year = ?;
