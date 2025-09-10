package entry

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var entries []Entry

func RegisterRoutes(r *gin.Engine) {
	r.POST("/habits/:id/tracking", CreateTracking)
	r.GET("/habits/:id/tracking", GetTrackings)
}

func CreateTracking(c *gin.Context) {
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
}

func GetTrackings(c *gin.Context) {
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
}
