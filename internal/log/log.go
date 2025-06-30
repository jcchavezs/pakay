package log

import (
	"context"
	"log/slog"
)

var Logger = slog.Default()

// SetHandler for the logger
func SetHandler(handler slog.Handler) {
	Logger = slog.New(handler)
}

// Ported from slog.DiscardHandler which makes this available from Go 1.24
var DiscardHandler slog.Handler = discardHandler{}

type discardHandler struct{}

func (dh discardHandler) Enabled(context.Context, slog.Level) bool  { return false }
func (dh discardHandler) Handle(context.Context, slog.Record) error { return nil }
func (dh discardHandler) WithAttrs(attrs []slog.Attr) slog.Handler  { return dh }
func (dh discardHandler) WithGroup(name string) slog.Handler        { return dh }
