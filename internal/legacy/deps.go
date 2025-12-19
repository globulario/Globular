package legacy

import (
	"io"
	"log/slog"
	"time"
)

// Deps captures the runtime dependencies needed by legacy execution helpers.
type Deps struct {
	Logger *slog.Logger
	Stdout io.Writer
	Stderr io.Writer

	Domain   string
	Protocol string
	Email    string
	Timeout  time.Duration
}
