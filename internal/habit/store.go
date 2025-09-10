package habit

import (
	"database/sql"
	"log"
	"time"

	"github.com/google/uuid"
)

type Store struct {
	DB *sql.DB
}

func InitSchema(db *sql.DB) {
	db, err := sql.Open("sqlite3", "./habits.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create tables if not exist
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS habits (
			id TEXT PRIMARY KEY,
			name TEXT,
			description TEXT,
			frequency TEXT,
			start_date DATETIME
		);
		CREATE TABLE IF NOT EXISTS entries (
			id TEXT PRIMARY KEY,
			habit_id TEXT,
			timestamp DATETIME,
			note TEXT
		);
	`)
	if err != nil {
		log.Fatal(err)
	}
}

func (s *Store) GetAll() ([]Habit, error) {
	rows, err := s.DB.Query("SELECT id, name, description, frequency, start_date FROM habits")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var habits []Habit
	for rows.Next() {
		var h Habit
		var start string
		if err := rows.Scan(&h.ID, &h.Name, &h.Description, &h.Frequency, &start); err != nil {
			return nil, err
		}
		h.StartDate, _ = time.Parse(time.RFC3339, start)
		habits = append(habits, h)
	}
	return habits, nil
}

func (s *Store) GetById(id uuid.UUID) (Habit, error) {
	var h Habit
	var start string

	err := s.DB.QueryRow(
		"SELECT id, name, description, frequency, start_date FROM habits WHERE id = ?",
		id.String(),
	).Scan(&h.ID, &h.Name, &h.Description, &h.Frequency, &start)

	if err != nil {
		return Habit{}, err
	}

	h.StartDate, _ = time.Parse(time.RFC3339, start)
	return h, nil
}

func (s *Store) Create(h *Habit) error {
	h.ID = uuid.New()
	h.StartDate = time.Now()
	_, err := s.DB.Exec(
		"INSERT INTO habits (id, name, description, frequency, start_date) VALUES (?, ?, ?, ?, ?)",
		h.ID.String(), h.Name, h.Description, h.Frequency, h.StartDate.Format(time.RFC3339),
	)
	return err
}

func (s *Store) Delete(id uuid.UUID) error {
	result, err := s.DB.Exec("DELETE FROM habits WHERE id = ?", id.String())
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
