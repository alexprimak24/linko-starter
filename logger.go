package main

import (
	"bufio"
	"fmt"
	"io"
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
	// if logfile exists, write to file too
	if logFile != "" {
		file, err := os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0o644)
		if err != nil {
			// cleanup nil as file is probably nil
			return nil, nil, fmt.Errorf("failed to open log file: %v", err)
		}
		bufferedFile := bufio.NewWriterSize(file, 8192)

		multiWriter := io.MultiWriter(os.Stderr, bufferedFile)

		return slog.New(slog.NewTextHandler(multiWriter, nil)), func() error {
			flushErr := bufferedFile.Flush()
			closeErr := file.Close()

			if flushErr != nil {
				return flushErr
			}

			return closeErr
		}, nil
	}

	return slog.New(slog.NewTextHandler(os.Stderr, nil)), func() error {
		return nil
	}, nil
}

func requestLogger(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
			logger.Info(fmt.Sprintf("Served request: %s %s", r.Method, r.URL.Path))
		})
	}
}
