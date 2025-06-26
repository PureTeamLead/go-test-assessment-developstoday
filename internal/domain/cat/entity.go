package cat

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"time"
)

type Cat struct {
	ID          uuid.UUID
	Name        string `validate:"required"`
	YearsXP     int    `validate:"required"`
	Breed       string `validate:"required"`
	SalaryCents int64  `validate:"required"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// TODO: validate update params in service
type UpdateCatParams struct {
	Name        *string
	YearsXP     *int
	Breed       *string
	SalaryCents *int64
}

func NewEntity(name string, experience int, breed string, salary int64) *Cat {
	return &Cat{
		Name:        name,
		YearsXP:     experience,
		Breed:       breed,
		SalaryCents: salary,
	}
}

func (c *Cat) Validate() error {
	const op = "cat.validate"

	validate := validator.New(validator.WithRequiredStructEnabled())

	if err := validate.Struct(c); err != nil {
		return fmt.Errorf("%s: %s", op, err.Error())
	}

	return nil
}
