package entry

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

	r.POST("/habits/:id/tracking", h.CreateTracking)
	r.GET("/habits/:id/tracking", h.GetTrackings)
}

func (h *Handler) CreateTracking(c *gin.Context) {
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

	if err := h.Store.Create(&newEntry); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusCreated, newEntry)
}

func (h *Handler) GetTrackings(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID"})
		return
	}

	// Insert validation that habit is valid?

	entries, err := h.Store.GetByHabitId(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch entries"})
	}

	c.JSON(http.StatusOK, entries)
}
