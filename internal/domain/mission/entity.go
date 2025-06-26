package mission

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"time"
)

type Mission struct {
	ID        uuid.UUID
	CatID     uuid.UUID `validate:"required"`
	State     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewEntity(catID uuid.UUID) *Mission {
	return &Mission{
		CatID: catID,
	}
}

func (m *Mission) Validate() error {
	const op = "mission.validate"

	validate := validator.New(validator.WithRequiredStructEnabled())

	if err := validate.Struct(m); err != nil {
		return fmt.Errorf("%s: %s", op, err.Error())
	}

	return nil
}
