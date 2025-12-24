-- ============================================
-- REPORT_SETTING queries
-- ============================================

-- name: GetVacationDuration :one
SELECT vacation_duration
FROM report_setting
WHERE id = 1

