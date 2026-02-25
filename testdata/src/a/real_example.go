package a

import (
	"context"
	"log/slog"
	"net/http"
)

// Example of a real HTTP server with logging
func StartServer() error {
	slog.Info("Starting HTTP server on :8080") // want "log message should start with lowercase letter"

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		slog.Info("handling request")
		w.Write([]byte("OK"))
	})

	slog.Info("server is ready")
	return http.ListenAndServe(":8080", nil)
}

// Example of database connection with logging
func ConnectDB(password string) error {
	slog.Info("connecting to database")

	// Bad: logging sensitive data
	slog.Debug("connection string with password: " + password) // want "log message may contain sensitive data \\(keyword: password\\)"

	slog.Info("database connected successfully")
	return nil
}

// Example with context
func ProcessRequest(ctx context.Context, userID string) {
	slog.InfoContext(ctx, "processing user request")

	// Bad: Cyrillic text
	slog.InfoContext(ctx, "обработка запроса") // want "log message should be in English only"

	slog.InfoContext(ctx, "request processed successfully")
}

// Example with emojis
func NotifyUser() {
	slog.Info("sending notification 📧") // want "log message should not contain emojis"
	slog.Info("notification sent!!!")   // want "log message should not contain excessive punctuation or special characters"
}
