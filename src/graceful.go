package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"
)

func gracefulShutdown(wait chan os.Signal, server *http.Server) {
	<-wait
	log.Println("shutting down")
	ctx, _ := context.WithTimeout(context.Background(), time.Second)
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("failed shutting down gracefully")
	}
}
