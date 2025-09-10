package habit

import (
	"time"

	"github.com/google/uuid"
)

type Habit struct {
	ID          uuid.UUID
	Name        string
	Description string
	Frequency   Frequency
	StartDate   time.Time
}

type Frequency string

const (
	Daily   Frequency = "Daily"
	Weekly  Frequency = "Weekly"
	Monthly Frequency = "Monthly"
)

func (f Frequency) IsValid() bool {
	switch f {
	case Daily, Weekly, Monthly:
		return true
	}
	return false
}
