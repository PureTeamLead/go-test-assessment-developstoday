package cat

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"time"
)

type Cat struct {
	Name      string `validate:"required"`
	YearsXP   int    `validate:"required"`
	Breed     string `validate:"required"`
	Salary    int64  `validate:"required"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewEntity(name string, experience int, breed string, salary int64) *Cat {
	// CreatedAt and UpdatedAt are assigned for data consistency -> see trigger function in the database migrations
	return &Cat{
		Name:      name,
		YearsXP:   experience,
		Breed:     breed,
		Salary:    salary,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
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
