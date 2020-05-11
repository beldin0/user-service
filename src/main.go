package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/beldin0/users/src/logging"
	"github.com/beldin0/users/src/routing"
	"github.com/jmoiron/sqlx"
	"github.com/kelseyhightower/envconfig"
	_ "github.com/lib/pq"
)

const defaultPort = 8080

func main() {
	logger := logging.NewLogger()
	var c config
	err := envconfig.Process("", &c)
	db, err := sqlx.Connect("postgres", c.ConnString())
	if err != nil {
		logger.Sugar().With("error", err).Fatal("problem connecting to database")
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
	logger.Sugar().With("port", defaultPort).Info("listening")
	err = server.ListenAndServe()
	if err != http.ErrServerClosed {
		logger.Sugar().With("error", err).Fatal("problem with server")
	}
}
