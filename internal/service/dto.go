package service

import (
	"github.com/PureTeamLead/go-test-assessment-developstoday/internal/domain/cat"
	"github.com/PureTeamLead/go-test-assessment-developstoday/internal/domain/target"
	"github.com/google/uuid"
	"time"
)

type FullMission struct {
	ID        uuid.UUID
	Cat       *cat.Cat
	Targets   []*target.Target
	State     string
	CreatedAt time.Time
	UpdatedAt time.Time
}
