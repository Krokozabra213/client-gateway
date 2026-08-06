package config

import (
	"fmt"
	"log/slog"
)

type SlogConfig struct {
	Structured bool `yaml:"structured" env-default:"true"`
	Level      int  `yaml:"level" env-default:"0"`
}

func (c SlogConfig) Validate() error {
	allowed := map[slog.Level]bool{
		slog.LevelDebug: true,
		slog.LevelInfo:  true,
		slog.LevelWarn:  true,
		slog.LevelError: true,
	}

	level := slog.Level(c.Level)
	if !allowed[level] {
		return fmt.Errorf("invalid log level %d: must be one of debug(%d), info(%d), warn(%d), error(%d)",
			c.Level, slog.LevelDebug, slog.LevelInfo, slog.LevelWarn, slog.LevelError)
	}

	return nil
}
