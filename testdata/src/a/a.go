package a

import (
	"context"
	"log/slog"
)

func TestLowercaseRule() {
	// Rule 1: Messages should start with lowercase
	slog.Info("Starting server on port 8080")   // want "log message should start with lowercase letter"
	slog.Error("Failed to connect to database") // want "log message should start with lowercase letter"

	// Correct usage
	slog.Info("starting server on port 8080")
	slog.Error("failed to connect to database")
}

func TestEnglishOnlyRule() {
	// Rule 2: Messages should be in English only
	slog.Info("Запуск сервера")                    // want "log message should be in English only"
	slog.Error("Ошибка подключения к базе данных") // want "log message should be in English only"

	// Correct usage
	slog.Info("starting server")
	slog.Error("failed to connect to database")
}

func TestEnglishOnlyRuleContext() {
	ctx := context.Background()

	// Rule 2: Messages should be in English only (context methods)
	slog.InfoContext(ctx, "Обработка запроса") // want "log message should be in English only"

	// Correct usage
	slog.InfoContext(ctx, "processing request")
}

func TestSpecialCharsRule() {
	// Rule 3: No special characters or emojis
	slog.Info("server started!🚀")                 // want "log message should not contain emojis"
	slog.Error("connection failed!!!")            // want "log message should not contain excessive punctuation or special characters"
	slog.Warn("warning: something went wrong...") // want "log message should not contain excessive punctuation or special characters"

	// Correct usage
	slog.Info("server started")
	slog.Error("connection failed")
	slog.Warn("something went wrong")
}

func TestSensitiveDataRule() {
	password := "secret123"
	apiKey := "key123"
	token := "token123"

	// Rule 4: No sensitive data
	slog.Info("user password: " + password) // want "log message may contain sensitive data \\(keyword: password\\)"
	slog.Debug("api_key=" + apiKey)         // want "log message may contain sensitive data \\(keyword: api_key\\)"
	slog.Info("token: " + token)            // want "log message may contain sensitive data \\(keyword: token\\)"

	// Correct usage
	slog.Info("user login successful")
	slog.Debug("api request completed")
	slog.Info("session validated")
}

func TestContextMethods() {
	ctx := context.Background()

	// Test context-aware methods
	slog.InfoContext(ctx, "Starting service")   // want "log message should start with lowercase letter"
	slog.ErrorContext(ctx, "Failed to process") // want "log message should start with lowercase letter"

	// Correct usage
	slog.InfoContext(ctx, "starting service")
	slog.ErrorContext(ctx, "failed to process")
}

func TestValidMessages() {
	// All valid messages - should not trigger any warnings
	slog.Info("server started successfully")
	slog.Debug("processing request")
	slog.Warn("connection timeout")
	slog.Error("failed to read file")
}
