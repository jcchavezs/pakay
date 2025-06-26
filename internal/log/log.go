package log

import "log/slog"

var Logger = slog.Default()

// SetHandler for the logger
func SetHandler(handler slog.Handler) {
	Logger = slog.New(handler)
}
