package dto

type CreateCatInput struct {
	Name              string `json:"name"`
	Breed             string `json:"breed"`
	ExperienceInYears int    `json:"years_exp"`
	SalaryCents       int64  `json:"salary_cents"`
}

type UpdateCatInput struct {
	Name              *string `json:"name,omitempty"`
	Breed             *string `json:"breed,omitempty"`
	ExperienceInYears *int    `json:"years_exp,omitempty"`
	SalaryCents       *int64  `json:"salary_cents,omitempty"`
}
