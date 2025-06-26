package cat

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/PureTeamLead/go-test-assessment-developstoday/internal/utils"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// TODO: dollar placeholder?

type Repository struct {
	db      *pgxpool.Pool
	builder sq.StatementBuilderType
}

type AddRepoDTO struct {
	Name        string
	YearsXP     int
	Breed       string
	SalaryCents int64
}

const (
	tableName    = "cats"
	idColumn     = "id"
	nameColumn   = "name"
	expColumn    = "experience"
	breedColumn  = "breed"
	salaryColumn = "salary"
	allColumns   = "*"
)

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{db: pool}
}

func (r *Repository) GetCatByID(ctx context.Context, id uuid.UUID) (*Cat, error) {
	const op = "cat.Repository.GetCatByID"
	var cat Cat

	query, args, err := r.builder.
		Select(allColumns).
		From(tableName).
		Where(sq.Eq{idColumn: id}).
		Limit(1).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	err = r.db.QueryRow(ctx, query, args...).Scan(
		&cat.ID,
		&cat.Name,
		&cat.YearsXP,
		&cat.Breed,
		&cat.SalaryCents,
		&cat.CreatedAt,
		&cat.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, utils.ErrCatNotFound)
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &cat, nil
}

func (r *Repository) GetCatByName(ctx context.Context, name string) (*Cat, error) {
	const op = "cat.Repository.GetCatByName"
	var cat Cat

	query, args, err := r.builder.
		Select(allColumns).
		From(tableName).
		Where(sq.Eq{nameColumn: name}).
		Limit(1).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	err = r.db.QueryRow(ctx, query, args...).Scan(
		&cat.ID,
		&cat.Name,
		&cat.YearsXP,
		&cat.Breed,
		&cat.SalaryCents,
		&cat.CreatedAt,
		&cat.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, utils.ErrCatNotFound)
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &cat, nil
}

func (r *Repository) GetCats(ctx context.Context) ([]*Cat, error) {
	const op = "cat.Repository.GetCats"
	cats := make([]*Cat, 0)

	query, args, err := r.builder.
		Select(allColumns).
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
		var cat Cat

		err = rows.Scan(&cat.ID, &cat.Name, &cat.YearsXP, &cat.Breed, &cat.SalaryCents, &cat.CreatedAt, &cat.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		cats = append(cats, &cat)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return cats, nil
}

func (r *Repository) AddCat(ctx context.Context, cat AddRepoDTO) (uuid.UUID, error) {
	const op = "cat.Repository.AddCat"
	var id uuid.UUID

	query, args, err := r.builder.
		Insert(tableName).
		Columns(nameColumn, expColumn, breedColumn, salaryColumn).
		Values(cat.Name, cat.YearsXP, cat.Breed, cat.SalaryCents).
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

func (r *Repository) DeleteCat(ctx context.Context, id uuid.UUID) error {
	const op = "cat.Repository.DeleteCat"

	query, args, err := r.builder.Delete(tableName).Where(sq.Eq{idColumn: id}).ToSql()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if _, err = r.db.Exec(ctx, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return utils.ErrCatNotFound
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *Repository) UpdateCat(ctx context.Context, cat *Cat) error {
	const op = "cat.Repository.UpdateCat"

	query, args, err := r.builder.Update(tableName).
		Set(nameColumn, cat.Name).
		Set(salaryColumn, cat.SalaryCents).
		Set(breedColumn, cat.Breed).
		Set(expColumn, cat.YearsXP).
		Where(sq.Eq{idColumn: cat.ID}).
		ToSql()

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	var res pgconn.CommandTag
	res, err = r.db.Exec(ctx, query, args...)
	if err != nil {
		var pgErr *pgconn.PgError
		if ok := errors.As(err, &pgErr); ok {
			if pgErr.Code == "23505" {
				return fmt.Errorf("%s: %w", op, utils.ErrConflictingData)
			}
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	if res.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", op, utils.ErrCatNotFound)
	}

	return nil
}
