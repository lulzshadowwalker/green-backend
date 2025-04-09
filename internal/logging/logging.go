package logging

import (
	"fmt"
	"log/slog"
	"time"

	rl "github.com/lestrrat-go/file-rotatelogs"
)

func init() {
	path := "./logs/app"

	writer, err := rl.New(
		path+"-%Y-%m-%d.log",
		rl.WithLinkName(path+".log"),
		rl.WithRotationTime(24*time.Hour),
		rl.WithMaxAge(7*24*time.Hour),
	)
	if err != nil {
		panic(fmt.Errorf("failed to initialize log file: %w", err))
	}

	handler := slog.NewJSONHandler(writer, &slog.HandlerOptions{
		Level:     slog.LevelInfo,
		AddSource: true,
	})

	logger := slog.New(handler)
	slog.SetDefault(logger)
}
