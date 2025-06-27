package dto

type CreateMissionReq struct {
	Targets []CreateTargetReq `json:"targets"`
}
