package service

import (
	"github.com/PureTeamLead/go-test-assessment-developstoday/internal/domain/cat"
	"github.com/PureTeamLead/go-test-assessment-developstoday/internal/domain/target"
	"github.com/google/uuid"
	"time"
)

// FullMission aggregated structure
type FullMission struct {
	ID        uuid.UUID
	Cat       *cat.Cat
	Targets   []*target.Target
	State     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// CreateUpdateTargetSvc Notes field is optional
type CreateUpdateTargetSvc struct {
	Name    string
	Country string
	Notes   *string
}

func MapTargetSvcToEntity(tarReq CreateUpdateTargetSvc) *target.Target {
	var tar target.Target
	tar.Name = tarReq.Name
	tar.Country = tarReq.Country

	if tarReq.Notes != nil {
		tar.Notes = *(tarReq.Notes)
	}

	return &tar
}
