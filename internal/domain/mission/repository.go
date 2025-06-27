package mission

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/PureTeamLead/go-test-assessment-developstoday/internal/domain/cat"
	"github.com/PureTeamLead/go-test-assessment-developstoday/internal/utils"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	tableName       = "missions"
	idColumn        = "id"
	catIDColumn     = "cat_id"
	stateColumn     = "state"
	createdAtColumn = "created_at"
	updatedAtColumn = "updated_at"
	completedState  = "completed"
	startedState    = "started"
)

type Repository struct {
	db      *pgxpool.Pool
	builder sq.StatementBuilderType
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	builder := sq.StatementBuilderType{}
	builder = builder.PlaceholderFormat(sq.Dollar)
	return &Repository{
		db:      pool,
		builder: builder,
	}
}

// TODO: check
func (r *Repository) AddMission(ctx context.Context) (uuid.UUID, error) {
	const op = "mission.Repository.AddMission"
	var id uuid.UUID

	query, args, err := r.builder.
		Insert(tableName).
		Columns(stateColumn).
		Values(startedState).
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

func (r *Repository) DeleteMission(ctx context.Context, id uuid.UUID) error {
	const op = "mission.Repository.DeleteMission"

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	checkQuery, checkArgs, err := r.builder.Select(catIDColumn).From(tableName).Where(sq.Eq{idColumn: id}).ToSql()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	var catID sql.Null[uuid.UUID]
	if err = tx.QueryRow(ctx, checkQuery, checkArgs...).Scan(&catID); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if catID.Valid {
		return utils.ErrCatAssigned
	}

	delQuery, delArgs, err := r.builder.Delete(tableName).Where(sq.Eq{idColumn: id}).ToSql()
	if err != nil {
		tx.Rollback(ctx)
		return fmt.Errorf("%s: %w", op, err)
	}

	if _, err = tx.Exec(ctx, delQuery, delArgs...); err != nil {
		defer tx.Rollback(ctx)
		if errors.Is(err, pgx.ErrNoRows) {
			return utils.ErrMissionNotFound
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (r *Repository) SetMissionCompleted(ctx context.Context, id uuid.UUID) error {
	const op = "mission.Repository.SetMissionCompleted"

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
		return fmt.Errorf("%s: %w", op, utils.ErrMissionNotFound)
	}

	return nil
}

func (r *Repository) AddCatID(ctx context.Context, missionID uuid.UUID, catID uuid.UUID) error {
	const op = "mission.Repository.AddCatID"

	query, args, err := r.builder.
		Update(tableName).
		Set(catIDColumn, catID).
		Where(sq.Eq{idColumn: missionID}).
		ToSql()

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if _, err = r.db.Exec(ctx, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return utils.ErrMissionNotFound
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *Repository) GetMissions(ctx context.Context) ([]*Mission, error) {
	const op = "mission.Repository.GetMissions"
	missions := make([]*Mission, 0)

	query, args, err := r.builder.
		Select(idColumn, stateColumn, createdAtColumn, updatedAtColumn).
		From(tableName).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var rows pgx.Rows
	rows, err = r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	for rows.Next() {
		var mission Mission

		err = rows.Scan(&mission.ID, &mission.State, &mission.CreatedAt, &mission.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		missions = append(missions, &mission)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return missions, nil
}

func (r *Repository) GetMissionByCatID(ctx context.Context, catID uuid.UUID) (*Mission, error) {
	const op = "mission.Repository.GetMissionsByCatID"

	query, args, err := r.builder.
		Select(idColumn, stateColumn, createdAtColumn, updatedAtColumn).
		From(tableName).
		Where(sq.Eq{catIDColumn: catID}).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var mission Mission
	err = r.db.QueryRow(ctx, query, args...).Scan(&mission.ID, &mission.State, &mission.CreatedAt, &mission.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, utils.ErrCatNotFound
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &mission, nil
}

func (r *Repository) GetMissionByID(ctx context.Context, id uuid.UUID) (*Mission, error) {
	const op = "mission.Repository.GetMissionByID"
	var mission Mission

	query, args, err := r.builder.
		Select(idColumn, stateColumn, createdAtColumn, updatedAtColumn).
		From(tableName).
		Where(sq.Eq{idColumn: id}).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	err = r.db.QueryRow(ctx, query, args...).Scan(
		&mission.ID,
		&mission.State,
		&mission.CreatedAt,
		&mission.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, utils.ErrMissionNotFound)
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &mission, nil
}

func (r *Repository) GetAssignedCat(ctx context.Context, missionID uuid.UUID) (*cat.Cat, error) {
	const op = "cat.Repository.GetAssignedCat"
	var assignedCat cat.Cat

	query := fmt.Sprintf(`SELECT c.id, c.name, c.experience, c.breed, c.salary, c.created_at, c.updated_at" +
		"FROM cats c JOIN missions m ON m.cat_id = c.id WHERE m.id = $1`)

	args := []interface{}{
		missionID,
	}

	err := r.db.QueryRow(ctx, query, args...).Scan(
		&assignedCat.ID,
		&assignedCat.Name,
		&assignedCat.YearsXP,
		&assignedCat.Breed,
		&assignedCat.SalaryCents,
		&assignedCat.CreatedAt,
		&assignedCat.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, utils.ErrCatNotFound
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &assignedCat, nil
}
