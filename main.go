package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

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

type Habit struct {
	ID          uuid.UUID
	Name        string
	Description string
	Frequency   Frequency
	StartDate   time.Time
}

type Entry struct {
	ID        uuid.UUID
	HabitId   uuid.UUID
	Timestamp time.Time
	Note      string
}

func main() {
	server := gin.Default()

	// Temp storage of habits & entries (slice = dynamic size array)
	var habits = []Habit{
		{ID: uuid.New(), Name: "Running", Description: "", Frequency: Daily, StartDate: time.Date(2025, 8, 8, 0, 0, 0, 0, time.UTC)},
	}
	var entries = []Entry{
		{ID: uuid.New(), HabitId: habits[0].ID, Timestamp: time.Date(2025, 9, 9, 0, 0, 0, 0, time.UTC), Note: ""},
	}

	// Get all habits
	server.GET("/habits", func(c *gin.Context) {
		c.JSON(http.StatusOK, habits)
	})

	// Get habit by id
	server.GET("/habits/:id", func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := uuid.Parse(idParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID in request"})
			return
		}
		for _, habit := range habits {
			if habit.ID == id {
				c.JSON(http.StatusOK, habit)
				return
			}
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "No habit with matching ID found"})
	})

	// Add a new habit
	server.POST("/habits", func(c *gin.Context) {
		var newHabit Habit

		// BindJSON binds the received JSON to newHabit
		if err := c.BindJSON(&newHabit); err != nil {
			return
		}

		// Validate frequency
		if !newHabit.Frequency.IsValid() {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid frequency, must be Daily, Weekly, or Monthly"})
			return
		}

		newHabit.ID = uuid.New()
		newHabit.StartDate = time.Now()

		habits = append(habits, newHabit)
		c.IndentedJSON(http.StatusCreated, newHabit)
	})

	// Updates habit
	server.PATCH("/habits/:id", func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := uuid.Parse(idParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID"})
			return
		}

		var update Habit
		if err := c.ShouldBindJSON(&update); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
			return
		}
		log.Println(update)

		for i := range habits {
			if habits[i].ID == id {
				if update.Name != "" {
					habits[i].Name = update.Name
				}
				if update.Description != "" {
					habits[i].Description = update.Description
				}
				if update.Frequency != "" {
					if !update.Frequency.IsValid() {
						c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid frequency"})
						return
					}
					habits[i].Frequency = update.Frequency
				}
				if !update.StartDate.IsZero() {
					habits[i].StartDate = update.StartDate
				}

				c.JSON(http.StatusOK, habits[i])
				return
			}
		}

		c.JSON(http.StatusNotFound, gin.H{"error": "Habit not found"})
	})

	server.DELETE("/habits/:id", func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := uuid.Parse(idParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID"})
			return
		}
		for i, habit := range habits {
			if habit.ID == id {
				habits = append(habits[:i], habits[i+1:]...)
				c.Status(http.StatusNoContent)
				return
			}
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "Habit not found"})
	})

	// Add a new entry
	server.POST("/habits/:id/tracking", func(c *gin.Context) {

		idParam := c.Param("id")
		id, err := uuid.Parse(idParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID"})
			return
		}

		// Insert validation that habit is valid?

		var newEntry Entry

		// BindJSON binds the received JSON to newEntry
		if err := c.BindJSON(&newEntry); err != nil {
			return
		}

		newEntry.ID = uuid.New()
		newEntry.HabitId = id
		if newEntry.Timestamp.IsZero() {
			newEntry.Timestamp = time.Now()
		}

		entries = append(entries, newEntry)
		c.IndentedJSON(http.StatusCreated, newEntry)
	})

	// Get all trackings of a habit
	server.GET("/habits/:id/tracking", func(c *gin.Context) {

		idParam := c.Param("id")
		id, err := uuid.Parse(idParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID"})
			return
		}

		// Insert validation that habit is valid?

		var entriesToReturn []Entry

		for _, entry := range entries {
			if entry.HabitId == id {
				entriesToReturn = append(entriesToReturn, entry)
			}
		}

		c.JSON(http.StatusOK, entriesToReturn)
	})

	server.Run(":8080")
}
