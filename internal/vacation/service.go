package vacation

import (
	repo "TimeTrack/internal/adapter/mysql/sqlc"
	"context"
	"database/sql"
	"fmt"
	"time"
)

type Service interface {
	List(ctx context.Context, userID string, year int32) (*[]vacationRow, error)
	ListAll(ctx context.Context, year int32) (*[]vacationRow, error)
	Stats(ctx context.Context, userID string, year int32) (*vacationStats, error)
	Create(ctx context.Context, prm repo.CreateVacationParams) (*repo.GetVacationByIdRow, error)
	ChangeStatus(ctx context.Context, prm repo.UpdateVacationStatusParams) error
	Years(ctx context.Context, userID string) (*[]int32, error)
	Delete(ctx context.Context, id string) error
}

type service struct {
	repo repo.Querier
	db   *sql.DB
}

func NewService(repo repo.Querier, db *sql.DB) Service {
	return &service{repo: repo, db: db}
}

type vacationStats struct {
	Approved      int32 `json:"approved"`
	Consideration int32 `json:"consideration"`
	Free          int32 `json:"free"`
	All           int32 `json:"all"`
}

func (s *service) Stats(ctx context.Context, userID string, year int32) (*vacationStats, error) {
	all, err := s.repo.GetSettingVacationDuration(ctx)
	if err != nil {
		return nil, err
	}

	vacations, err := s.List(ctx, userID, year)
	if err != nil {
		return nil, err
	}

	var approved int32
	var consideration int32
	free := all

	for _, vacation := range *vacations {
		days := int32(vacation.CountDay)

		switch vacation.Status {
		case repo.ReportVacationStatusApproved:
			approved += days
			free -= days

		case repo.ReportVacationStatusConsideration:
			consideration += days
			free -= days
		}
	}

	if free < 0 {
		free = 0
	}

	return &vacationStats{
		Approved:      approved,
		Consideration: consideration,
		Free:          free,
		All:           all,
	}, nil
}

type vacationRow struct {
	ID          string                             `json:"id"`
	UserID      string                             `json:"userId"`
	StartDate   time.Time                          `json:"startDate"`
	EndDate     time.Time                          `json:"endDate"`
	Year        int32                              `json:"year"`
	Description string                             `json:"description"`
	Status      repo.ReportVacationStatus          `json:"status"`
	CountDay    int16                              `json:"countDay"`
	Holidays    []repo.GetCalendarDaysAllByTypeRow `json:"holidays"`
	CreateAt    time.Time                          `json:"createAt"`
}

func (s *service) List(ctx context.Context, userID string, year int32) (*[]vacationRow, error) {
	vacations, err := s.repo.GetVacationsByYear(ctx, repo.GetVacationsByYearParams{UserID: userID, Year: year})
	if err != nil {
		return nil, err
	}

	holidays, err := s.repo.GetCalendarDaysAllByType(ctx, repo.GetCalendarDaysAllByTypeParams{Year: year, SystemName: "holiday"})
	if err != nil {
		return nil, err
	}

	// Создаём map для быстрого поиска: "MM-DD" -> holiday
	holidayMap := make(map[string]repo.GetCalendarDaysAllByTypeRow)
	for _, h := range holidays {
		// Ключ формата "01-15" (месяц-день)
		key := fmt.Sprintf("%02d-%02d", h.Month, h.Day)
		holidayMap[key] = h
	}

	vacationRows := make([]vacationRow, len(vacations))

	for i, v := range vacations {
		vacationHolidays := findHolidaysInRange(holidayMap, v.StartDate, v.EndDate)

		countDay := countVacationDays(
			holidayMap,
			v.StartDate,
			v.EndDate,
		)

		vacationRows[i] = vacationRow{
			ID:          v.ID,
			UserID:      v.UserID,
			StartDate:   v.StartDate,
			EndDate:     v.EndDate,
			Year:        v.Year,
			Description: v.Description,
			Status:      v.Status,
			CountDay:    countDay,
			Holidays:    vacationHolidays,
			CreateAt:    v.CreateAt,
		}
	}

	return &vacationRows, nil
}

func (s *service) ListAll(ctx context.Context, year int32) (*[]vacationRow, error) {
	vacations, err := s.repo.GetAdminVacationsByYear(ctx, year)
	if err != nil {
		return nil, err
	}

	holidays, err := s.repo.GetCalendarDaysAllByType(ctx, repo.GetCalendarDaysAllByTypeParams{Year: year, SystemName: "holiday"})
	if err != nil {
		return nil, err
	}

	// Создаём map для быстрого поиска: "MM-DD" -> holiday
	holidayMap := make(map[string]repo.GetCalendarDaysAllByTypeRow)
	for _, h := range holidays {
		// Ключ формата "01-15" (месяц-день)
		key := fmt.Sprintf("%02d-%02d", h.Month, h.Day)
		holidayMap[key] = h
	}

	vacationRows := make([]vacationRow, len(vacations))

	for i, v := range vacations {
		vacationHolidays := findHolidaysInRange(holidayMap, v.StartDate, v.EndDate)

		countDay := countVacationDays(
			holidayMap,
			v.StartDate,
			v.EndDate,
		)

		vacationRows[i] = vacationRow{
			ID:          v.ID,
			UserID:      v.UserID,
			StartDate:   v.StartDate,
			EndDate:     v.EndDate,
			Year:        v.Year,
			Description: v.Description,
			Status:      v.Status,
			CountDay:    countDay,
			Holidays:    vacationHolidays,
			CreateAt:    v.CreateAt,
		}
	}

	return &vacationRows, nil
}

func findHolidaysInRange(holidayMap map[string]repo.GetCalendarDaysAllByTypeRow, startDate, endDate time.Time) []repo.GetCalendarDaysAllByTypeRow {
	var result []repo.GetCalendarDaysAllByTypeRow

	// Итерируемся по каждому дню в диапазоне отпуска
	for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {
		key := fmt.Sprintf("%02d-%02d", d.Month(), d.Day())
		if holiday, exists := holidayMap[key]; exists {
			result = append(result, holiday)
		}
	}

	return result
}

func countVacationDays(
	holidayMap map[string]repo.GetCalendarDaysAllByTypeRow,
	startDate, endDate time.Time,
) int16 {
	var count int16

	for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {
		key := fmt.Sprintf("%02d-%02d", d.Month(), d.Day())

		holiday, exists := holidayMap[key]

		// Если это праздник и он НЕ входит в оплачиваемый отпуск — пропускаем
		if exists && !holiday.IsPaidVacation {
			continue
		}

		count++
	}

	return count
}

func (s *service) Create(ctx context.Context, prm repo.CreateVacationParams) (*repo.GetVacationByIdRow, error) {
	err := s.repo.CreateVacation(ctx, prm)
	if err != nil {
		return nil, err
	}

	vacation, err := s.repo.GetVacationById(ctx, prm.ID)

	if err != nil {
		return nil, err
	}

	return &vacation, nil
}

func (s *service) ChangeStatus(ctx context.Context, prm repo.UpdateVacationStatusParams) error {
	err := s.repo.UpdateVacationStatus(ctx, prm)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) Years(ctx context.Context, userID string) (*[]int32, error) {
	years, err := s.repo.GetYearsVacation(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &years, nil
}

func (s *service) Delete(ctx context.Context, id string) error {
	if err := s.repo.DeleteVacation(ctx, id); err != nil {
		return err
	}
	return nil
}
