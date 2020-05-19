package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/beldin0/users/src/logging"
	pb "github.com/beldin0/users/src/user"
	"github.com/beldin0/users/src/userhandler"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/jmoiron/sqlx"
	"github.com/kelseyhightower/envconfig"
	_ "github.com/lib/pq"
)

const defaultPort = 8080
const grpcPort = 9000

func main() {
	logger := logging.NewLogger()
	var c config
	err := envconfig.Process("", &c)
	db, err := sqlx.Connect("postgres", c.ConnString())
	if err != nil {
		logger.Sugar().
			With("error", err).
			With("connection_string", c.ConnString()).
			Fatal("problem connecting to database")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if run(ctx, db) != http.ErrServerClosed {
		logger.Sugar().With("error", err).Fatal("problem with server")
	}
}

func run(ctx context.Context, db *sqlx.DB) error {
	logger := logging.NewLogger()

	mux := runtime.NewServeMux()
	err := pb.RegisterUserServiceHandlerServer(ctx, mux, userhandler.New(db))
	if err != nil {
		return err
	}

	server := &http.Server{
		Addr:    fmt.Sprint(":", defaultPort),
		Handler: mux,
	}

	// Prepare for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	go gracefulShutdown(ctx, quit, server.Shutdown)

	// Start server
	logger.Sugar().With("port", defaultPort).Info("listening http")
	logger.Sugar().With("port", grpcPort).Info("listening grpc")
	return server.ListenAndServe()
}
