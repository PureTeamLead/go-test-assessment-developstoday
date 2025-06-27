package cat

import (
	"context"
	"fmt"
	"github.com/PureTeamLead/go-test-assessment-developstoday/internal/api/breedapi"
	"github.com/PureTeamLead/go-test-assessment-developstoday/internal/utils"
	"github.com/google/uuid"
)

type Repo interface {
	GetCatByID(ctx context.Context, id uuid.UUID) (*Cat, error)
	GetCatByName(ctx context.Context, name string) (*Cat, error)
	GetCats(ctx context.Context) ([]*Cat, error)
	AddCat(ctx context.Context, cat *Cat) (uuid.UUID, error)
	DeleteCat(ctx context.Context, id uuid.UUID) error
	UpdateCat(ctx context.Context, cat *Cat) error
}

type Service struct {
	repo Repo
}

func NewService(repo Repo) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateCat(ctx context.Context, req CreateCatSvc) (uuid.UUID, error) {
	const op = "cat.Service.CreateCat"

	// validate breed via API
	if err := breedapi.ValidateBreed(req.Breed); err != nil {
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}

	cat := NewEntity(req.Name, req.YearsExp, req.Breed, req.Salary)
	if err := cat.Validate(); err != nil {
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}

	id, err := s.repo.AddCat(ctx, cat)
	if err != nil {
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *Service) DeleteCat(ctx context.Context, id uuid.UUID) error {
	const op = "cat.Service.CreateCat"

	if err := s.repo.DeleteCat(ctx, id); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Service) ListCats(ctx context.Context) ([]*Cat, error) {
	const op = "cat.Service.ListCats"

	cats, err := s.repo.GetCats(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return cats, nil
}

func (s *Service) GetCatByID(ctx context.Context, id uuid.UUID) (*Cat, error) {
	const op = "cat.Service.GetCatByID"

	cat, err := s.repo.GetCatByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return cat, nil
}

func (s *Service) GetCatByName(ctx context.Context, name string) (*Cat, error) {
	const op = "cat.Service.GetCatByName"

	cat, err := s.repo.GetCatByName(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return cat, nil
}

func (s *Service) UpdateCat(ctx context.Context, params UpdateCatParams) (*Cat, error) {
	const op = "cat.Service.UpdateCat"

	cat, err := s.repo.GetCatByID(ctx, params.ID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if params.Breed != nil {
		err = breedapi.ValidateBreed(*(params.Breed))
		if err != nil {
			return nil, utils.ErrInvalidBreed
		}
	}

	cat.Update(params)
	if err = cat.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	err = s.repo.UpdateCat(ctx, cat)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return cat, nil
}
