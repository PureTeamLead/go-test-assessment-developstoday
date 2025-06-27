package dto

type CreateMissionReq struct {
	Targets []CreateTargetReq `json:"targets"`
}

type UpdateMissionReq struct {
	Notes string `json:"notes"`
}

type AssignToMissionReq struct {
	CatID string `json:"cat_id"`
}
