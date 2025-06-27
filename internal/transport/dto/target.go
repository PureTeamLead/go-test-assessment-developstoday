package dto

import "github.com/PureTeamLead/go-test-assessment-developstoday/internal/domain/target"

type CreateTargetReq struct {
	Name    string  `json:"name"`
	Country string  `json:"country"`
	Notes   *string `json:"notes,omitempty"`
}

func MapTargetsToRaw(reqTargets []CreateTargetReq) []target.CreateUpdateTargetSvc {
	targets := make([]target.CreateUpdateTargetSvc, 0, len(reqTargets))

	for _, rawTarget := range reqTargets {
		tar := target.CreateUpdateTargetSvc{Name: rawTarget.Name, Country: rawTarget.Country, Notes: rawTarget.Notes}

		targets = append(targets, tar)
	}

	return targets
}
