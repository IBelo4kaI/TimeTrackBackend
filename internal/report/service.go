package report

import (
	repo "TimeTrack/internal/adapter/mysql/sqlc"
	"context"
	"database/sql"
	"fmt"
)

type Service interface {
	List(ctx context.Context, prm repo.GetUserMonthReportParams) (*[]repo.GetUserMonthReportRow, error)
	Create(ctx context.Context, prm CreateReportParams) (*ReportResponse, error)
	Update(ctx context.Context, prm UpdateReportParams) (*ReportResponse, error)
	Delete(ctx context.Context, prm repo.DeleteUserReportParams) error
}

type service struct {
	repo repo.Querier
	db   *sql.DB
}

func NewService(repo repo.Querier, db *sql.DB) Service {
	return &service{repo: repo, db: db}
}

type ReportResponse struct {
	ID             string  `json:"id"`
	UserID         string  `json:"userId"`
	Day            int32   `json:"day"`
	Month          int32   `json:"month"`
	Year           int32   `json:"year"`
	Hours          float64 `json:"hours"`
	TypeID         string  `json:"typeId"`
	TypeName       string  `json:"typeName"`
	TypeSystemName string  `json:"typeSystemName"`
}

type CreateReportParams struct {
	ID     string  `json:"id"`
	UserID string  `json:"userId"`
	Day    int32   `json:"day"`
	Month  int32   `json:"month"`
	Year   int32   `json:"year"`
	Hours  float64 `json:"hours"`
	Type   string  `json:"typeSystemName"`
}

type UpdateReportParams struct {
	ID    string  `json:"id"`
	Hours float64 `json:"hours"`
	Type  string  `json:"typeSystemName"`
}

func (s *service) List(ctx context.Context, prm repo.GetUserMonthReportParams) (*[]repo.GetUserMonthReportRow, error) {
	reports, err := s.repo.GetUserMonthReport(ctx, prm)
	if err != nil {
		return nil, fmt.Errorf("get user month report: %w", err)
	}

	return &reports, nil
}

func (s *service) Create(ctx context.Context, prm CreateReportParams) (*ReportResponse, error) {
	reportType, err := s.repo.GetReportTypeBySystemName(ctx, prm.Type)
	if err != nil {
		return nil, fmt.Errorf("get report type: %w", err)
	}

	if err := s.repo.CreateUserReport(ctx, repo.CreateUserReportParams{
		ID:     prm.ID,
		UserID: prm.UserID,
		Day:    prm.Day,
		Month:  prm.Month,
		Year:   prm.Year,
		Hours:  prm.Hours,
		TypeID: reportType.ID,
	}); err != nil {
		return nil, fmt.Errorf("create user report: %w", err)
	}

	return s.buildReportResponse(ctx, prm.ID)
}

func (s *service) Update(ctx context.Context, prm UpdateReportParams) (*ReportResponse, error) {
	reportType, err := s.repo.GetReportTypeBySystemName(ctx, prm.Type)
	if err != nil {
		return nil, fmt.Errorf("get report type: %w", err)
	}

	if err := s.repo.UpdateUserReport(ctx, repo.UpdateUserReportParams{
		ID:     prm.ID,
		Hours:  prm.Hours,
		TypeID: reportType.ID,
	}); err != nil {
		return nil, fmt.Errorf("update user report: %w", err)
	}

	return s.buildReportResponse(ctx, prm.ID)
}

func (s *service) Delete(ctx context.Context, prm repo.DeleteUserReportParams) error {
	if err := s.repo.DeleteUserReport(ctx, prm); err != nil {
		return fmt.Errorf("delete user report: %w", err)
	}
	return nil
}

// monthStats содержит агрегированную статистику за месяц
type monthStats struct {
	totalHours  float64
	workDays    int64
	medicalDays int64
}

// getMonthStats получает всю статистику за месяц одним вызовом
func (s *service) GetMonthStats(ctx context.Context, userID string, month, year int32) (*monthStats, error) {
	// Получение общих часов
	totalHours, err := s.repo.GetMonthTotalHours(ctx, repo.GetMonthTotalHoursParams{
		UserID: userID,
		Month:  month,
		Year:   year,
	})
	if err != nil {
		return nil, fmt.Errorf("get total hours: %w", err)
	}

	// Подсчет рабочих дней
	workDays, err := s.repo.CountDaysWork(ctx, repo.CountDaysWorkParams{
		UserID: userID,
		Month:  month,
		Year:   year,
	})
	if err != nil {
		return nil, fmt.Errorf("count work days: %w", err)
	}

	// Получение типа больничных
	medicalType, err := s.repo.GetReportTypeBySystemName(ctx, "medical")
	if err != nil {
		return nil, fmt.Errorf("get medical type: %w", err)
	}

	// Подсчет дней больничных
	medicalDays, err := s.repo.CountDaysByType(ctx, repo.CountDaysByTypeParams{
		UserID: userID,
		Month:  month,
		Year:   year,
		TypeID: medicalType.ID,
	})
	if err != nil {
		return nil, fmt.Errorf("count medical days: %w", err)
	}

	return &monthStats{
		totalHours:  totalHours,
		workDays:    workDays,
		medicalDays: medicalDays,
	}, nil
}

// buildReportResponse создает ответ с отчетом и статистикой
func (s *service) buildReportResponse(ctx context.Context, reportID string) (*ReportResponse, error) {
	report, err := s.repo.GetUserDayReport(ctx, reportID)
	if err != nil {
		return nil, fmt.Errorf("get user day report: %w", err)
	}

	return &ReportResponse{
		ID:             report.ID,
		UserID:         report.UserID,
		Day:            report.Day,
		Month:          report.Month,
		Year:           report.Year,
		Hours:          report.Hours,
		TypeID:         report.TypeID,
		TypeName:       report.TypeName,
		TypeSystemName: report.TypeSystemName,
	}, nil
}
