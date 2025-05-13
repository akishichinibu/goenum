package gen

import (
	"context"
	"log/slog"
)

var Logger = slog.Default()

func init() {
	c := context.Background()
	Logger.Handler().Enabled(c, slog.LevelDebug)
}
