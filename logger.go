package main

import (
	"bufio"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
)

type cleanupFunc func() error

// approach of returning error and avoiding log.Fatal
// is better for unit testing + cleanups
// as if we fail there and just return, how we
// are supposed to do some cleaunups in main
func initializeLogger(logFile string) (logger *slog.Logger, cleanup cleanupFunc, err error) {
	handlers := []slog.Handler{
		slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}),
	}

	cleanups := []cleanupFunc{}
	// if logfile exists, write to file too
	if logFile != "" {
		file, err := os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0o644)
		if err != nil {
			// cleanup nil as file is probably nil
			return nil, nil, fmt.Errorf("failed to open log file: %v", err)
		}
		bufferedFile := bufio.NewWriterSize(file, 8192)

		cleanup := func() error {
			if err := bufferedFile.Flush(); err != nil {
				return fmt.Errorf("failed to flush log file: %w", err)
			}
			if err := file.Close(); err != nil {
				return fmt.Errorf("failed to close log file: %w", err)
			}
			return nil
		}

		handlers = append(handlers, slog.NewTextHandler(bufferedFile, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}))
		cleanups = append(cleanups, cleanup)
	}

	cleanupFunc := func() error {
		var errs []error
		for _, cleanup := range cleanups {
			if err := cleanup(); err != nil {
				errs = append(errs, err)
			}
		}

		return errors.Join(errs...)
	}

	return slog.New(slog.NewMultiHandler(handlers...)), cleanupFunc, nil
}

func requestLogger(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
			logger.Info(fmt.Sprintf("Served request: %s %s", r.Method, r.URL.Path))
		})
	}
}
