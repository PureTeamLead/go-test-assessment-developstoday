package handler

import (
	"context"
	"github.com/PureTeamLead/go-test-assessment-developstoday/internal/domain/cat"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CatService interface {
	CreateCat(ctx context.Context, req cat.CreateCatSvc) (uuid.UUID, error)
	DeleteCat(ctx context.Context, id uuid.UUID) error
	ListCats(ctx context.Context) ([]*cat.Cat, error)
	GetCatByID(ctx context.Context, id uuid.UUID) (*cat.Cat, error)
	GetCatByName(ctx context.Context, name string) (*cat.Cat, error)
	UpdateCat(ctx context.Context, params cat.UpdateCatParams) (*cat.Cat, error)
}

func (h *Handler) GetCats(c *gin.Context) {
	
}

func (h *Handler) CreateCat(c *gin.Context) {

}

func (h *Handler) GetCat(c *gin.Context) {

}

func (h *Handler) UpdateCat(c *gin.Context) {

}

func (h *Handler) DeleteCat(c *gin.Context) {

}
