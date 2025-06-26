package cat

import (
	"context"
	"github.com/PureTeamLead/go-test-assessment-developstoday/internal/transport/dto"
	"github.com/google/uuid"
)

type Repo interface {
	GetCatByID(ctx context.Context, id uuid.UUID) (*Cat, error)
	GetCatByName(ctx context.Context, name string) (*Cat, error)
	GetCats(ctx context.Context) ([]*Cat, error)
	AddCat(ctx context.Context, cat AddRepoDTO) (uuid.UUID, error)
	DeleteCat(ctx context.Context, id uuid.UUID) error
	UpdateCat(ctx context.Context, id uuid.UUID, params UpdateCatParams) (*Cat, error)
}

type Service struct {
	repo Repo
}

func NewService(repo Repo) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateCat(ctx context.Context, req dto.CreateCatReq) (uuid.UUID, error) {
	// validate breed via API
}

func (s *Service) DeleteCat(ctx context.Context) {

}

func (s *Service) ListCats(ctx context.Context) ([]*Cat, error) {

}

func (s *Service) GetCatByID(ctx context.Context, id uuid.UUID) (*Cat, error) {

}

func (s *Service) GetCatByName(ctx context.Context, name string) (*Cat, error) {

}

func (s *Service) UpdateCat(ctx context.Context, params UpdateCatParams) (*Cat, error) {

}
