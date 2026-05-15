package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gentle-ai/gentle-ai/internal/config"
	"github.com/gentle-ai/gentle-ai/internal/server"
)

// Version is set at build time via ldflags
var Version = "dev"

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle graceful shutdown on OS signals
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		sig := <-sigCh
		log.Printf("received signal %s, shutting down...", sig)
		cancel()
	}()

	srv, err := server.New(cfg)
	if err != nil {
		log.Fatalf("failed to initialize server: %v", err)
	}

	fmt.Printf("gentle-ai %s starting on %s\n", Version, cfg.ListenAddr)

	if err := srv.Run(ctx); err != nil {
		log.Fatalf("server exited with error: %v", err)
	}

	log.Println("shutdown complete")
}
