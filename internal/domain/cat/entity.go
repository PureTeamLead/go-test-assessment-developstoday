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
	SalaryCents int64  `validate:"min=0 required"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type CreateCatSvc struct {
	Name     string
	YearsExp int
	Breed    string
	Salary   int64
}

type UpdateCatParams struct {
	ID          uuid.UUID
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

func (c *Cat) Update(params UpdateCatParams) {
	if params.Name != nil {
		c.Name = *(params.Name)
	}

	if params.Breed != nil {
		c.Breed = *(params.Breed)
	}

	if params.YearsXP != nil {
		c.YearsXP = *(params.YearsXP)
	}

	if params.SalaryCents != nil {
		c.SalaryCents = *(params.SalaryCents)
	}
}
