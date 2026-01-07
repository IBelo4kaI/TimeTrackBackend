package types

import (
	repo "TimeTrack/internal/adapter/mysql/sqlc"
	"context"
	"database/sql"
)

type Service interface {
	List(ctx context.Context) (*[]repo.ReportType, error)
}

type service struct {
	repo repo.Querier
	db   *sql.DB
}

func NewService(repo repo.Querier, db *sql.DB) Service {
	return &service{repo: repo, db: db}
}

func (s *service) List(ctx context.Context) (*[]repo.ReportType, error) {
	types, err := s.repo.GetTypeAll(ctx)
	if err != nil {
		return nil, err
	}

	return &types, nil
}
