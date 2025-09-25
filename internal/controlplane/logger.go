// controlplane/logger.go
package controlplane

import (
	"fmt"
	"io"

	"github.com/globulario/Globular/internal/logsink"
	"github.com/globulario/services/golang/log/logpb"
)

type Logger struct {
	Debug bool
	out   *logsink.ServiceLogWriter
}

// Debugf implements log.Logger.
func (l Logger) Debugf(format string, args ...interface{}) {
	if l.Debug {
		_, _ = l.out.Write([]byte("[DEBUG] " + fmt.Sprintf(format, args...)))
	}
}

// Errorf implements log.Logger.
func (l Logger) Errorf(format string, args ...interface{}) {
	if l.Debug {
		_, _ = l.out.Write([]byte("[ERROR] " + fmt.Sprintf(format, args...)))
	}
}

// Infof implements log.Logger.
func (l Logger) Infof(format string, args ...interface{}) {
	if l.Debug {
		_, _ = l.out.Write([]byte("[INFO] " + fmt.Sprintf(format, args...)))
	}
}

// Warnf implements log.Logger.
func (l Logger) Warnf(format string, args ...interface{}) {
	if l.Debug {
		_, _ = l.out.Write([]byte("[WARN] " + fmt.Sprintf(format, args...)))
	}
}

func NewLogger(address, app string, echo io.Writer, debug bool) *Logger {
	return &Logger{
		Debug: debug,
		out: logsink.NewServiceLogWriter(
			address,   // Globular address
			app,       // application name, e.g. "xds.ControlPlane"
			"sa",      // user
			"/logger", // method
			logpb.LogLevel_INFO_MESSAGE,
			echo,
		),
	}
}
