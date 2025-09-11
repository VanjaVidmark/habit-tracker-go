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
	_, err = db.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := habit.InitSchema(db); err != nil {
		log.Fatal(err)
	}
	if err := entry.InitSchema(db); err != nil {
		log.Fatal(err)
	}
	habitStore := &habit.Store{DB: db}
	entryStore := &entry.Store{DB: db}

	// Start server
	server := gin.Default()

	habit.RegisterRoutes(server, habitStore)
	entry.RegisterRoutes(server, entryStore)

	server.Run(":8080")
}
