package models

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"os/signal"
)

var database *sql.DB

func init() {

	var err error
	fmt.Println("Connecting to database...")
	database, err = sql.Open("sqlite3", "./resources/MessageBoard.db")
	if err != nil {
		panic(err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go awaitShutdown(c)

}

func awaitShutdown(c chan os.Signal) {
	for range c {
		fmt.Println("Disconnecting from to database...")
		database.Close()
		os.Exit(0)
	}
}
