package entry

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Store struct {
	DB *sql.DB
}

func InitSchema(db *sql.DB) error {
	_, err := db.Exec(`
        CREATE TABLE IF NOT EXISTS entries (
            id TEXT PRIMARY KEY,
			habit_id TEXT,
            timestamp DATETIME,
			note TEXT,
			FOREIGN KEY (habit_id) REFERENCES habits(id) ON DELETE CASCADE
        );
    `)
	return err
}

func (s *Store) Create(e *Entry) error {
	_, err := s.DB.Exec(
		"INSERT INTO entries (id, habit_id, timestamp, note) VALUES (?, ?, ?, ?)",
		e.ID.String(), e.HabitId.String(), e.Timestamp.Format(time.RFC3339), e.Note,
	)
	return err
}

func (s *Store) GetByHabitId(habitID uuid.UUID) ([]Entry, error) {
	rows, err := s.DB.Query("SELECT id, habit_id, timestamp, note FROM entries WHERE habit_id = ?", habitID.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []Entry
	for rows.Next() {
		var e Entry
		var start string
		if err := rows.Scan(&e.ID, &e.HabitId, &e.Timestamp, &e.Note); err != nil {
			return nil, err
		}
		e.Timestamp, _ = time.Parse(time.RFC3339, start)
		entries = append(entries, e)
	}
	return entries, nil
}
