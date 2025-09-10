package habit

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	Store *Store
}

func RegisterRoutes(r *gin.Engine, store *Store) {
	h := &Handler{Store: store}

	r.GET("/habits", h.GetAll)
	r.GET("/habits/:id", h.GetByID)
	r.POST("/habits", h.CreateHabit)
	// r.PATCH("/habits/:id", h.UpdateHabit)
	r.DELETE("/habits/:id", h.DeleteHabit)
}

func (h *Handler) GetAll(c *gin.Context) {
	habits, err := h.Store.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, habits)
}

func (h *Handler) GetByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID in request"})
		return
	}

	habit, err := h.Store.GetById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, habit)
}

func (h *Handler) CreateHabit(c *gin.Context) {
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

	err := h.Store.Create(&newHabit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusCreated, newHabit)
}

/*
func (h *Handler) UpdateHabit(c *gin.Context) {
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
}
*/

func (h *Handler) DeleteHabit(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID"})
		return
	}
	error := h.Store.Delete(id)
	if error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": error.Error()})
		return
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "Habit not found"})
}
