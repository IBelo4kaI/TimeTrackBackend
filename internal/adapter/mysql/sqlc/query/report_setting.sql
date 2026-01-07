-- ============================================
-- REPORT_SETTING queries
-- ============================================

-- name: GetSettingVacationDuration :one
SELECT vacation_duration
FROM report_setting
WHERE id = 1

