package mission

import (
	"github.com/google/uuid"
	"time"
)

type Mission struct {
	ID        uuid.UUID
	State     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewEntity : here is the empty entity, because all of the values are set on db level
func NewEntity() *Mission {
	return &Mission{}
}
