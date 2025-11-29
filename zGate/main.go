package main

import (
	"context"
	"encoding/hex"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"golang.org/x/sync/errgroup"

	"github.com/zGate-Team/zGate-Platform/internal/api"
	"github.com/zGate-Team/zGate-Platform/internal/store"
	"github.com/zGate-Team/zGate-Platform/internal/utils"
)

const (
	shutdownTimeout = 30 * time.Second
	defaultDBPath   = "data/zgate.db"
	storePathEnvVar = "ZGATE_STORE_PATH"
	storeKeyEnvVar  = "ZGATE_STORE_KEY"
	portEnvVar      = "ZGATE_PORT"
)

func main() {
	if err := loadDotEnv(); err != nil {
		utils.Logger.Warn("failed to load .env file", "error", err)
	}

	// Get default port from environment variable
	defaultPort := os.Getenv(portEnvVar)
	if defaultPort == "" {
		fmt.Fprintf(os.Stderr, "%s environment variable is required\n", portEnvVar)
		os.Exit(1)
	}
	if defaultPort[0] != ':' {
		// Ensure port has colon prefix if not provided
		defaultPort = ":" + defaultPort
	}

	// Parse command-line flags
	apiAddr := flag.String("api-addr", defaultPort, "address for API server")
	flag.Parse()

	// Initialize structured logger
	if err := utils.InitLogger(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to initialize logger: %v\n", err)
		os.Exit(1)
	}

	utils.Logger.Info("starting zGate platform")

	dataStore, err := initStore()
	if err != nil {
		utils.Logger.Error("failed to initialize store", "error", err)
		os.Exit(1)
	}
	defer func() {
		if err := dataStore.Close(); err != nil {
			utils.Logger.Warn("failed to close store", "error", err)
		}
	}()

	logStoreInventory(dataStore)

	// Initialize API server (includes proxy manager)
	apiServer, err := api.NewServer(*apiAddr, dataStore)
	if err != nil {
		utils.Logger.Error("failed to initialize API server", "error", err)
		os.Exit(1)
	}

	// Set up signal handling for graceful shutdown
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// Use errgroup to manage all services
	g, gctx := errgroup.WithContext(ctx)

	// Start API server
	g.Go(func() error {
		if err := apiServer.Start(); err != nil {
			utils.Logger.Error("API server failed", "error", err)
			return fmt.Errorf("API server: %w", err)
		}
		return nil
	})

	// Background task: Clean up expired tokens periodically
	g.Go(func() error {
		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()

		utils.Logger.Info("token cleanup task started")

		for {
			select {
			case <-gctx.Done():
				utils.Logger.Info("token cleanup task stopped")
				return nil
			case <-ticker.C:
				if err := dataStore.CleanupExpiredTokens(); err != nil {
					utils.Logger.Warn("failed to cleanup expired tokens", "error", err)
				} else {
					utils.Logger.Debug("expired tokens cleaned up")
				}
			}
		}
	})

	// Gracefully shutdown API server when context is canceled
	g.Go(func() error {
		<-gctx.Done()

		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer shutdownCancel()

		if err := apiServer.Shutdown(shutdownCtx); err != nil {
			utils.Logger.Error("API server shutdown error", "error", err)
			return fmt.Errorf("API server shutdown: %w", err)
		}
		utils.Logger.Info("API server stopped")
		return nil
	})

	// Wait for all services to complete or error
	if err := g.Wait(); err != nil {
		utils.Logger.Error("application stopped with error", "error", err)
		os.Exit(1)
	}

	utils.Logger.Info("zGate shut down successfully")
}

func loadDotEnv() error {
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

func initStore() (*store.Store, error) {
	keyHex := os.Getenv(storeKeyEnvVar)
	if keyHex == "" {
		return nil, fmt.Errorf("%s environment variable is required", storeKeyEnvVar)
	}

	key, err := hex.DecodeString(keyHex)
	if err != nil {
		return nil, fmt.Errorf("failed to decode %s as hex: %w", storeKeyEnvVar, err)
	}
	if len(key) != 32 {
		return nil, fmt.Errorf("decoded %s must be exactly 32 bytes, got %d", storeKeyEnvVar, len(key))
	}

	dbPath := os.Getenv(storePathEnvVar)
	if dbPath == "" {
		dbPath = defaultDBPath
	}

	utils.Logger.Info("opening metadata store", "path", dbPath)
	return store.NewStore(dbPath, key)
}

func logStoreInventory(s *store.Store) {
	databases, err := s.ListDatabases()
	if err != nil {
		utils.Logger.Warn("failed to list databases", "error", err)
	}
	roles, err := s.ListRoles()
	if err != nil {
		utils.Logger.Warn("failed to list roles", "error", err)
	}
	users, err := s.ListUsers()
	if err != nil {
		utils.Logger.Warn("failed to list users", "error", err)
	}

	utils.Logger.Info("store ready",
		"databases", len(databases),
		"roles", len(roles),
		"users", len(users),
	)
}
