package entry

import (
	"time"

	"github.com/google/uuid"
)

type Entry struct {
	ID        uuid.UUID
	HabitId   uuid.UUID
	Timestamp time.Time
	Note      string
}
