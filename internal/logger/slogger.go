package logger

import (
	"log/slog"
	"os"
)

type Options struct {
	Level       int
	Structured  bool
	ServiceName string
	Version     string
}

func Init(opts Options) *slog.Logger {
	log := setupLogger(opts)
	slog.SetDefault(log)
	return log
}

func setupLogger(opts Options) *slog.Logger {
	var level slog.Level
	switch opts.Level {
	case -4:
		level = slog.LevelDebug
	case 0:
		level = slog.LevelInfo
	case 4:
		level = slog.LevelWarn
	case 8:
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	handlerOpts := &slog.HandlerOptions{
		Level:     level,
		AddSource: level == slog.LevelDebug,
	}

	var handler slog.Handler
	if opts.Structured {
		handler = slog.NewJSONHandler(os.Stdout, handlerOpts)
	} else {
		handler = slog.NewTextHandler(os.Stdout, handlerOpts)
	}

	var attrs []slog.Attr
	if opts.ServiceName != "" {
		attrs = append(attrs, slog.String("service", opts.ServiceName))
	}
	if opts.Version != "" {
		attrs = append(attrs, slog.String("version", opts.Version))
	}

	if len(attrs) > 0 {
		handler = handler.WithAttrs(attrs)
	}

	return slog.New(handler)
}
