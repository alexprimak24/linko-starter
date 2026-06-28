package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

// approach of returning error and avoiding log.Fatal
// is better for unit testing + cleanups
// as if we fail there and just return, how we
// are supposed to do some cleaunups in main
func initializeLogger(logFile string) (logger *log.Logger, err error) {
	// if logfile exists, write to file too
	if logFile != "" {
		file, err := os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0o644)
		if err != nil {
			return nil, fmt.Errorf("failed to open log file: %v", err)
		}
		bufferedFile := bufio.NewWriterSize(file, 8192)

		multiWriter := io.MultiWriter(os.Stderr, bufferedFile)

		return log.New(multiWriter, "", log.LstdFlags), nil
	}

	return log.New(os.Stderr, "", log.LstdFlags), nil
}

func requestLogger(logger *log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
			logger.Printf("Served request: %s %s", r.Method, r.URL.Path)
		})
	}
}
