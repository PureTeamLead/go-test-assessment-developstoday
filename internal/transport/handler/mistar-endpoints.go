package handler

import (
	"context"
	"github.com/PureTeamLead/go-test-assessment-developstoday/internal/domain/target"
	"github.com/PureTeamLead/go-test-assessment-developstoday/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type MisTargetService interface {
	CreateMission(ctx context.Context, rawTargets []target.CreateUpdateTargetSvc) (uuid.UUID, error)
	DeleteMission(ctx context.Context, id uuid.UUID) error
	UpdateMissionState(ctx context.Context, id uuid.UUID) error
	SetMissionTargetState(ctx context.Context, missionID uuid.UUID, targetID uuid.UUID) error
	UpdateMissionTargetNotes(ctx context.Context, missionID uuid.UUID, targetID uuid.UUID, notes string) error
	DeleteTargetFromMission(ctx context.Context, id uuid.UUID) error
	AddTargetToMission(ctx context.Context, targetID uuid.UUID, missionID uuid.UUID) error
	AssignCatToMission(ctx context.Context, missionID uuid.UUID, catID uuid.UUID) error
	ListMissions(ctx context.Context) ([]*service.FullMission, error)
	GetMission(ctx context.Context, id uuid.UUID) (*service.FullMission, error)
}

func (h *Handler) CreateMission(c *gin.Context) {

}

func (h *Handler) DeleteMission(c *gin.Context) {

}

func (h *Handler) UpdateMission(c *gin.Context) {

}

func (h *Handler) UpdateMissionTarget(c *gin.Context) {

}

func (h *Handler) DeleteMissionTarget(c *gin.Context) {

}

func (h *Handler) AddMissionTarget(c *gin.Context) {

}

func (h *Handler) AssignMission(c *gin.Context) {

}

func (h *Handler) ListMissions(c *gin.Context) {

}

func (h *Handler) GetMission(c *gin.Context) {

}
