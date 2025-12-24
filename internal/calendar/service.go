package calendar

import (
	repo "TimeTrack/internal/adapter/mysql/sqlc"
	"context"
	"database/sql"
)

type Service interface {
	List(ctx context.Context, prm repo.GetDaysParams) ([]repo.GetDaysRow, error)
}

type service struct {
	repo repo.Querier
	db   *sql.DB
}

func NewService(repo repo.Querier, db *sql.DB) Service {
	return &service{repo: repo, db: db}
}

func (s *service) List(ctx context.Context, prm repo.GetDaysParams) ([]repo.GetDaysRow, error) {
	calendar, err := s.repo.GetDays(ctx, prm)
	if err != nil {
		return nil, err
	}

	return calendar, nil
}
