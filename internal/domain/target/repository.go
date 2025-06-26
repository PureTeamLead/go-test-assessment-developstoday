package target

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db      *pgxpool.Pool
	builder sq.StatementBuilderType
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{db: pool}
}

func (r *Repository) AddTarget() {}

func (r *Repository) UpdateTargetStatus() {}

func (r *Repository) UpdateTargetNotes() {}

func (r *Repository) AddMissionID() {}

func (r *Repository) GetTargetsByMissionID() {}

func (r *Repository) GetTargetByID() {}
