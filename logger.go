package main

import (
	"io"
	"log"
	"net/http"
	"os"
)

func initializeLogger(logFile string) (logger *log.Logger) {
	// if logfile exists, write to file too
	if logFile != "" {
		file, err := os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0o644)
		if err != nil {
			log.Fatalf("failed to open log file: %v", err)
		}
		multiWriter := io.MultiWriter(os.Stderr, file)

		return log.New(multiWriter, "", log.LstdFlags)
	}

	return log.New(os.Stderr, "", log.LstdFlags)
}

func requestLogger(logger *log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
			logger.Printf("Served request: %s %s", r.Method, r.URL.Path)
		})
	}
}
