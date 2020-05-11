package main

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/beldin0/users/src/logging"
)

func gracefulShutdown(wait chan os.Signal, server *http.Server) {
	logger := logging.NewLogger()
	<-wait
	logger.Info("shutting down")
	ctx, _ := context.WithTimeout(context.Background(), time.Second)
	if err := server.Shutdown(ctx); err != nil {
		logger.Fatal("failed shutting down gracefully")
	}
}
