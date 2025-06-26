package target

import (
	"context"
	"errors"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/PureTeamLead/go-test-assessment-developstoday/internal/utils"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db      *pgxpool.Pool
	builder sq.StatementBuilderType
}

const (
	tableName       = "targets"
	idColumn        = "id"
	missionIDColumn = "mission_id"
	nameColumn      = "name"
	countryColumn   = "country"
	notesColumn     = "notes"
	stateColumn     = "state"
	createdAtColumn = "created_at"
	updatedAtColumn = "updated_at"
	completedState  = "completed"
)

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{db: pool}
}

func (r *Repository) AddTarget(ctx context.Context, missionID uuid.UUID, target *Target) (uuid.UUID, error) {
	const op = "target.Repository.AddTarget"
	var id uuid.UUID

	query, args, err := r.builder.
		Insert(tableName).
		Columns(missionIDColumn, nameColumn, countryColumn, notesColumn).
		Values(missionID, target.Name, target.Country, target.Notes).
		Suffix("RETURNING " + idColumn).
		ToSql()

	if err != nil {
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}

	err = r.db.QueryRow(ctx, query, args...).Scan(&id)
	if err != nil {
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (r *Repository) SetTargetCompleted(ctx context.Context, id uuid.UUID) error {
	const op = "target.Repository.SetTargetCompleted"

	query, args, err := r.builder.Update(tableName).
		Set(stateColumn, completedState).
		Where(sq.Eq{idColumn: id}).
		ToSql()

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	var res pgconn.CommandTag
	res, err = r.db.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if res.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", op, utils.ErrTargetNotFound)
	}

	return nil
}

func (r *Repository) UpdateTargetNotes(ctx context.Context, id uuid.UUID, notes string) error {
	const op = "target.Repository.UpdateTargetNotes"

	query, args, err := r.builder.Update(tableName).
		Set(notesColumn, notes).
		Where(sq.Eq{idColumn: id}).
		ToSql()

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	var res pgconn.CommandTag
	res, err = r.db.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if res.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", op, utils.ErrTargetNotFound)
	}

	return nil
}

func (r *Repository) GetTargetsByMissionID(ctx context.Context, missionID uuid.UUID) ([]*Target, error) {
	const op = "target.Repository.GetTargetsByMissionID"

	query, args, err := r.builder.
		Select(idColumn, nameColumn, countryColumn, notesColumn, stateColumn, createdAtColumn, updatedAtColumn).
		From(tableName).
		Where(sq.Eq{missionIDColumn: missionID}).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var rows pgx.Rows
	targets := make([]*Target, 0)

	rows, err = r.db.Query(ctx, query, args...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, utils.ErrMissionNotFound
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	for rows.Next() {
		var target Target

		err = rows.Scan(&target.ID, &target.Name, &target.Country, &target.Notes, &target.State, &target.CreatedAt, &target.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		targets = append(targets, &target)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return targets, nil
}

func (r *Repository) GetTargetByID(ctx context.Context, id uuid.UUID) (*Target, error) {
	const op = "target.Repository.GetTargetByID"
	var target Target

	query, args, err := r.builder.
		Select(idColumn, nameColumn, countryColumn, notesColumn, stateColumn, createdAtColumn, updatedAtColumn).
		From(tableName).
		Where(sq.Eq{idColumn: id}).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	err = r.db.QueryRow(ctx, query, args...).Scan(
		&target.ID,
		&target.Name,
		&target.Country,
		&target.Notes,
		&target.State,
		&target.CreatedAt,
		&target.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, utils.ErrTargetNotFound)
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &target, nil
}

func (r *Repository) DeleteTarget(ctx context.Context, id uuid.UUID) error {
	const op = "target.Repository.DeleteTarget"

	query, args, err := r.builder.Delete(tableName).Where(sq.Eq{idColumn: id}).ToSql()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if _, err = r.db.Exec(ctx, query, args...); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return utils.ErrTargetNotFound
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
