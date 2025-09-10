package main

import (
	"log"

	"database/sql"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"

	"github.com/VanjaVidmark/habit-tracker-go/internal/entry"
	"github.com/VanjaVidmark/habit-tracker-go/internal/habit"
)

func main() {
	// Open db connection
	db, err := sql.Open("sqlite3", "./habits.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	habit.InitSchema(db)
	store := &habit.Store{DB: db}

	// Start server
	server := gin.Default()

	habit.RegisterRoutes(server, store)
	entry.RegisterRoutes(server)

	server.Run(":8080")
}
