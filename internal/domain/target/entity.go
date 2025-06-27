package target

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"time"
)

type Target struct {
	ID        uuid.UUID
	Name      string `validate:"required"`
	Country   string `validate:"required"`
	Notes     string
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

func NewEntity(name, country string, notes string) *Target {
	return &Target{
		Name:    name,
		Country: country,
		Notes:   notes,
	}
}

func (t *Target) Validate() error {
	const op = "target.validate"

	validate := validator.New(validator.WithRequiredStructEnabled())

	if err := validate.Struct(t); err != nil {
		return fmt.Errorf("%s: %s", op, err.Error())
	}

	return nil
}
