package graceful

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
)

// RunWithShutdown blocks until a specific signal is received.
// It executes cleanup logic before exiting the program entirely.
func RunWithShutdown(logger *zap.Logger, cleanup func()) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Block until signal is received
	sig := <-sigChan
	logger.Info("Received stop signal, initiating graceful shutdown...", zap.String("signal", sig.String()))

	// Run cleanup logic
	cleanup()

	logger.Info("Cleanup finished. Exiting.")
}

// WatchContext returns a context that cancels when OS termination signals are received.
func WatchContext() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		cancel()
	}()

	return ctx, cancel
}
