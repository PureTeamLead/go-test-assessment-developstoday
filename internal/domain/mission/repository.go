package mission

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db      *pgxpool.Pool
	builder sq.StatementBuilderType
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{
		db: pool,
	}
}

func (r *Repository) AddMission() {

}

func (r *Repository) DeleteMission() {

}

func (r *Repository) UpdateMission() {

}

func (r *Repository) AddCatID() {

}

func (r *Repository) GetMissions() {

}

func (r *Repository) GetMissionsByCatID() {

}

func (r *Repository) GetMissionByID() {}
