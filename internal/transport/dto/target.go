package dto

import (
	"github.com/PureTeamLead/go-test-assessment-developstoday/internal/service"
)

type CreateTargetReq struct {
	Name    string  `json:"name"`
	Country string  `json:"country"`
	Notes   *string `json:"notes,omitempty"`
}

func MapTargetsToRaw(reqTargets []CreateTargetReq) []service.CreateUpdateTargetSvc {
	targets := make([]service.CreateUpdateTargetSvc, 0, len(reqTargets))

	for _, rawTarget := range reqTargets {
		tar := MapTargetToRaw(rawTarget)
		targets = append(targets, tar)
	}
	return targets
}

func MapTargetToRaw(reqTarget CreateTargetReq) service.CreateUpdateTargetSvc {
	return service.CreateUpdateTargetSvc{Name: reqTarget.Name, Country: reqTarget.Country, Notes: reqTarget.Notes}
}
