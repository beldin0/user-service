package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/beldin0/users/src/routing"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

const defaultPort = 8080

func main() {
	db, err := sqlx.Connect("sqlite3", ":memory:")
	if err != nil {
		log.Fatal(err)
	}
	server := &http.Server{
		Addr:    fmt.Sprint(":", defaultPort),
		Handler: routing.NewRouting(db),
	}

	// Prepare for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	go gracefulShutdown(quit, server)

	// Start server
	log.Println("listening on port", defaultPort)
	err = server.ListenAndServe()
	if err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
