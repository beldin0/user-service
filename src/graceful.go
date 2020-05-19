package main

import (
	"context"
	"os"
	"time"

	"github.com/beldin0/users/src/logging"
)

func gracefulShutdown(mainCtx context.Context, wait chan os.Signal, shutdown func(context.Context) error) {
	logger := logging.NewLogger()
	select {
	case <-wait:
	case <-mainCtx.Done():
	}
	logger.Info("shutting down")
	ctx, _ := context.WithTimeout(context.Background(), time.Second)
	if err := shutdown(ctx); err != nil {
		logger.Fatal("failed shutting down gracefully")
	}
}
