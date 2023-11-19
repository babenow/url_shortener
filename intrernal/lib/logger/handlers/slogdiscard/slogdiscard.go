package slogdiscard

import (
	"context"
	"log/slog"
)

func NewSlogDiscardLogger() *slog.Logger {
	return slog.New(NewDiscardHandler())
}

type DiscardHandler struct {
}

func NewDiscardHandler() *DiscardHandler {
	return &DiscardHandler{}
}

func (h *DiscardHandler) Enabled(context.Context, slog.Level) bool {
	return true
}

func (h *DiscardHandler) Handle(context.Context, slog.Record) error {
	return nil
}

func (h *DiscardHandler) WithGroup(string) slog.Handler {
	return h
}

func (h *DiscardHandler) WithAttrs([]slog.Attr) slog.Handler {
	return h
}
