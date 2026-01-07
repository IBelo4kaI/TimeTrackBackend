package calendar

import (
	repo "TimeTrack/internal/adapter/mysql/sqlc"
	"context"
	"database/sql"
)

type Service interface {
	ListMonth(ctx context.Context, prm repo.GetCalendarDaysParams) (*[]repo.GetCalendarDaysRow, error)
	ListYear(ctx context.Context, year int32) (*[]repo.GetCalendarDaysAllRow, error)
	Create(ctx context.Context, prm repo.CreateCalendarDayParams) (*repo.GetCalendarDayRow, error)
	Delete(ctx context.Context, id string) error
}

type service struct {
	repo repo.Querier
	db   *sql.DB
}

func NewService(repo repo.Querier, db *sql.DB) Service {
	return &service{repo: repo, db: db}
}

func (s *service) ListMonth(ctx context.Context, prm repo.GetCalendarDaysParams) (*[]repo.GetCalendarDaysRow, error) {
	calendar, err := s.repo.GetCalendarDays(ctx, prm)
	if err != nil {
		return nil, err
	}

	return &calendar, nil
}

func (s *service) ListYear(ctx context.Context, year int32) (*[]repo.GetCalendarDaysAllRow, error) {
	calendar, err := s.repo.GetCalendarDaysAll(ctx, year)
	if err != nil {
		return nil, err
	}

	return &calendar, nil
}

func (s *service) Create(ctx context.Context, prm repo.CreateCalendarDayParams) (*repo.GetCalendarDayRow, error) {
	err := s.repo.CreateCalendarDay(ctx, prm)
	if err != nil {
		return nil, err
	}

	calendar, err := s.repo.GetCalendarDay(ctx, repo.GetCalendarDayParams{Day: prm.Day, Month: prm.Month, Year: prm.Year})
	if err != nil {
		return nil, err
	}

	return &calendar, nil
}

func (s *service) Delete(ctx context.Context, id string) error {
	if err := s.repo.DeleteCalendarDay(ctx, id); err != nil {
		return err
	}
	return nil
}
