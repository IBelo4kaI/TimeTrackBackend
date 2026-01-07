package standard

import (
	repo "TimeTrack/internal/adapter/mysql/sqlc"
	"context"
	"database/sql"
)

type Service interface {
	ListForSetting(ctx context.Context, year int32) (*[]repo.ReportStandard, error)
	Create(ctx context.Context, prm repo.CreateStandardParams) (*repo.ReportStandard, error)
	Update(ctx context.Context, prm repo.UpdateStandardParams) error
}

type service struct {
	repo repo.Querier
	db   *sql.DB
}

func NewService(repo repo.Querier, db *sql.DB) Service {
	return &service{repo: repo, db: db}
}

func (s *service) ListForSetting(ctx context.Context, year int32) (*[]repo.ReportStandard, error) {
	standards, err := s.repo.GetStandardByYear(ctx, year)
	if err != nil {
		return nil, err
	}

	return &standards, nil
}

func (s *service) Create(ctx context.Context, prm repo.CreateStandardParams) (*repo.ReportStandard, error) {
	err := s.repo.CreateStandard(ctx, prm)
	if err != nil {
		return nil, err
	}

	standard, err := s.repo.GetStandard(ctx, repo.GetStandardParams{Month: prm.Month, Year: prm.Year, GenderID: prm.GenderID})
	if err != nil {
		return nil, err
	}

	return &standard, nil
}

func (s *service) Update(ctx context.Context, prm repo.UpdateStandardParams) error {
	if err := s.repo.UpdateStandard(ctx, prm); err != nil {
		return err
	}

	return nil
}
