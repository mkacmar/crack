package cli

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func cancelOnSignal(ctx context.Context) context.Context {
	ctx, cancel := context.WithCancel(ctx)
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		select {
		case <-sigCh:
			fmt.Fprintln(os.Stderr, "\nInterrupted, shutting down.")
			cancel()
		case <-ctx.Done():
		}
		signal.Stop(sigCh)
	}()
	return ctx
}
